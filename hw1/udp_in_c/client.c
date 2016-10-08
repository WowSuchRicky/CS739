#include <stdlib.h>
#include <stdio.h>
#include <sys/socket.h>
#include <sys/types.h>
#include <netinet/in.h>
#include <string.h>
#include <time.h>
#include <assert.h>
#include <unistd.h>
#include "udp_lib.h"

#define BUFFER_SIZE 2000
#define SERVER_SOCKET 10000
#define CLIENT_SOCKET 20000

// send message to hostname, repeat n_times
int send_message(char* hostname, int msg_n_bytes, int n_times, long time_to_retry) {

  int sd, rc, i;
  long total_time;
  struct sockaddr_in addr;
  struct timespec start_time, end_time;
  assert (clock_getres(CLOCK_REALTIME, &start_time) != -1);
  assert (clock_getres(CLOCK_REALTIME, &end_time) != -1);

  // create message
  char message[msg_n_bytes];
  for (i = 0; i < msg_n_bytes; i++) message[i] = 'a';

  // port to send messages
  sd = udp_open(CLIENT_SOCKET);
  assert (sd >= 0);

  // get address info for hostname
  if (udp_fill_addr(&addr, hostname, SERVER_SOCKET) < 0) return -1;

  // send actual (msg_n_bytes) bytes n_times , and time the total amount
  assert (clock_gettime(CLOCK_REALTIME, &start_time) != -1);
  for (i = 0; i < n_times; i++) {
    rc = udp_write_reliable(sd, &addr, message, msg_n_bytes, time_to_retry); // it's blocking until receive ack
    if (rc < 0) return -1;
  }
  assert (clock_gettime(CLOCK_REALTIME, &end_time) != -1);
  total_time = 1000000000 * (end_time.tv_sec - start_time.tv_sec) + end_time.tv_nsec - start_time.tv_nsec;

  // print results
  int total_bytes = msg_n_bytes * n_times + n_times;
  printf("%d,%d,%d,%ld\n", msg_n_bytes, n_times, total_bytes, total_time);
  
  return 0;
}


int main(int argc, char* argv[]) {

  if (argc < 5) {
    printf("Correct usage is <client SERVER_HOSTNAME MSG_N_BYTES N_TIMES TIME_TO_RETRY>\n");
    return 0;
  }

  char* hostname = argv[1];
  int msg_n_bytes = atoi(argv[2]);
  int n_times = atoi(argv[3]);
  long time_to_retry = atol(argv[4]);

  if (send_message(hostname, msg_n_bytes, n_times, time_to_retry) < 0) 
    return -1;

  return 0;
}
