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
	name := "aheh"
	buf := make([]byte, 100)
	fd, _ := syscall.Open("main.go", syscall.O_RDONLY, 0)
	op := nonblocking.Operation{nonblocking.READ, fd, name, buf, 0}
	c := make(chan nonblocking.Operation)
	nonblocking.InitScheduler(c)
	c <- op // send the read operation
	time.Sleep(2 * time.Second)
	c <- op // send another read
	time.Sleep(2 * time.Second)
}
