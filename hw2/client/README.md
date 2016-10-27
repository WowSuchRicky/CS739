NFS client

------
Supports:
   - Dir exploration: 	   ls, cd
   - Reading and writing:  cat, echo redirected into a file works
   - File removal:         rm (with -r for dir or without) work
   - File moving: 	   mv


------
To enable/disable write buffer optimization:
Set the boolean in client/main.go to enable/disable (true is enable)


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
