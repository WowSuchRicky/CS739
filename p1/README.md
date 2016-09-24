How did we take these measurements, or why did we not take one?


-Mutex lock: mutex/bench-mutex.c w RDTSCP
-Read 1MB Sequentially from SSD/HDD (Different Machines): dd -if /tmp/1MB.dat -of /dev/null
-Packet from CA->Netherlands->CA, used Ping, got 112ms every time. 
