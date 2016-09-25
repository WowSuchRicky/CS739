0;136;0cCC = gcc
CFLAGS = -Wall -O0
DEPS = udp_lib.h
LDFLAGS = -lrt

SERVER_HOST = "betelgeuse.cs.wisc.edu"
CLIENT_MSG = "Hello m'server!"
SERVER_MSG = "Get away from me client..."

CLIENT_MSG_MAX = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"

SEND_MSG_N_TIMES = 100		
TIME_TO_RETRY_NS = 1000000
PERCENT_CHANCE_READ_DROP = 50
LOG = "measurements/no_compiler_optimizations/15B_message_size/different_machines/$(PERCENT_CHANCE_READ_DROP)_percent_chance_drop.log"

all: server.o client.o udp_lib.o
	gcc $(CFLAGS) -o server server.o udp_lib.o $(LDFLAGS)
	gcc $(CFLAGS) -o client client.o udp_lib.o $(LDFLAGS)

%.o: %.c %(DEPS)
	%(CC) $(CFLAGS) -c -o $@ $< $(LDFLAGS)

run_client_experiment: client
	./client $(SERVER_HOST) $(CLIENT_MSG) $(SEND_MSG_N_TIMES) $(TIME_TO_RETRY_NS) $(LOG)

run_client: all
	./client $(SERVER_HOST) $(CLIENT_MSG) 1 $(TIME_TO_RETRY_NS)

run_server: all
	./server $(SERVER_MSG) $(PERCENT_CHANCE_READ_DROP)

clean: 	
	rm -f *.o
	rm -f server client
