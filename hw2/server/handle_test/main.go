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

    *mid = mount_id;
    return *fhp;
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
    pathname := C.CString("test")
    var mount_id C.int
	fhp := C.get_fhp(pathname, &mount_id)
	fmt.Printf("fhp: %v\n", fhp);
	C.read_fhp(&fhp);
}
