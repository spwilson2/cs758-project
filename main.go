package main

import (
	"log"
	"os"
	"syscall"

	blocking "github.com/spwilson2/cs758-project/scheduler-blocking"
	nonblocking "github.com/spwilson2/cs758-project/scheduler-nonblocking"
)

func blocking_write(data []byte) {
	// create if doesn't exist, and read/write
	mode := syscall.O_CREAT | syscall.O_RDWR

	// open file ./test.txt
	fd, err := syscall.Open("./test.txt", mode, 0666)

	if err != nil {
		log.Fatal("File failed to open, exiting...")
		os.Exit(1)
	}

	// do our blocking write
	blocking.Write(fd, data)

}

func main() {

	// remove these once we use the packages, get compiler errs otherwise.
	_ = blocking.TestExport
	_ = nonblocking.TestExport

	// write "Hello world!" to file
	var d []byte
	d = []byte("Hello world!\n")

	blocking_write(d)
}
