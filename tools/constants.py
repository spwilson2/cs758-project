import re
import sys
import os
from util import joinpath as _joinpath

DIR = os.path.dirname(os.path.realpath(__file__))
BASE_DIR = os.path.realpath(_joinpath(DIR, '..'))
GENDIR = _joinpath(BASE_DIR, 'generated')
MODDED_COMPILER_VERSION = 'go version go1.7.3 linux/amd64'
SRC_DIR = _joinpath(BASE_DIR, 'src')
RESULT_DIR = _joinpath(BASE_DIR, 'results')
CSV_DIR = _joinpath(RESULT_DIR, 'csv')
PLOT_DIR = _joinpath(RESULT_DIR, 'plots')

##############################################################################
## File Generation Constants                                                ##
##############################################################################
MAIN_SRC = _joinpath(SRC_DIR, 'main.go')
REPO_REGEX = re.compile("//SCHEDULER_UNDER_TEST")
SCHEDULER_IMPORT_NAME = 'sut'

ASYNC_SRC = _joinpath(GENDIR, 'amain.go')
ASYNC_BIN = _joinpath(GENDIR, 'amain')
ASYNC_REPO = '"github.com/spwilson2/cs758-project/scheduler-nonblocking"'

BLOCKING_SRC = _joinpath(GENDIR, 'bmain.go')
BLOCKING_BIN = _joinpath(GENDIR, 'bmain')
BLOCKING_REPO = '"github.com/spwilson2/cs758-project/scheduler-blocking"'

GENERATED_WARING = '''\
/*
* This file has been automatically genererated using the following
* command %s
*
* Do not modify!
*/''' % ' '.join(sys.argv)
##############################################################################

##############################################################################
## Go Program Constants                                                     ##
##############################################################################
class Go(object):
    TEST_FLAGS = {
            "threads",
            "rsize",
            "wsize",
            "nwrites",
            "nreads",
            "roff",
            "woff",
            "nfiles",
            "mixed"
            }

    IO_TYPE_KEY = "IoType"
    OP_KEY = "Operation"
    LENGTH_KEY = "Length"
##############################################################################
