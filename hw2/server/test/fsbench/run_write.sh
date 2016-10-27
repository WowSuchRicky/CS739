#!/bin/bash

mkdir -p $1

for i in `seq 1 1000`; do
    for j in `seq 0 21`; do
        n=$(echo 2^${j} | bc)
        ./write -i blah -n 2 -b $n >> write_times_b_${j}.txt
    done
done
