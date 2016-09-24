#include <stdlib.h>
#include <stdio.h>
#include <sys/socket.h>
#include <sys/types.h>
#include <netinet/in.h>
#include <string.h>
#include <unistd.h>
#include <netdb.h>
#include <time.h>
#include <assert.h>

// Basic UDP on top of sockets implementation referenced from:
// https://www.cs.rutgers.edu/~pxk/417/notes/sockets/udp.html
// http://pages.cs.wisc.edu/~remzi/OSTEP/dist-intro.pdf

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

// returns the bytes read, or -1 if error
int udpRead(int sockfd, struct sockaddr_in* addr, char* buffer, int buffer_len) {
  int addrLen = sizeof(struct sockaddr_in);
  return recvfrom(sockfd, 
		  buffer, buffer_len, 
		  MSG_DONTWAIT,      // non-blocking read (change to 0 for blocking)
		  (struct sockaddr*)addr, 
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

//////////////////////
// RELIABLE VERSIONS!!
//////////////////////

// Reliable UDP write
// 1) Send message
// 2) Wait for timeout secs while reading (essentially polling) for reply from receiver
// 3) If no acknowledgment from receiver in that time, return to step 1
int udpWriteRel(int sockfd, struct sockaddr_in* dest, char* buffer, int buffer_len, int timeout) {

  int rc, startTime;
  int ackRec = 0;

  struct timespec *time;
  time = (struct timespec *)malloc(sizeof(struct timespec));
  assert (clock_getres(CLOCK_REALTIME, time) != -1);

  char* ack = (char*)malloc(sizeof(char));
  int ack_len = sizeof(char);
  memset(ack, 0, ack_len);

  // repeat send attempt every 'timeout' seconds until receive ack
  while (!ackRec) {
    printf("Attempting send.\n");
    rc = udpWrite(sockfd, dest, buffer, buffer_len);

    // repeatedly check for acknowledgment until timoeut    
    assert (clock_gettime(CLOCK_REALTIME, time) != -1);
    startTime = time->tv_sec;
    while (time->tv_sec - startTime < timeout) {

      // use unreliable read, we don't want to send ack of an ack
      udpRead(sockfd, dest, ack, ack_len);
      if (*ack == 1) {
	ackRec = 1;
	break;
      }
      assert (clock_gettime(CLOCK_REALTIME, time) != -1);
    }

    if (ackRec) break;
  }
  
  free(ack);
  free(time);

  return rc;
}


// Reliable UDP read (with percent chance of dropping it i.e. taking no action)
// 1) read the data
// 2) send a single byte with value 1 as acknowledgment back to the sender
int udpReadRel(int sockfd, struct sockaddr_in* addr, char* buffer, int buffer_len, int dropPercentage) {

  // simulate a dropped call by clearing socket buffer and returning (don't send ack)
  int randNum = rand() % 100;
  if (randNum < dropPercentage) {
    udpRead(sockfd, addr, buffer, buffer_len);
    return -1;
  }

  int rc = udpRead(sockfd, addr, buffer, buffer_len);
  if (rc <= 0) return rc;

  // send single-byte acknowledgement only if we read something
  printf("Message successfully read, sending ack.\n");
  char* ack = (char*)malloc(sizeof(char));
  memset(ack, 1, sizeof(char));
  udpWrite(sockfd, addr, ack, sizeof(char));
  
  free(ack);
  return rc;
}
