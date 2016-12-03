#!/usr/bin/python

from util import *
from constants import *

def generateFile(output_file, scheduler_repo):
    main = open(MAIN_SRC, 'r')
    genmain = open(output_file, 'w')

    for line in main.readlines():
        line = REPO_REGEX.sub(SCHEDULER_IMPORT_NAME + ' ' + scheduler_repo, line)
        genmain.write(line)

    main.close()
    genmain.close()

def build_project():
    output = command('go version')
    compiler_version = output.strip()
    if compiler_version != MODDED_COMPILER_VERSION:
        fail('Not using the modded compiler!\n'+'Your compiler version is:' + compiler_version)

    # Generate both copies of main for each scheduler type.
    for args in ((ASYNC_SRC, ASYNC_REPO), (BLOCKING_SRC, BLOCKING_REPO)):
        generateFile(*args)

    command(' '.join(('go build -o', ASYNC_BIN, ASYNC_SRC)))
    command(' '.join(('go build -o', BLOCKING_BIN, BLOCKING_SRC)))

def make():
    command(' '.join(('make -f', joinpath(DIR,'build.mk'), '-C', DIR)))

if __name__ == '__main__':
    build_project()
