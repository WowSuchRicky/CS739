#include <stdlib.h>
#include <stdio.h>
#include <sys/socket.h>
#include <sys/types.h>
#include <netinet/in.h>
#include <string.h>
#include "udp_lib.h"

#define BUFFER_SIZE 1024
#define SERVER_SOCKET 10000
#define CLIENT_SOCKET 20000

int main(int argc, char* argv[]) {

  if (argc != 3) {
    printf("Correct usage is <client SERVER_HOSTNAME MESSAGE>\n");
    return 0;
  }

  int sd, rc;
  struct sockaddr_in addr, addr2;
  sd = udpOpen(CLIENT_SOCKET);
  
  // fill addr with server information and send message to server
  rc = udpFillAddr(&addr, argv[1], SERVER_SOCKET);
  rc = udpWriteRel(sd, &addr, argv[2], BUFFER_SIZE, 10);

  // read and print server reply
  if (rc > 0) {
    char recMessage[BUFFER_SIZE];
    memset(recMessage, '\0', BUFFER_SIZE);
    rc = udpReadRel(sd, &addr2, recMessage, BUFFER_SIZE, 0);
    if (rc > 0) printf("Message received from server: %s\n", recMessage);
  }
  return 0;
}
