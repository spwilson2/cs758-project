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

class Test():
    def __init__(self, blocking, **kwargs):
        for flag in kwargs.keys():
            assert flag in TEST_FLAGS

        if blocking:
            self.program = BLOCKING_BIN
        else:
            self.program = ASYNC_BIN

        self.flags = kwargs
        self.flags['path'] = BASE_DIR

    def run(self):
        program = self.program
        args = " "
        args = args.join(('-' + flag + ' ' + val for flag, val in self.flags.items()))
        result = command(program + ' ' + args)

def setupProject():
    build_project()

def main():
    setupProject()
    Test(blocking=True, rsize='1000', nreads='10', nfiles='1').run()
    Test(blocking=False, rsize='1000', nreads='10', nfiles='1').run()

if __name__ == '__main__':
    parse_args()
    main()
