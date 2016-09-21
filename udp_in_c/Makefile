1;2501;0cCC = gcc
CFLAGS = -Wall
DEPS = udp_lib.h

all: server.o client.o udp_lib.o
	gcc $(CFLAGS) -o server server.o udp_lib.o
	gcc $(CFLAGS) -o client client.o udp_lib.o

%.o: %.c %(DEPS)
	%(CC) $(CFLAGS -c -o $@ $<

run_client: client
	./client "betelgeuse.cs.wisc.edu" "Hello m'server!"

run_server: server
	./server "Get away from me client!!!"

clean: 	
	rm -f *.o
	rm -f server client
