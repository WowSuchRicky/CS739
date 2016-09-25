

## What's here?
1. Simple UDP library (unreliable + reliable APIs)
2. Client for sending message to server
3. Server for receiving messages and sending acknowledgment
4. Measurements results from various experiments (see below)

## To use:

1. make all

2. On the machine that you want to be server (keep in mind hostname), run:
   make run_server

3. Change SERVER_HOST variable in Makefile to the hostname from step 2.

3. On the machine that you want to be client, run:
   make run_client 

4. If successful:
   - client will send a single message to the server
   - server will print the message
   - server will send an acknowledgment in return
   - client will print notification that ack was received
   - the time in seconds and nanoseconds for roundtrip (client -> server -> client) will be recorded in a log file

## Other notes:
- If you need to change socket port numbers, edit server.c and client.c directly
- Most other interesting parameters can be changed in the Makefile

## To perform measurements:
1. Tweak settings in Makefile as desired (optimizations, time between send retries, etc.)
2. Follow the above usage steps, but use "make run_client_experiment" at step 3.

## Completed measurements:
1. Roundtrip length of 100 messages sent from process on one machine to process on the same machine in succession with (0/10/20/30/50/90)% chance of dropping
2. Roundtrip length of 100 messages sent from process on one machine to process on a different machine (across network) in succession with (0/10/20/30/50/90%) chance of dropping
3. Repeated #1 and #2 with -O3 optimizations (as opposed to -O0 optimizations, i.e. none)
4. Bandwidth measurements (when sending large number of max-sized packets) on one machine
5. Bandwidth measurements across machines
6. Repeated #4 and #5 bandwidth measurements with -O3 optimizations

## TODO:
1. Graph those results
2. Answer question (think U-net paper): how much overhead is there to send a message?
3. Answer questions: what limits the bandwidth, and how could we do better?
