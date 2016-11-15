package blocking

// Blocking IO Scheduler

import (
	"fmt"
	"log"
	"os"
	"syscall"
)

var TestExport int

/* writes "Hello world!\n" to test.txt */
func Write() {

	// create if doesn't exist, and read/write
	mode := syscall.O_CREAT | syscall.O_RDWR

	// open file ./test.txt
	fd, err := syscall.Open("./test.txt", mode, 0666)

	if err != nil {
		log.Fatal("File failed to open, exiting...")
		os.Exit(1)
	}

	// write "Hello world!" to file
	var p []byte
	p = []byte("Hello world!\n")

	n, err := syscall.Write(fd, p)

	if err != nil {
		log.Fatal("File failed to write, exiting...")
		os.Exit(1)
	}

	fmt.Printf("Bytes written: %d\n", n)

}
