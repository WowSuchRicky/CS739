Write buffer optimization

References:
Could be good => http://www.scs.stanford.edu/nyu/02fa/notes/l3.pdf
Although I don't think we'll do the exact same thing

-----
Premise: many writes are small, so we're incurring a lot of I/O overhead that could be avoided if we batch them into a single large write.

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
