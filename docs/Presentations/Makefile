CC=gcc
CCFLAGS+=-laio -ggdb
TARGETS=aio-example io-example odirect-aio
.PHONY:all clean
all: $(TARGETS)
%:%.c Makefile
	$(CC) -o $@ $< $(CCFLAGS) 
clean:
	rm -rf $(TARGETS)
