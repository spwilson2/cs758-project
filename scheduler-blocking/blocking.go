package blocking

// Blocking IO Scheduler

import (
	"log"
	"os"
	"syscall"
)

var TestExport int

/* writes "Hello world!\n" to test.txt */
func do_write(fd int, data []byte, c chan int) {

	// write data to file
	n, err := syscall.Write(fd, data)

	if err != nil || n == 0 {
		log.Fatal("File failed to write, exiting...")
		os.Exit(1)
	}

	// fsync data
	err = syscall.Fsync(fd)

	if err != nil {
		log.Fatal("fsync failed, exiting, fd: ...", fd)
		os.Exit(1)
	}

	//fmt.Printf("Bytes written: %d\n", n)

	close(c) // done with the write
}

/* writes data to fd */
func Write(fd int, data []byte) {
	c := make(chan int)
	go do_write(fd, data, c)

	// wait until channel is closed
	for i := range c {
		log.Println("i: ", i)
	}

}
