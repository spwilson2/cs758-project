package blocking

// blocking IO Scheduler

import (
	"errors"
	"log"
	"os"
	"runtime"
	"syscall"
)

var TestExport int
var initialized = false
var channel chan Operation

var fail = errors.New("")

const CTX_SIZE = 10
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

	Ret_Valid *bool  // true when operation is done, false otherwise.
	Ret_N     *int   // return value, you must specify pointer for it to write to.
	Ret_Err   *error // ^
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
func Open(path string, mode int, perm uint32) (fd int, err error) {
	return syscall.Open(path, mode, perm)
}

/* opens file w/ default params */
func OpenFile(path string) (fd int, err error) {
	mode := syscall.O_RDWR
	perm := uint32(0)
	return syscall.Open(path, mode, perm)
}

/* Writes p to file opened w/  fd */
func Write(fd int, p []byte) (n int, err error) {
	return syscall.Write(fd, p)
}

func WriteAt(fd int, off int, p []byte) (n int, err error) {

	// seek for non-zero offset
	if off != 0 {
		// 0 == SEEK_SET (from beginning)
		syscall.Seek(op.Fd, op.Off, 0)
	}

	return syscall.Write(fd, p)
}

/* Reads from file fd */
func Read(fd int, p []byte) (n int, err error) {
	return syscall.Read(fd, p)
}

/* Read from file, at offset off */
func ReadAt(fd int, off int, p []byte) (n int, err error) {

	if off != 0 {
		// 0 == SEEK_SET (from beginning)
		syscall.Seek(op.Fd, op.Off, 0)
	}

	return syscall.Read(op.Fd, op.Buf)
}

/* Create file */
func Creat(path string, mode uint32) (fd int, err error) {
	return syscall.Creat(path, mode)
}

var channel chan Operation

func InitScheduler(c chan Operation) {

	// scheduler was already initialized, ignore this call
	if initialized != false {
		return
	}

	// set up goroutine for scheduler to run, with the passed channel
	channel = c // global state
	go scheduler(c)
	initialized = true
}
