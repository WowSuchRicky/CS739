CC = gcc
CFLAGS = -Wall -O0
DEPS = udp_lib.h
LDFLAGS = -lrt

SERVER_HOST = "royal-02.cs.wisc.edu"
SERVER_MSG = "Get away from me client..."
# msg_max before was 2000 char
# small msg was 15 char

N_CHAR_TO_SEND = 2000
SEND_MSG_N_TIMES = 100		
TIME_TO_RETRY_NS = 1000000
PERCENT_CHANCE_READ_DROP = 0
LOG = results.csv
# LOG = "measurements/with_compiler_optimizations/512B_message_size/different_machines/$(PERCENT_CHANCE_READ_DROP)_percent_chance_drop.log"

all: server.o client.o udp_lib.o
	gcc $(CFLAGS) -o server server.o udp_lib.o $(LDFLAGS)
	gcc $(CFLAGS) -o client client.o udp_lib.o $(LDFLAGS)

%.o: %.c %(DEPS)
	%(CC) $(CFLAGS) -c -o $@ $< $(LDFLAGS)

run_client_experiment: client
	./client $(SERVER_HOST) $(N_CHAR_TO_SEND) $(SEND_MSG_N_TIMES) $(TIME_TO_RETRY_NS) $(LOG)

run_client: all
	./client $(SERVER_HOST) $(N_CHAR_TO_SEND) 1 $(TIME_TO_RETRY_NS)

run_server: all
	./server $(SERVER_MSG) $(PERCENT_CHANCE_READ_DROP)

clean: 	
	rm -f *.o
	rm -f server client
