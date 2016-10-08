#!/bin/bash

vary_data=streaming_times_vary_data_per_stream.csv

rm -f $vary_data touch $vary_data

echo "N_STREAMS,N_BYTES_PER_STREAM,TOTAL_BYTES,ROUNDTRIP_TIME_NS" >> $vary_data

for i in `seq 0 10 10000`; do
    ./greeter_client 100 $i >> $vary_data
done
