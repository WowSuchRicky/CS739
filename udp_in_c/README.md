

## What's Here?
1. Simple UDP library (it's not reliable yet!)
2. Client for sending message to server
3. Server for receiving messages and sending reply

## To Use:

1. make all

2. On the machine that you want to be server (keep in mind hostname)...
   make run_server

3. Edit run_client rule in Makefile; change the existing hostname to hostname from step 2

3. On the machine that you want to be client...
   make run_client 

4. If successful:
   - server instance will print the message
   - client will print server's reply

## Other Notes:
Editing run_client and run_server in makefile will let you:
- change the client -> server message
- change the server -> client message
- change the server hostname for client to connect to

If you need to change socket port numbers, edit server.c and client.c directly.


