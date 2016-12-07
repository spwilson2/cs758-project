CC=gcc
CCFLAGS+=-laio -ggdb
TARGETS= $(basename $(wildcard *.c))
.PHONY:all clean
all: $(TARGETS)
%:%.c Makefile
	$(CC) -o $@ $< $(CCFLAGS) 
clean:
	rm -rf $(TARGETS)
