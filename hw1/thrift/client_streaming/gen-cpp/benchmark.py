import procnetdev as pndev
import time
import sys
import subprocess

def main():
    
    if len(sys.argv) != 2:
        print "usage: ./benchmark.py <server-name>"
        sys.exit(1)

   
    print "TOTAL_BYTES,TIME_NS"
    for i in range(0, 1000000, 10000):
        # call the Calculator client
        output = subprocess.check_output(['./Calculator_client', sys.argv[1], str(i)])
        print output


main()
