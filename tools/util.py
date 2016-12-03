from __future__ import print_function
from argparse import ArgumentParser
import os
import sys
import subprocess

##############################################################################
## Flags                                                                    ##
##############################################################################
V_SPEW  = 3
V_TRACE = 2
V_INFO  = 1
V_DATA  = 0
VERBOSITY = None
##############################################################################

##############################################################################
## Utility Functions                                                        ##
##############################################################################
joinpath = os.path.join
def parse_args():

    parser = ArgumentParser()
    parser.add_argument("-v", "--verbosity",
                        action="count", default=0,
                        help="increase output verbosity")
    flags = parser.parse_args()

    global VERBOSITY
    VERBOSITY = flags.verbosity

    printv(V_SPEW, 'Input Flags: ' + str(flags.__dict__))


def command(program, *args):
    printv(V_SPEW, 'Running: ' + program)
    output = subprocess.check_output(program, *args, shell=True,
            executable='/bin/bash', env=os.environ)
    if type(output) is bytes:
        output = output.decode('utf-8')
    printv(V_SPEW, 'Output: ' + output)

    return output

def printv(verbosity, *args):
    if VERBOSITY >= verbosity:
        print(*args)

def fail(*args):
    print(*args)
    sys.exit(-1)
##############################################################################
