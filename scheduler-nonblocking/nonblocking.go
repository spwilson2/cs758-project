package nonblocking

// Non-blocking IO Scheduler

import (
	aio "github.com/spwilson2/cs758-project/libaio"
	"log"
	"os"
)

var TestExport int
var initialized = 0

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
	Op   int
	Name string // name argument for certain ops
	Buf  []byte // input or output buffer, depends on op
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

		switch {
		case op.Op == READ:
			log.Println("READing: ", op)
			log.Println(aio.CMD_PREAD) // need to compile for now
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
