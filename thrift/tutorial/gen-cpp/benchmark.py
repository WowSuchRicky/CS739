import procnetdev as pndev
import time
import subprocess

def main():
    pnd1 = pndev.ProcNetDev()
    start_bytes = pnd1['p3p1']['receive']['bytes']
    start_time = time.time()
    
    # call the Calculator client
    output = subprocess.check_output(['./Calculator_client', 'royal-16', '100000000'])
    end_time = time.time()

    print "Time elapsed: ", end_time - start_time
    
    pnd2 = pndev.ProcNetDev()
    end_bytes = pnd2['p3p1']['receive']['bytes']

    print "Data sent: ", end_bytes - start_bytes 

    print "Bandwidth: ", (end_bytes - start_bytes) / (end_time - start_time), " Bytes/sec" 

main()
