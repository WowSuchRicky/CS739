#include <stdio.h>
#include <stdlib.h>

void usage(char* s) {
    char* cmd_args = "<-i infile> <-o outfile> <-n runs>";
    fprintf(stderr, "Usage: %s %s\n" , argv[0], cmd_args);
    exit(1);
}

int main(int argc, char* argv[]) {
    char *ifname, *ofname;
    int ifd, ofd, runs;
    
    // input params
    int c;
    opterr = 0;
    while ((c = getopt(argc, argv, "i:o:n:")) != -1) {
        switch (c) {
            case 'i':
                ifname = strdup(optarg);
                break;
            case 'o':
                ofname = strdup(optarg);
                break;
            case 'n':
                runs = atoi(optarg);
                break;
            default:
                usage(argv[0]);
    }
    exit(0);
}
