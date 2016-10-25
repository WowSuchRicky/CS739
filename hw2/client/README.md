NFS client

------
Supports:
   - Dir exploration: 	   ls, cd
   - Reading and writing:  cat, echo redirected into a file works
   - File removal:         rm (with -r for dir or without) work
   - File moving: 	   mv


------
To run FUSE client:

0. (run the server first, using 'go run main.go' in server dir)
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

- ensure that if serve crashes and reboots immediately:
  client doesn't notice anything beyond slow down
  can continue to retry until succeed (idempotency)

- ensure that when one file is open and client A deletes it,
  client B should be able to continue to modify, save it.

- write buffer optimization

- measurements & graphs (after all else is done)

- answer some questions that we'll be asked during meeting
  1) why did we pick gRPC?
  2) others??

- make sure file descriptors are closing in server calls that use them

- (if we have time) fix this: using text editor like emacs doesn't work
   it's saving the buffer correctly (with enclosing pound symbols)
   doesn't seem to be able to successfully overwrite the original
