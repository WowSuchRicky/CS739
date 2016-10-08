#!/bin/bash

mkdir -p int_times

for i in `seq 1 1000`; do
    for j in `seq 0 31`; do
        n=$(echo 2^${j} | bc)
        ./greeter_client 2 $n  >> int_times/${j}.txt
    done
done
