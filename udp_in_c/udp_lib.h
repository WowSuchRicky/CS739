#include <sys/socket.h>
#include <sys/types.h>

#ifndef H_UDP_LIB

// Open a port and return a file (socket) descriptor for it
int udp_open(int port);

// Send buffer contents to dest
int udp_write(int sockfd, struct sockaddr_in* dest, char* buffer, int buffer_len);

// Read message on socket and 1) fill buffer & 2) fill addr with information on sending machine (so you can reply)
int udp_read(int sockfd, struct sockaddr_in* addr, char* buffer, int buffer_len);

// Fill addr with host and port information
int udp_fill_addr(struct sockaddr_in* addr, char* hostName, int port);

// Reliable version of UDP write 
// Wait for acknowledge (blocking), if no ack after timeout_ns nanoseconds, retry send
int udp_write_reliable(int sockfd, struct sockaddr_in* dest, char* buffer, int buffer_len, 
		       long timeout_ns);

// Reliable version of UDP read 
// Sends acknowledgment after every read
int udp_read_reliable(int sockfd, struct sockaddr_in* addr, char* buffer, int buffer_len, 
		      int dropPercentage);

#define H_UDP_LIB
#endif
