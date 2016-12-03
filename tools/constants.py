import re
import sys
import os
_joinpath = os.path.join

DIR = os.path.realpath(__file__)
BASE_DIR = os.path.dirname(os.path.realpath(_joinpath(DIR, '..')))
GENDIR = _joinpath(BASE_DIR, 'generated')
MODDED_COMPILER_VERSION = 'go version go1.7.3 linux/amd64'
SRC_DIR = 'src'

##############################################################################
## File Generation Constants                                                ##
##############################################################################
MAIN_SRC = _joinpath(BASE_DIR, SRC_DIR, 'main.go')
REPO_REGEX = re.compile("SCHEDULER_UNDER_TEST")
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

TEST_FLAGS = {
        "t",
        "rsize",
        "wsize",
        "nwrites",
        "nreads",
        "roff",
        "woff",
        "nfiles",
        "mixed"
        }
