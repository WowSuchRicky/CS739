#include <sys/socket.h>
#include <sys/types.h>

#ifndef H_UDP_LIB

// Open a port and return a file (socket) descriptor for it
int udpOpen(int port);

// Send buffer contents to dest
int udpWrite(int sockfd, struct sockaddr_in* dest, char* buffer, int buffer_len);

// Read message on socket and 1) fill buffer & 2) fill addr with information on sending machine (so you can reply)
int udpRead(int sockfd, struct sockaddr_in* addr, char* buffer, int buffer_len);

// Fill addr with host and port information
int udpFillAddr(struct sockaddr_in* addr, char* hostName, int port);

// Reliable version of UDP write (wait timeout seconds until retry send, blocks until acknowledgment received)
int udpWriteRel(int sockfd, struct sockaddr_in* dest, char* buffer, int buffer_len, int timeout);

// Reliable version of UDP read (sends acknowledgment after each read)
int udpReadRel(int sockfd, struct sockaddr_in* addr, char* buffer, int buffer_len, int dropPercentage);

#define H_UDP_LIB
#endif
