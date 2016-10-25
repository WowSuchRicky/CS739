main.go is the FUSE client.
main_old.go is just used to test NFS protocol via RPC directly.


------
Supports:
   - Dir exploration: 	   ls, cd
   - Reading and writing:  cat, echo redirected into a file works
   - File removal:         rm (with -r for dir or without) work


------
To run FUSE client:

1. ./run.sh
2. ./client test (this will mount in test/ within the current directory)

-OR-

1. To start the client process (continuously running), run:
   go run main.go test
   (this will mount the file system exported by the server to the folder called test in this directory
   This also is a continuously running process.)
2. In a different process, you can now access that directory.
3. When you're done, we need to unmount so run:
   fusermount -u test


------
TODO:

2. (Alex) using text editor like emacs doesn't work
   it's saving the temp file correctly (with enclosing pound symbols)
   doesn't seem to be able to successfully overwrite the original

3. ensure that if serve crashes and reboots immediately:
   client doesn't notice anything beyond slow down
   can continue to retry until succeed (idempotency)

4. ensure that when one file is open and client A deletes it,
   client B should be able to continue to modify, save it.

5. write buffer optimization

6. measurements & graphs


DONE:
1. (solved) rename isn't working fully, so 'mv' doesn't work
   problem identified: need to access the destination dir filehandle, and this is not straightforward in FUSE for rename
