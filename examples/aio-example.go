package main

// package main must be the first thing in this program for the test to work.

import "fmt"
import "unsafe"
import "os"
import "errors"
import "syscall"
import aio "github.com/spwilson2/cs758-project/libaio"

const TESTFILE string = "aio-example.go"

//const TESTFILE string = "Hello.txt"

var fail = errors.New("")

//const BUFSIZE = 5000000
const BUFSIZE = 512

func chk_err(err error) {
	if err != nil {
		fmt.Printf("(FAILED) %s\n", os.Args[0])
		panic(err) //os.Exit(-1)
	}
}

func main() {

	var ctx syscall.AioContext_t
	var err error

	chk_err(syscall.IoSetup(128, &ctx))
	defer func() {
		chk_err(syscall.IoDestroy(ctx))
	}()

	fd, err := syscall.Open(TESTFILE, syscall.O_RDONLY|syscall.O_DIRECT, 0)
	chk_err(err)

	var iocb syscall.Iocb
	var iocbp = &iocb

	var buffer = make([]byte, BUFSIZE, BUFSIZE)

	if (uintptr(unsafe.Pointer(&buffer[0])) & 0x01ff) > 0 {
		fmt.Printf("Buffer must be 512 byte-aligned..\n")
		fmt.Printf("0x%x\n", &buffer[0])
		chk_err(fail)
	}

	aio.PrepPread(iocbp, fd, buffer, len(buffer), 0)

	// Submit our request.
	chk_err(syscall.IoSubmit(ctx, 1, &iocbp))

	var event syscall.IoEvent
	var timeout syscall.Timespec
	timeout.Sec = 1

	events := syscall.IoGetevents(ctx, 1, 1, &event, &timeout)
	if events == 0 {
		chk_err(fail)
	} else if events < 0 {
		chk_err(fail)
	}

	// Check the result..
	valid_string := "package main"

	if string(buffer[:len(valid_string)]) != valid_string {
		fmt.Printf("Expected:%s, Found:%s\n", valid_string, buffer)
		chk_err(fail)
	}

	fmt.Printf("(OK) %s\n", os.Args[0])
}
