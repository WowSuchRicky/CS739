#!/bin/bash

mkdir -p $1

for i in `seq 1 1000`; do
    ./close -i blah -n 2 >> ${1}/close_times.txt
done
