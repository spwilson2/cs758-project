#!/bin/bash
runBenchmark() {
    for i in 1 10 100 1000
    do
        writeSum=0
        readSum=0
        for j in `seq 1 10`;
        do
            echo "Running $1 benchmark with $i KB writes and reads"
            output=$(./main -$2 -$3 -size $(($i * 1000)))
            
            writeResult=$(echo "$output" | head -n 1)
            let writeSum=writeSum+$writeResult

            readResult=$(echo "$output" | tail -1)
            let readSum=readSum+$readResult

            if [ $j -eq 1 ] ; then 
                printf "%d" $(($writeSum/10)) >> $2.csv
                printf "%d" $(($readSum/10)) >> $3.csv
            else
                printf ",%d" $(($writeSum/10)) >> $2.csv
                printf ",%d" $(($readSum/10)) >> $3.csv
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
