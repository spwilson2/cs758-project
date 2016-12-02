#!/usr/bin/python

from __future__ import print_function
import subprocess
import re
import sys
import os

V_SPEW=2
V_INFO=1
V_DATA=0

VERBOSITY = V_SPEW
ITERATIONS = 10

def run(program, *args):
    output = subprocess.check_output(program, *args, shell=True,
            executable='/bin/bash', env=os.environ)
    return output

def printv(verbosity, *args):
    if VERBOSITY >= verbosity:
        print(*args)

def run_benchmark(benchmark_text, benchmark, size):
    size = str(size)
    printv(V_SPEW, 'Running '+ benchmark_text +' becnhmark with '+ size +'Kb writes and reads')
    output = run('./main -'+benchmark +' -size '+size)
    printv(V_SPEW, output)

    # Return result - first value
    return output.split()[0]

def gather_benchmarks(benchmark_text, benchmark):
    total = 0
    for size in [1,10,100]:
        for j in range(ITERATIONS):
            total += int(run_benchmark(benchmark_text, benchmark, size))
    return total/ITERATIONS

def main():
    printv(V_SPEW,'SAW: ' + str(gather_benchmarks("sequential non-blocking",
        'SAR')))

    printv(V_SPEW,'SAW: ' + str(gather_benchmarks("sequential non-blocking",
        'SAR')))

    gather_benchmarks("sequential blocking", 'SBW')
    gather_benchmarks("sequential blocking", 'SBR')


if __name__ == '__main__':
    main()
