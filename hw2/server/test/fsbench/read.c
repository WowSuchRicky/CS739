#include <fcntl.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/stat.h>
#include <sys/types.h>
#include <time.h>
#include <unistd.h>

void usage(char* s) {
    char* cmd_args = "<-i infile> <-n runs> <-b bytes to read>";
    fprintf(stderr, "Usage: %s %s\n" , s, cmd_args);
    exit(1);
}

int main(int argc, char* argv[]) {
    char *fname;
    int bytes, c, i, fd, runs;
    double time;
    
    // input params
    opterr = 0;
    while ((c = getopt(argc, argv, "i:b:n:")) != -1) {
        switch (c) {
            case 'i':
                fname = strdup(optarg);
                break;
            case 'b':
                bytes = atoi(optarg);
                break;
            case 'n':
                runs = atoi(optarg);
                break;
            default:
                usage(argv[0]);
        }
    }

    // Generate bytes
    char *buf = (char *) malloc(bytes*sizeof(char));

    
    // Time each file opening
    struct timespec start, end;
    for (i = 0; i < runs; i++) {
        int rc;
        if ((fd = open(fname, O_RDONLY, 0666)) < 0) {
            perror("open input file failed");
            exit(1); // If open fails, we should know about it & debug!
        }
        clock_gettime(CLOCK_REALTIME, &start);
        rc = read(fd, buf, bytes*sizeof(char));
        clock_gettime(CLOCK_REALTIME, &end);

        if (rc < bytes) {
            fprintf(stderr, "read only read %d expected %d\n", rc, bytes);
            exit(1);
        }
        close(fd); // B/c I don't think we have lseek

        time = 1000000000*(end.tv_sec-start.tv_sec)+(end.tv_nsec-start.tv_nsec);
        fprintf(stdout, "%f\n", time);
    }

    free(fname);
    exit(0);
}
