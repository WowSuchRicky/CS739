#include <stdlib.h>
#include <stdio.h>
#include <sys/socket.h>
#include <sys/types.h>
#include <netinet/in.h>
#include <assert.h>
#include <string.h>
#include "udp_lib.h"

#define BUFFER_SIZE 1024
#define SERVER_SOCKET 10000
#define MSG_DROP_PERCENT 50

// What do we need?
// A client (sender) that will send msg, wait, and send again if didn't receive ack
// A server (receiver) that will read a message and send an acknowledgment

int main(int argc, char* argv[]) {

  if (argc != 2) {
    printf("Correct usage is <server REPLY_MESSAGE>\n");
    return 0;
  }

  int sd = udpOpen(SERVER_SOCKET);
  assert(sd > -1);
  char message[BUFFER_SIZE];
  int msgN = 0;

  while (1) {
    struct sockaddr_in addr;

    // read message from client and send acknowledgment
    int rc = udpReadRel(sd, &addr, message, BUFFER_SIZE, MSG_DROP_PERCENT); 

    if (rc > 0) { 
      printf("Message %d from client: %s\n", msgN++, message);
      memset(message, '\0', BUFFER_SIZE);
    }
  }

  return 0;
}
