#include <stdlib.h>
#include <stdio.h>
#include <sys/types.h>
#include <sys/stat.h>
#include <fcntl.h>

#include <string.h>
#include <errno.h>

int main(int argc, char *argv[]){
	int rc = creat("test/test2.txt", O_CREAT);

	printf("errno %s\n", strerror(errno));
}