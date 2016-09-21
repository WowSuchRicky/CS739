#include <stdlib.h>
#include <stdio.h>
#include <sys/socket.h>
#include <sys/types.h>
#include <netinet/in.h>
#include <assert.h>
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
    int rc = udpRead(sd, &addr, message, BUFFER_SIZE); 
    printf("Message from client: %s\n", message);
    
    // reply to client
    if (rc > 0) rc = udpWrite(sd, &addr, argv[1], BUFFER_SIZE);
  }

  return 0;
}
