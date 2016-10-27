Write buffer optimization

-----
Premise:
Many writes are small, so we're incurring a lot of I/O overhead that could be avoided if we batch them into a single large write.

-----
How it works:
Each client maintains a queue.
This queue contains information for writes that have been sent to but
not necessarily persisted on the server.
There is a maximum size to this queue.
The server has a queue as well; it stores all the writes that it hasn't
persisted yet in that queue. The server's queue is unbounded.

When a client writes, what we do depends on this queue and the write:
1. If the write is huge (> capacity of the queue), it is sent directly
   to the server as a stable (synchronous) write. It does not go in the queue.
2. If the queue has space for the write, the client will send an
   unstable (asynchronous) write to the server, and put that request in
   the queue as well. This allows it to be retried later.
3. If the write is not greater than queue size, but it won't currently fit in
   the queue, the client will issue a COMMIT to the server, telling it
   to persist everything in its queue. When this commit returns, the client
   frees its queue, and places that write inside it

Every WRITE and COMMIT call will return a 'writeverf3' number from the server.
This number is persisted on the server, and increments each time the server
restarts.
The client keeps track of the last known writeverf3 it received.
If it ever sends a WRITE or COMMIT request and notices the writeverf3 number
has increased, it will resend everything in its queue; that increase means
the server crashed/reset and thus lost its volatile state. Doing this will catch the server up to the client.


NOTE: A read request from the client will flush the buffer.
      This means some benefit is lost if there's even a single read in the middle
      of a write-heavy workload.
      The reason this is done is because on the server, we might be reading
      a file that has been written to but the writes haven't been persisted yet,
      so we would need to take them into account to ensure correctness.
      The easy way is to flush the buffer to disk, like we're doing, and then
      read the file.
      Alternatives could include bringing the file contents into memory, and
      then applying the buffered write requests. This means we would be
      going through the server write queue on every read request. This
      penalty would be lessened if we keep a queue for every file that had
      a write request.
      

NOTE: The server can have multiple clients.
      Since the server only stores one queue, it doesn't differentiate between
      clients. When it receives a commit request, it persists everything, from all
      clients. This means that the client queues might get out of sync from the
      server queue without the server crashing.
      There's a way around this but it's not currently implemented.
      The idea is to keep an additional number in the server which increments on
      every commit. This would be returned by various requests, and the client keeps
      track of the last one it saw. If at any point the client sees an increase
      but the client didn't call COMMIT itself, it knows that the server queue
      was persisted and thus it should free its queue as well.






-----------------
IGNORE BELOW THIS 
-----------------

References:
Could be good => http://www.scs.stanford.edu/nyu/02fa/notes/l3.pdf
Although I don't think we'll do the exact same thing

-----
How should this work?
 - should be transparent to client
 - application calls write(), and FUSE will pick up that write; we need to control
   operation in a specific way from there

NFS protocol:
 - Extend existing NFS WRITE protocol to
   Include a STABLE field (boolean) in its args
   Include a writeverf3 field in its return value
 - Add a COMMIT call to the protocol which:
   Includes nothing in its args
   Includes a writeverf3 field (integer) in its return value
 - FSYNC on client should call NFS COMMIT
 
CLIENT procedure (in FUSE):
  1. Client maintains a write buffer which is essentially a list of NFS WriteArgs; this
     buffer has a predetermined maximum size.

  2. If the write is greater than that threshold, we use NFS STABLE WRITE

  3. If the write is smaller than that threshold:
     a) if the available memory in the client data structure + memory of this write
     	is less than the capacity of that data structure, store it in the data structure
	and issue an UNSTABLE write
     b) if the available memory is less, issue a COMMIT; once acknowledged, free the
     	entire data structure and issue this as an UNSTABLE write

  (aside: every NFS WRITE and COMMIT returns the writeverf3 number;
  	  client maintains the verf number
  	  that it saw at last commit; if a write or commit returns a different one,
	  it should resend everything that is in its buffer.

  *NOTE: we do all of this because the client must maintain all the data to be
  	 written which aren't known to have committed in case it must resend;
	 the application itself shouldn't need to make sure that data exists
	 still, so it's saved in the NFS client

SERVER procedure:
  1. When server boots up, read the last writeverf3 number from disk.
     Increment it and persist it to disk immediately.
     
  2. Server maintains some structures:
     List of WriteRequests; these are intermingled without regard to client.
     This list is ordered by when the request was received by server.

  3. When server receives NFS write, do the following:
     Check stable field ->
       If UNSTABLE,
         Add the entry to the WriteRequest list
         Return ACK
       If STABLE, write with normal NFS semantics.

  4. When server receives NFS commit, do the following:
     Go through the list for the committing client
     Perform all the I/Os (NOTE: this entire optimization hinges on the fact
     	     	     	  	 that these I/Os can be batched at this point;
				 how do we do that?!)
     Free list
     Send acknowledgment

ASSUMPTIONS
  The client has an explicit maximum write buffer size.
  We assume the following always holds:
     available memory on server > number of clients * client buffer capacity
  We could keep track of clients, but that requires maintaing additional state and
  requires generating unique numbers (even IP address wouldn't suffice, because you
  could have multiple clients from the same machine). Instead, the server will
  commit all changes from all clients when a single commit is received.


TODO:
 - (done) basic process working
 - (done) add fsync to fuse client; should just call commit
 - (done) change reads so each one persists; ensure the read takes into account everything, even
 - 	  the data that is persisted at the beginning of the read (this requires a modification to
 -	  get attr, because we call that to know the size of the file to read, in the case of
 -	  ReadAll. WE could do better here.
 - utilize writever3 num so that a client will know to resend everything in it's buffer if server crashed
 - add nfs_Commit number so clients know they should flush their own buffers
 - perform measurements

This optimization is good for write-heavy; if a lot of files are written at once, it will work well. If files are read in-between, it would require more work to be efficient.

 - (ignore) change reads so that on each read, the write buffer on server will be scanned; if any changes to the specific to-be-read file are found, apply them in memory and return it