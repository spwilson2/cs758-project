package nonblocking

// Non-blocking IO Scheduler

import (
	"errors"
	aio "github.com/spwilson2/cs758-project/libaio"
	"log"
	"os"
	"syscall"
)

var TestExport int
var initialized = 0

var fail = errors.New("")

const VALID_STRING string = "package main"
const TESTFILE string = "aio-example.go"

/* struct and const of Operations for AIO to do */
const (
	WRITE    int = 0
	WRITEAT  int = 1
	READ     int = 2
	READAT   int = 3
	CREATE   int = 4
	OPEN     int = 5
	OPENFILE int = 6
)

type Operation struct {
	Op   int    // the operation to do (as per consts above)
	Fd   int    // file descriptor, if doing rd/wr, need this instead
	Name string // name argument for certain ops (open, openfile, create)
	Buf  []byte // input (read, readat) or output (write, writeat) buf.
	Off  int64  // offset, used for READAT, WRITEAT
}

/* extend os.File so we can actually make our own methods */
type File struct {
	os.File
}

/* easily check for errors and panic */
func chk_err(err error) {
	if err != nil {
		log.Printf("(FAILED) %s\n", os.Args[0])
		panic(err) //os.Exit(-1)
	}
}

/* opens file  */
func Open(name string) (*os.File, error) {
	return os.Open(name)
}

/* opens file w/ specified perms and flags*/
func OpenFile(name string, flag int, perm os.FileMode) (*os.File, error) {
	return os.OpenFile(name, flag, perm)
}

/* Writes p to file opened w/  fd */
func (f *File) Write(b []byte) (int, error) {
	return f.Write(b)
}

/* Writes to a specific byte of a file */
func (f *File) WriteAt(b []byte, off int64) (int, error) {
	return f.WriteAt(b, off)
}

/* Reads from file */
func (f *File) Read(b []byte) (int, error) {
	return f.Read(b)
}

/* Read from file, at offset off */
func (f *File) ReadAt(b []byte, off int64) (int, error) {
	return f.ReadAt(b, off)
}

/* Create file */
func Create(name string) (*os.File, error) {
	return os.Create(name)
}

/* Scheduler function, reads op from channel and does it */

func scheduler(c chan Operation) {
	for {
		op := <-c
		log.Println("Received operation: ", op)

		// @TODO: Handle operations.

		var ctx syscall.AioContext_t

		// set up AIO
		chk_err(syscall.IoSetup(1, &ctx)) // 1 == num of AIO in-flight
		defer chk_err(syscall.IoDestroy(ctx))

		var iocb syscall.Iocb
		var iocbp = &iocb

		switch {
		case op.Op == READ:
			log.Println("READing: ", op)
			log.Println(aio.CMD_PREAD) // need to compile for now

			// begin read
			aio.PrepPread(iocbp, op.Fd, op.Buf, uint64(len(op.Buf)), 0)
			chk_err(syscall.IoSubmit(ctx, 1, &iocbp))

			// validate read

			// check to see if we actually got valid results back
			var event syscall.IoEvent
			var timeout syscall.Timespec
			events := syscall.IoGetevents(ctx, 1, 1, &event, &timeout)

			if events <= 0 {
				chk_err(fail)
			}

			if string(op.Buf[:len(VALID_STRING)]) != VALID_STRING {
				log.Printf("Expected %s, found %s\n", VALID_STRING, op.Buf)
				chk_err(fail)
			}

			log.Println("Read succeeded... waiting for next op. ")

		default:
			log.Println("Op not found ", op.Op)
		}

	}
}

/*
   Called once upon creation, stays running and reading from channel for directions on what to do
*/
func InitScheduler(c chan Operation) {

	// scheduler was already initialized, ignore this call
	if initialized != 0 {
		return
	}

	// set up goroutine for scheduler to run, with the passed channel
	go scheduler(c)
	initialized = 1
}
