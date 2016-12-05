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
        args = args.join(('-' + flag + ' ' + val for flag, val in self.flags.items()))
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

    def saveResults(self, results, chartName, sortParameter):
        result_filename = genFilename(pfx=self.name+'-' if self.name is not None else '')
        plot.save_csv(results, joinpath(CSV_DIR, result_filename + '-results.csv'))
        plot.bar(results, file_=joinpath(PLOT_DIR, result_filename + '-results.png'), sortParameter=sortParameter, title=chartName, ylab="execution time (ns)")

def setupProject():
    make()

def createAndRunTest(testName, blocking, rsize, nreads, nfiles, nwrites, wsize, sortParameter, gomaxprocs=None):
    test = Test(blocking=blocking, name=testName, GOMAXPROCS=gomaxprocs, rsize=rsize, nreads=nreads, nfiles=nfiles, nwrites=nwrites, wsize=wsize)
    test.saveResults(test.getResults(), testName, sortParameter)

def main():
    setupProject()

    #createAndRunTest("In Order Mixed nonblocking", False, '1000', '10', '1', '10', '1000', Go.OP_KEY)
    createAndRunTest("In Order Mixed blocking", True, '1000', '10', '1', '10', '1000', Go.OP_KEY)

if __name__ == '__main__':
    parse_args()
    main()