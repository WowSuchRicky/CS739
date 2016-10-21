main.go is the FUSE client.
main_old.go is just used to test NFS protocol via RPC directly.

------
To run FUSE client:

1. To start the client process (continuously running), run:
   go run main.go test

   This will mount the file system exported by the server to the folder called test in this directory
   This also is a continuously running process.

2. In a different process, you can now access that directory.
   - ls seems to work right now (although still adding attribute functionality)
   - you can use a text editor to access some info (but it's only a subset of the file since not taking into account the count)
     - related, cat does not work right now??

3. When you're done, we need to unmount so run:
   fusermount -u test