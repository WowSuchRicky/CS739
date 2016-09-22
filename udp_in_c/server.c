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

int main(int argc, char* argv[]) {

  if (argc != 2) {
    printf("Correct usage is <server REPLY_MESSAGE>\n");
    return 0;
  }

  int sd = udpOpen(SERVER_SOCKET);
  assert(sd > -1);

  while (1) {
    struct sockaddr_in addr;

    // read the message from client (addr is filled in with client addr info during read)
    char message[BUFFER_SIZE];
    memset(message, '\0', BUFFER_SIZE);

    int rc = udpReadRel(sd, &addr, message, BUFFER_SIZE, 0); 
    if (rc > 0) {
      printf("Message from client: %s\n", message);
      udpWriteRel(sd, &addr, argv[1], BUFFER_SIZE, 5);
    }
  }

  return 0;
}
