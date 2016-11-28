package main

import (
	"log"
	"syscall"
	"time"

	blocking "github.com/spwilson2/cs758-project/scheduler-blocking"
	nonblocking "github.com/spwilson2/cs758-project/scheduler-nonblocking"
)

func main() {

	// remove these once we use the packages, get compiler errs otherwise.
	_ = blocking.TestExport
	_ = nonblocking.TestExport

	// initialize scheduler
	log.Printf("Testing nonblocking")
	c := make(chan nonblocking.Operation)
	nonblocking.InitScheduler(c)

	// send various operations to scheduler
	name := "hello.txt"
	buf := make([]byte, 100)
	off := 0
	fd, _ := nonblocking.Open(name, syscall.O_RDWR, 0)

	// READ
	nonblocking.Read(fd, buf)
	time.Sleep(2 * time.Second)

	// READAT
	nonblocking.ReadAt(fd, off, buf)
	time.Sleep(2 * time.Second)

	// WRITE
	buf = []byte("World Hello")
	nonblocking.Write(fd, buf)
	time.Sleep(2 * time.Second)

	// WRITEAT
	buf = []byte("World")
	off = 5
	nonblocking.WriteAt(fd, off, buf)
	time.Sleep(2 * time.Second)

}
