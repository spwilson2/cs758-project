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
        args = args.join(('-' + str(flag) + ' ' + str(val) for flag, val in self.flags.items()))
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
                result_dict[key] = val

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
    def saveResults(results, name, sortParameter):
        plot.save_csv(results, joinpath(CSV_DIR, name + '-results.csv'))
        plot.bar(results, file_=joinpath(PLOT_DIR, name + '-results.png'), sortParameter=sortParameter, title=name, ylab="execution time (ns)")

def setupProject():
    make()


def executeTestAndReturnResults(blocking, opType, offset, size, threadCount, numops, nfiles):
    ioType = "blocking" if blocking else "nonblocking"
    test = None
    if(opType == "reads"):
        test = Test(blocking, (ioType + " " + opType + "(offset: " + str(offset) + ", size: " + str(size) + ", threads: " + str(threadCount) + ")"), roff=offset, rsize=size, nreads=numops, nfiles = 1, threads=threadCount, GOMAXPROCS=threadCount)
    else:
        test = Test(blocking, (ioType + " " + opType + "(offset: " + str(offset) + ", size: " + str(size) + ", threads: " + str(threadCount) + ")"), woff=offset, wsize=size, nwrites=numops, nfiles = 1, threads=threadCount, GOMAXPROCS=threadCount)
    return test.getResults()

def main():
    setupProject()

    for offset in range(0, 1):
        for exponent in range(3, 7):
            readSize = 10 ** exponent
            for threadCount in [1, 2, 4, 8, 16]:
                blockingResults = executeTestAndReturnResults(True, "reads", offset, threadCount, readSize, 10, 1)         
                nonblockingResults = executeTestAndReturnResults(False, "reads", offset, threadCount, readSize, 10, 1)

                for result in blockingResults:
                    nonblockingResults.append(result)
                Test.saveResults(nonblockingResults, ("reads(offset: " + str(offset) + ", size: " + str(readSize) + ", threads: " + str(threadCount) + ")"), Go.IO_TYPE_KEY)

    for offset in range(0, 1):
        for exponent in range(3, 7):
            writeSize = 10 ** exponent
            for threadCount in [1, 2, 4, 8, 16]:
                blockingResults = executeTestAndReturnResults(True, "writes", offset, threadCount, writeSize, 10, 1)
                nonblockingResults =  executeTestAndReturnResults(False, "writes", offset, threadCount, readSize, 10, 1)
                
                for result in blockingResults:
                    nonblockingResults.append(result)
                Test.saveResults(nonblockingResults, ("writes(offset: " + str(offset) + ", size: " + str(writeSize) + ", threads: " + str(threadCount) + ")"), Go.IO_TYPE_KEY)
                

if __name__ == '__main__':
    parse_args()
    main()