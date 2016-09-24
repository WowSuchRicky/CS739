

## What's Here?
1. Simple UDP library (it's not reliable yet!)
2. Client for sending message to server
3. Server for receiving messages and sending reply

## To Use:

1. make all

2. On the machine that you want to be server (keep in mind hostname), run:
   make run_server

3. Change SERVER_HOST variable in Makefile to the hostname from step 2.

3. On the machine that you want to be client, run:
   make run_client 

4. If successful:
   - server instance will print the message
   - client will print notification that ack was received

## Other Notes:
If you need to change socket port numbers, edit server.c and client.c directly.


