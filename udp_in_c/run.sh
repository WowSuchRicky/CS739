#!/bin/bash

HOSTNAME=royal-02.cs.wisc.edu
TIME_TO_RETRY_NS=1000000
ITERATIONS_OF_EXPERIMENT=5
vary_data=udp_vary_data

for j in `seq 0 1 $ITERATIONS_OF_EXPERIMENT`; do
    filename=$vary_data$j.csv
    rm -f $filename
    echo "N_BYTES_PER_MSG,N_MSG,TOTAL_BYTES,TOTAL_TIME_NS" >> $filename
    for i in `seq 1 5 1000`; do
	./client $HOSTNAME $i 100 $TIME_TO_RETRY_NS >> $filename
    done
done
