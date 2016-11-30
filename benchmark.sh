#!/bin/bash
runBenchmark() {
    for i in 1 10 100 1000
    do
        writeSum=0
        readSum=0
        
        echo "Running $1 benchmark with $i KB writes and reads"
        for j in `seq 1 10`;
        do
            output=$(./main -$2 -$3 -size $(($i * 1000)))
            
            writeResult=$(echo "$output" | head -n 1)
            writeSum=$((writeSum+writeResult))

            readResult=$(echo "$output" | tail -1)
            readSum=$((readSum+readResult))
        done

        if [ $i -eq 1 ] ; then 
            printf "%d" $(($writeSum/10)) >> $2.csv
            printf "%d" $(($readSum/10)) >> $3.csv
        else
            printf ",%d" $(($writeSum/10)) >> $2.csv
            printf ",%d" $(($readSum/10)) >> $3.csv
        fi
    done

    Rscript generateGraph.R $2.csv
    Rscript generateGraph.R $3.csv
}

make

runBenchmark "sequential non-blocking" SAW SAR
#runBenchmark "sequential blocking" SBW SBR
#runBenchmark "random non-blocking" RAW RAR
#runBenchmark "random blocking" RBW RBR

#make clean
