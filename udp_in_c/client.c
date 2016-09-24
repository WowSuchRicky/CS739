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
  struct sockaddr_in addr;
  sd = udpOpen(CLIENT_SOCKET);
  
  int msg_size = 0;
  while (*(argv[2] + msg_size++) != '\0') { }
  msg_size++; // to account for null char

  // get address info for hostname
  rc = udpFillAddr(&addr, argv[1], SERVER_SOCKET);
  
  // send message to that host (TODO: wrap timer around this for measurements)
  rc = udpWriteRel(sd, &addr, argv[2], msg_size, 2);

  // write procedure blocks until ack, so we know it was received
  printf("Ack received thus message was received!\n");

  return 0;
}
