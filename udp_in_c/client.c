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

#define BUFFER_SIZE 1024
#define SERVER_SOCKET 10000
#define CLIENT_SOCKET 20000

#define WRITE_TO_CONSOLE 1
#define WRITE_TO_LOG 1
#define DEFAULT_LOG_PATH "./results.log"

// send message to hostname, repeat n_times
int send_message(char* hostname, char* message, int n_times, 
		 long time_to_retry, char* log_path) {

  int sd, rc, i, start_time_s, total_time_s;

  long start_time_ns, total_time_ns;
  struct sockaddr_in addr;
  struct timespec time;
  assert (clock_getres(CLOCK_REALTIME, &time) != -1);

  FILE* log_fd = fopen(log_path, "w"); // will truncate
  fprintf(log_fd, "seconds,nanoseconds\n");
  
  assert (log_fd != NULL);

  sd = udp_open(CLIENT_SOCKET);
  assert (sd >= 0);

  int msg_size = 0;
  while (*(message + msg_size++) != '\0') { }
  msg_size++; // to account for null char

  // get address info for hostname
  if (udp_fill_addr(&addr, hostname, SERVER_SOCKET) < 0) return -1;

  // repeat sending the actual message and recording results n_times
  for (i = 0; i < n_times; i++) {
    assert (clock_gettime(CLOCK_REALTIME, &time) != -1);
    start_time_ns = time.tv_nsec;
    start_time_s = time.tv_sec;

    rc = udp_write_reliable(sd, &addr, message, msg_size, time_to_retry);
    if (rc < 0) return -1;

    assert (clock_gettime(CLOCK_REALTIME, &time) != -1);
    total_time_ns = time.tv_nsec - start_time_ns; 
    total_time_s = time.tv_sec - start_time_s;

    // write will block until ack, so when we reach here we know it was received
    if (WRITE_TO_CONSOLE) 
      printf("Acknowledge received; trip time: %d sec, or %ld nsec.\n", 
	     total_time_s, total_time_ns);

    if (WRITE_TO_LOG) 
      fprintf(log_fd, "%d,%ld\n", total_time_s, total_time_ns);
  }

  fclose(log_fd);
  return 0;
}


int main(int argc, char* argv[]) {

  if (argc < 5) {
    printf("Correct usage is <client SERVER_HOSTNAME MESSAGE N_TIMES TIME_TO_RETRY>\n");
    return 0;
  }

  char* log_path;
  if (argc > 5) log_path = argv[5];
  else log_path = DEFAULT_LOG_PATH;

  char* hostname = argv[1];
  char* message_to_send = argv[2];
  int n_times = atoi(argv[3]);
  long time_to_retry = atol(argv[4]);


  if (send_message(hostname, message_to_send, n_times, time_to_retry, log_path) < 0) 
    return -1;

  return 0;
}
