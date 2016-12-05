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
    make()

def build(file_):
    output = command('go version')
    compiler_version = output.strip()
    if compiler_version != MODDED_COMPILER_VERSION:
        fail('Not using the modded compiler!\n'+'Your compiler version is:' + compiler_version)

    if file_ == 'bmain':
        command(' '.join(('go build -o', BLOCKING_BIN, BLOCKING_SRC)))
    elif file_ == 'amain':
        command(' '.join(('go build -o', ASYNC_BIN, ASYNC_SRC)))

def generate(file_):
    if file_ == 'bmain.go':
        generateFile(BLOCKING_SRC, BLOCKING_REPO)
    elif file_ == 'amain.go':
        generateFile(ASYNC_SRC, ASYNC_REPO)

def make():
    command(' '.join(('make -f', joinpath(DIR,'build.mk'), '-C', DIR)))

if __name__ == '__main__':
    build_project()
