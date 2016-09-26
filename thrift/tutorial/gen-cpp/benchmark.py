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

    time_elapsed = end_time - start_time;

    print "Time elapsed: ", time_elapsed
    
    pnd2 = pndev.ProcNetDev()
    end_bytes = pnd2['p3p1']['transmit']['bytes']

    bytes_sent = end_bytes - start_bytes
    kbytes_sent = bytes_sent / 1000.0;
    mbytes_sent = kbytes_sent / 1000.0;

    print "Data sent: ", mbytes_sent, " MB"
    print "Bandwidth: ", (mbytes_sent*8) / (time_elapsed/2), " Mbits/sec" 

main()
