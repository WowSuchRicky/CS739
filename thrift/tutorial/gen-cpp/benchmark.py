import procnetdev as pndev
import time
import sys
import subprocess

def main():
    
    if len(sys.argv) != 3:
        print "usage: ./benchmark.py <server-name> <#ints to send>"
        sys.exit(1)

    pnd1 = pndev.ProcNetDev()
    start_bytes = pnd1['p3p1']['transmit']['bytes']
    start_time = time.time()
    
    # call the Calculator client
    output = subprocess.check_output(['./Calculator_client', sys.argv[1], sys.argv[2]])
    end_time = time.time()

    print "Time elapsed: ", end_time - start_time
    
    pnd2 = pndev.ProcNetDev()
    end_bytes = pnd2['p3p1']['transmit']['bytes']

    print "Data sent: ", (end_bytes - start_bytes) / 1000, " MB"
    print "Bandwidth: ", (end_bytes - start_bytes) / ((end_time - start_time)/2) / 1000, " MBytes/sec" 

main()
