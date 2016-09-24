How did we take these measurements, or why did we not take one?

-Read 1MB Sequentially from Memory: `dd if=/dev/mem of=/tmp/739 bs=1024000 count=1`
-Read 4K Randomly from SSD/HDD: `dd if=/dev/sda of=/tmp/739 bs=512 count=8`
-Mutex lock: mutex/bench-mutex.c w RDTSCP
-Read 1MB Sequentially from SSD/HDD (Different Machines): `dd -if /tmp/1MB.dat -of /dev/null`
-Packet from CA->Netherlands->CA, used Ping, got 112ms every time. 
