IPs of VMs to run on:
104.197.218.40

----------------
Things to do:
1. (done) decide what should be in a filehandle; decided inode and genum
2. (done) Add way to retrieve inode
3. (done) Add way to retrieve genum
4. (in progress) Implement server-side calls - these are same as NFS protocol

   (done) lookup, create, remove, read, write
   (to be completed) readdir, ???
   
   Note: done here means partially done; still need to figure out
   	 some things before can be fully completed - those things are
	 outlined in comments in the file and essentially line up with the
	 next points

5. (done) Figure out what attributes are necessary to include in various calls
   - just using what is in the bazilfuse attr; this mostly overlaps with what stat
     returns, although stat does return even more
6. (done) Figure out how to use the given genum in each of the server-side calls to
   ensure correctness
7. Figure out how to ensure idempotency, or what design decisions we want to make
   (this somewhat goes along with properly using err's in server side functions)
8. Integrate client using FUSE


