package nfs

// #include <sys/ioctl.h>
// #include <sys/fcntl.h>
// #include <linux/fs.h>
// #include <stdio.h>
// #include <errno.h>
// unsigned long int pathtogen (char *path) {
//    int fileno = open(path, O_RDONLY);
//    unsigned int generation = 0;
//    if (ioctl(fileno, FS_IOC_GETVERSION, &generation)) {
//       printf("errno: %d\n", errno);
//    }
//    return generation;
// }
import "C"

func PathToGen(path string) (uint64, error) {
	return uint64(C.pathtogen(C.CString(path))), nil
}
