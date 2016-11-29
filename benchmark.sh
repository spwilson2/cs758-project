#!/bin/bash
make

for i in 1 10 100 1000 10000
do
    echo "Running sequentional non-blocking benchmark with $i KB writes and reads"
    ./main -SAW -SAR -size $(($i * 1000))
done

make clean