#!/bin/bash
runBenchmark() {
    for i in 1 10 100 1000 10000
    do
        writeSum=0
        readSum=0
        for j in 1..10
        do
            echo "Running $1 benchmark with $i KB writes and reads"
            output=$(./main -$2 -$2 -size $(($i * 1000)))
            
            writeResult=$(echo "$output" | head -n 1)
            $(($writeSum=$writeSum+$writeResult))
            
            readResult=$(echo "$output | tail -1)
            $(($readSum=$writeSum+$readResult))

            if [ ! -f "$2.csv" ]; then
                $(($writeSum / 10)) >> $2.csv
                $(($readSum / 10)) >> $3.csv
            else
                ,$(($writeSum / 10)) >> $2.csv
                ,$(($readSum / 10)) >> $3.csv
            fi
        done
        echo >> $2.csv
        echo >> $3.csv
    done
    
}

make

runBenchmark "sequential non-blocking" SAW SAR
#runBenchmark "sequential blocking" SBW SBR
#runBenchmark "random non-blocking" RAW RAR
#runBenchmark "random blocking" RBW RBR

#make clean