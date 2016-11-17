CC=gcc
CCFLAGS+=-lrt
TARGETS=aio-example io-example
.PHONY:all clean
all: $(TARGETS)
%:%.c Makefile
	$(CC) -o $@ $< $(CCFLAGS) 
clean:
	rm -rf $(TARGETS)
