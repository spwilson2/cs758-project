#!/usr/bin/python

from __future__ import print_function
from argparse import ArgumentParser
import subprocess
import re
import sys
import os

from util import *
from constants import *
from build import make
import plot

# Check constants.py for the most up to date version!
#TEST_FLAGS = {
#        "threads",
#        "rsize",
#        "wsize",
#        "nwrites",
#        "nreads",
#        "roff",
#        "woff",
#        "nfiles",
#        "mixed"
#        }

class Test():
    # Expect output to be 'key: value' pairs
    RESULT_REGEX = re.compile('(?P<key>\w+): (?P<val>\w+)')

    def __init__(self, blocking, name=None, GOMAXPROCS=None, **kwargs):
        for flag in kwargs.keys():
            assert flag in Go.TEST_FLAGS

        if blocking:
            self.program = BLOCKING_BIN
        else:
            self.program = ASYNC_BIN

        self.flags= {}
        self.flags['path'] = BASE_DIR
        self.flags.update(kwargs)

        self.GOMAXPROCS = None if GOMAXPROCS is None else str(GOMAXPROCS)
        self.name = None if name is None else str(name).replace(' ', '_')

    def run(self):
        program = self.program
        env = '' if self.GOMAXPROCS is None else 'GOMAXPROCS='+self.GOMAXPROCS
        args = " "
        args = args.join(('-' + flag + ' ' + str(val) for flag, val in self.flags.items()))
        result = command(' '.join((env, program, args)))
        return result

    @staticmethod
    def parseOutput(result):

        results = []
        for line in result.splitlines():
            matches = Test.RESULT_REGEX.finditer(line)
            result_dict = {}
            for match in matches:
                key, val = match.group('key'), match.group('val')
                result_dict[key.encode("ascii")] = val.encode()

            if result_dict:
                results.append(result_dict)

        return results

    def getResults(self):
        results = self.run()
        results = Test.parseOutput(results)

        # TODO: Remove once output from test is more verbose.
        for result in results:
            result.update(self.flags)
            result[Go.IO_TYPE_KEY] = 'blocking' if self.program == BLOCKING_BIN else 'nonblocking'
            result['GOMAXPROCS'] = self.GOMAXPROCS

        return results

    @staticmethod
    def saveResults(results, name):
        plot.save_csv(results, joinpath(CSV_DIR, name + '-results.csv'))
        plot.execution_bar(results, file_=joinpath(PLOT_DIR, name + 'execution-time-results.png'), title=name+" (Execution Time)", ylab="execution time (ns)")
        plot.tracedata_bar(results, file_=joinpath(PLOT_DIR, name + '-trace-events-results.png'), title=name+" (Trace Events)", ylab="execution time (ns)")

def setupProject():
    make()

def batch_readtest(rsize, nreads, nfiles, threads):

    for blocking in [False, True]:
        name = '-readtest'
        if not blocking:
            name = 'aio'+name
        else:
            name = 'blocking'+name

        test = Test(name=name, blocking=blocking, rsize=rsize, nreads=nreads, nfiles=nfiles, threads=threads)
        results = test.getResults()
        iosubmits = len([result for result in results if result[Go.OP_KEY] ==
                'IoSubmit'])
        print("Operations:", iosubmits)
        if results:
            test.saveResults(results, "batch_readtest")

def executeTestAndReturnResults(blocking, opType, offset, size, threadCount, numops, nfiles):
    ioType = "blocking" if blocking else "nonblocking"
    test = None
    if(opType == "reads"):
        test = Test(blocking, (ioType + " " + opType + "(offset: " + str(offset) + ", threads: " + str(threadCount) + ")"), roff=offset, rsize=size, nreads=numops, nfiles = 1, threads=threadCount, GOMAXPROCS=threadCount)
    else:
        test = Test(blocking, (ioType + " " + opType + "(offset: " + str(offset) + ", threads: " + str(threadCount) + ")"), woff=offset, wsize=size, nwrites=numops, nfiles = 1, threads=threadCount, GOMAXPROCS=threadCount)
    return test.getResults()

def main():
    setupProject()

    """
        Executes tests of various thread counts for various read/write sizes
    """
    for testType in ["reads", "writes"]:
        for offset in range(0, 1):
            results = []
            for threadCount in [1,2,4,8,16]:
                for exponent in range(3,7):
                    opSize = 10 ** exponent
                    blockingResults = executeTestAndReturnResults(True, testType, offset, opSize, threadCount, 10, 1)
                    nonblockingResults = executeTestAndReturnResults(False, testType, offset, opSize, threadCount, 10, 1)

                    for result in blockingResults:
                        results.append(result)
                    
                    for result in nonblockingResults:
                        results.append(result)
                    
                Test.saveResults(results, (testType + "(offset: " + str(offset) + ", threads: " + str(threadCount) + ")"))
                
    rsize=1000000
    nreads=20
    nfiles=4
    threads=8
    batch_readtest(rsize,nreads,nfiles,threads)


if __name__ == '__main__':
    parse_args()
    main()