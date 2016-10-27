#!/bin/bash

mkdir -p $1

for i in `seq 1 20`; do
    ./open -i blah -n 2 >> ${1}/open_times.txt
done
