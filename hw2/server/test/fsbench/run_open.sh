#!/bin/bash

mkdir -p $1

for i in `seq 1 1000`; do
    ./open -i blah -n 2 >> open_times.txt
done
