import procnetdev as pndev
import time
import sys
import subprocess

def main():
    
    if len(sys.argv) != 2:
        print "usage: ./benchmark.py <server-name>"
        sys.exit(1)

    # 1K, 10K, 100K, 1M, 10M, 100M, 500M
    nums_to_try = ['1000', '10000', '100000', '1000000', '10000000', '100000000', '500000000']

    
    for i in nums_to_try:
        # call the Calculator client
        output = subprocess.check_output(['./Calculator_client', sys.argv[1], i])
        print output


main()
