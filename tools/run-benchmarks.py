#!/usr/bin/python

from __future__ import print_function
from argparse import ArgumentParser
import subprocess
import re
import sys
import os

from util import *
from constants import *
from build import *

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
    def __init__(self, blocking, GOMAXPROCS=None, **kwargs):
        for flag in kwargs.keys():
            assert flag in TEST_FLAGS

        if blocking:
            self.program = BLOCKING_BIN
        else:
            self.program = ASYNC_BIN

        self.flags= {}
        self.flags['path'] = BASE_DIR
        self.flags.update(kwargs)

        self.GOMAXPROCS = None if GOMAXPROCS is None else str(GOMAXPROCS)

    def run(self):
        program = self.program
        env = '' if self.GOMAXPROCS is None else 'GOMAXPROCS='+self.GOMAXPROCS
        args = " "
        args = args.join(('-' + flag + ' ' + val for flag, val in self.flags.items()))
        result = command(' '.join((env, program, args)))

def setupProject():
    build_project()


def main():
    setupProject()
    Test(blocking=True, rsize='1000', nreads='10', nfiles='1').run()
    Test(blocking=False, rsize='1000', nreads='10', nfiles='1').run()
    Test(blocking=False, GOMAXPROCS=3, threads='2', rsize='1000', nreads='10', nfiles='1').run()

if __name__ == '__main__':
    parse_args()
    main()
