#include <stdlib.h>
#include <stdio.h>
#include <sys/socket.h>
#include <sys/types.h>
#include <netinet/in.h>
#include <assert.h>
#include <string.h>
#include "udp_lib.h"

#define BUFFER_SIZE 2000
#define SERVER_SOCKET 10000

// What do we need?
// A client (sender) that will send msg, wait, and send again if didn't receive ack
// A server (receiver) that will read a message and send an acknowledgment

int main(int argc, char* argv[]) {

  if (argc != 3) {
    printf("Correct usage is <server REPLY_MESSAGE MSG_DROP_PERCENT_CHANCE>\n");
    return 0;
  }

  int msg_drop_percent = atoi(argv[2]);
  int n, sd;
  char message[BUFFER_SIZE];

  n = 0;
  sd = udp_open(SERVER_SOCKET);
  assert(sd > -1);

  while (1) {
    struct sockaddr_in addr;

    // read message from client and send acknowledgment
    int rc = udp_read_reliable(sd, &addr, message, BUFFER_SIZE, msg_drop_percent); 

    if (rc > 0) { 
      printf("Message %d from client: %s\n", n++, message);
      memset(message, '\0', BUFFER_SIZE);
    }
  }

  return 0;
}
