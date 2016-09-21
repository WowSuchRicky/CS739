#include <stdlib.h>
#include <stdio.h>
#include <sys/socket.h>
#include <sys/types.h>
#include <netinet/in.h>
#include <string.h>
#include <unistd.h>
#include <netdb.h>

// Basic UDP on top of sockets implementation referenced from:
// https://www.cs.rutgers.edu/~pxk/417/notes/sockets/udp.html
// http://pages.cs.wisc.edu/~remzi/OSTEP/dist-intro.pdf

// What kind of API do we need?
//  1) open a port (which creates socket)
//  2) write to a socket
//  3) read from a socket
//  4) populate sockaddr_in structure

int udpOpen(int port) {
  
  int sd; // file descriptor for socket

  // create socket (AF_INET = IPv4, SOCK_DGRAM = unreliable datagrams)
  if ((sd = socket(AF_INET, SOCK_DGRAM, 0)) < 0) {
    fprintf(stderr, "Could not create socket.\n");
    return -1;
  }

  // name the socket (assign port number to it) i.e. bind the address
  // to bind, need a sockaddr structure (as opposed to a port; because interface is general)
  struct sockaddr_in saddr;
  memset((char*)&saddr, 0, sizeof(saddr));

  saddr.sin_family = AF_INET;                // address family (IPv4)
  saddr.sin_port = htons(port);              // allow OS to assign port number
  saddr.sin_addr.s_addr = htonl(INADDR_ANY); // current machine's IP addr
  if (bind(sd, (struct sockaddr*)&saddr, sizeof(saddr)) < 0) {
    close(sd);
    fprintf(stderr, "Bind failed.\n");
    return -1;
  }

  return sd;
}

int udpWrite(int sockfd, struct sockaddr_in* dest, char* buffer, int buffer_len) {
  return sendto(sockfd, buffer, buffer_len, 0, (struct sockaddr*)dest, 
		sizeof(struct sockaddr_in));
}

int udpRead(int sockfd, struct sockaddr_in* addr, char* buffer, int buffer_len) {
  int addrLen = sizeof(struct sockaddr_in);
  return recvfrom(sockfd, buffer, buffer_len, 0, (struct sockaddr*)addr, 
		  (socklen_t*)&addrLen);
}

int udpFillAddr(struct sockaddr_in* addr, char* hostName, int port) {
  memset((char*)addr, 0, sizeof(struct sockaddr_in));
  addr->sin_family = AF_INET;
  addr->sin_port = htons(port);
  
  struct in_addr* inAddr;
  struct hostent* hostEntry;

  // retrieve information for a specific hostname, pull the host addr out of that
  if ((hostEntry = gethostbyname(hostName)) == NULL) {
    fprintf(stderr, "Could not translate hostName\n");
    return -1;
  }
  inAddr = (struct in_addr*) hostEntry->h_addr;

  addr->sin_addr = *inAddr;
  return 0;
}
