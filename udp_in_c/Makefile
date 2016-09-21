CC = gcc
CFLAGS = -Wall
DEPS = udp_lib.h
LDFLAGS = -lrt

all: server.o client.o udp_lib.o
	gcc $(CFLAGS) -o server server.o udp_lib.o $(LDFLAGS)
	gcc $(CFLAGS) -o client client.o udp_lib.o $(LDFLAGS)

%.o: %.c %(DEPS)
	%(CC) $(CFLAGS) -c -o $@ $< $(LDFLAGS)

run_client: client
	./client "betelgeuse.cs.wisc.edu" "Hello m'server!"

run_server: server
	./server "Get away from me client!!!"

clean: 	
	rm -f *.o
	rm -f server client
