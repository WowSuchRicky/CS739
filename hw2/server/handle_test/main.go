package main

/*
#define _GNU_SOURCE
#include <errno.h>
#include <stdio.h>
#include <stdlib.h>
#include <sys/types.h>
#include <sys/stat.h>
#include <fcntl.h>

struct file_handle get_fhp(char* pathname, int *mid) {
    struct file_handle *fhp;
    int mount_id, fhsize, dirfd, flags;

    fhsize = sizeof(*fhp);
    if ((fhp = malloc(fhsize)) == NULL) {
        perror("malloc failed\n");
        return *fhp;
    }

    dirfd = AT_FDCWD;
    flags = 0;
    fhp->handle_bytes = 0;

    if (name_to_handle_at(dirfd, pathname, fhp, &mount_id, flags) != -1
        || errno != EOVERFLOW) {
        perror("name_to_handle_at failed\n");
        return *fhp;
    }

    fhsize = sizeof(struct file_handle) + fhp->handle_bytes;
    if ((fhp = realloc(fhp, fhsize)) == NULL) {
        perror("realloc failed\n");
        return *fhp;
    }

    if (name_to_handle_at(dirfd, pathname, fhp, &mount_id, flags) == -1) {
        perror("name_to_handle_at failed\n");
        return *fhp;
    }

    fprintf(stdout, "fhp->handle_bytes: %d\n", fhp->handle_bytes);
    fprintf(stdout, "fhp->handle_type: %d\n", fhp->handle_type);
    fprintf(stdout, "fhp->f_handle: %d\n", fhp->f_handle);
    fprintf(stdout, "mount_id: %d\n", mount_id);

	int mount_fd = open("/", O_RDONLY);
	int fd = open_by_handle_at(mount_fd, fhp, O_RDONLY);
	char buf[100];
	int nread = read(fd, buf, sizeof(buf));
	printf("fd in getting fd: %d\n", fd);
	printf("Read %s \n", buf);
	printf("Read %zd bytes\n", nread);

    *mid = mount_id;
    return *fhp;
}

int fh_open(struct file_handle fhp) {
	char buf[10];
	int mount_fd = open("test", O_RDONLY);
	int fd = open_by_handle_at(mount_fd, &fhp, O_RDONLY);
	printf("fd: %d\n", fd);
	int nread = read(fd, buf, sizeof(buf));
	printf("Read %s \n", buf);
	printf("Read %zd bytes\n", nread);

	
	return fd;
}

void read_fhp(struct file_handle *fhp) {
    fprintf(stdout, "fhp->handle_bytes: %d\n", fhp->handle_bytes);
    fprintf(stdout, "fhp->handle_type: %d\n", fhp->handle_type);
    fprintf(stdout, "fhp->f_handle: %d\n", fhp->f_handle);
}
*/
import "C"
import "fmt"

func main() {
    pathname := C.CString("test/test.txt")
    var mount_id C.int
	fhp := C.get_fhp(pathname, &mount_id)
	fmt.Printf("fhp: %v\n", fhp);
	C.read_fhp(&fhp);
	C.fh_open(fhp);
}
