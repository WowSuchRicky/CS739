#include <sys/ioctl.h>
#include <sys/fcntl.h>
#include <linux/fs.h>
#include <stdio.h>
#include <errno.h>

int main (int argc, char** argv) {
    int fileno = open(argv[1], O_RDONLY);
    printf("fileno: %d\n", fileno);
    unsigned int generation = 0;
    if (ioctl(fileno, FS_IOC_GETVERSION, &generation)) {
        printf("errno: %d\n", errno);
    }
    printf("generation: %u\n", generation);
}
