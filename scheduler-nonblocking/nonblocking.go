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
	// make sure sched is running
	if initialized != true {
		panic("Scheduler not initialized...")
	}

	log.Println("WRITE requested, sending to scheduler...")

	// create op and send it to scheduler
	op := Operation{WRITE, fd, "", p, 0}
	channel <- op

	//TODO: What do we do here?
	return 1, nil

	// orig:
	//return f.Write(fd, p)
}

func WriteAt(fd int, off int, p []byte) (n int, err error) {
	// make sure sched is running
	if initialized != true {
		panic("Scheduler not initialized...")
	}

	log.Println("WRITEAT requested, sending to scheduler...")

	// create op and send it to scheduler
	op := Operation{WRITEAT, fd, "", p, int64(off)}
	channel <- op

	//TODO: What do we do here?
	return 1, nil

	//return f.Write(fd, p)
}

/* Reads from file */
func Read(fd int, p []byte) (n int, err error) {
	// make sure sched is running
	if initialized != true {
		panic("Scheduler not initialized...")
	}

	log.Println("READ requested, sending to scheduler...")

	// create op and send it to scheduler
	op := Operation{READ, fd, "", p, 0}
	channel <- op

	//TODO: What do we do here?
	return 1, nil

	//return f.Read(fd, p)
}

/* Read from file, at offset off */
func ReadAt(fd int, off int, p []byte) (n int, err error) {
	// make sure sched is running
	if initialized != true {
		panic("Scheduler not initialized...")
	}

	log.Println("READAT requested, sending to scheduler...")

	// create op and send it to scheduler
	op := Operation{READAT, fd, "", p, int64(off)}
	channel <- op

	//TODO: What do we do here?
	return 1, nil

	//return f.Read(fd, p)
}

/* Create file */
func Create(name string) (*os.File, error) {
	return os.Create(name)
}

/* reference counting for context */
type Context struct {
	ctx        syscall.AioContext_t
	references uint
	maxsize    uint
}

var currentCtx Context

/* setup context for aio, manages as reference counter */
func GetCtx(num uint) (*syscall.AioContext_t, error) {
	// check if we have any more slots remaining in this context
	if currentCtx.references+num >= currentCtx.maxsize {
		log.Printf("GetCtx, have %d references so far\n", currentCtx.references)
		currentCtx.references += num
		return &currentCtx.ctx, nil
	} else {
		log.Printf("GetCtx, current ctx at %d capacity. Adding new ctx with %d capacity:\n", currentCtx.maxsize, num)
		var ctx syscall.AioContext_t
		chk_err(syscall.IoSetup(num, &ctx))
		currentCtx.ctx = ctx
		currentCtx.references = 0
		currentCtx.maxsize = num * 2
		return &currentCtx.ctx, nil
	}
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

		var iocb syscall.Iocb
		var iocbp = &iocb

		var offset = false

		switch {
		case op.Op == READAT:
			log.Println("READAT: ", op)
			offset = true
			fallthrough
		case op.Op == READ:
			log.Println("READ: ", op)

			if offset == false {
				op.Off = 0 // not using offset
			}

			// begin read
			aio.PrepPread(iocbp, op.Fd, op.Buf, uint64(len(op.Buf)), op.Off)
			chk_err(syscall.IoSubmit(ctx, 1, &iocbp))
			log.Println("Read submitted, waiting...")

			// check to see if we actually got valid results back
			var event syscall.IoEvent
			var timeout syscall.Timespec
			events := syscall.IoGetevents(ctx, 1, 1, &event, &timeout)

			if events <= 0 {
				chk_err(fail)
			}

			/*
				if offset == false && string(op.Buf[:len(VALID_STRING)]) != VALID_STRING {
					log.Printf("Expected %s, found %s\n", VALID_STRING, op.Buf)
					chk_err(fail)
				}
			*/

			log.Println("Read succeeded... waiting for next op. ")

		case op.Op == WRITEAT:
			log.Println("WRITEAT: ", op)
			offset = true
			fallthrough
		case op.Op == WRITE:
			if offset == false {
				op.Off = 0 //do not use offset
			}

			log.Println("WRITE: ", op)

			// begin read
			aio.PrepPwrite(iocbp, op.Fd, op.Buf, uint64(len(op.Buf)), op.Off)
			chk_err(syscall.IoSubmit(ctx, 1, &iocbp))
			log.Println("Write submitted...")

			// check to see if we actually got valid results back
			var event syscall.IoEvent
			var timeout syscall.Timespec
			events := syscall.IoGetevents(ctx, 1, 1, &event, &timeout)
			if events <= 0 {
				chk_err(fail)
			}

		default:
			log.Println("Op not found ", op.Op)
		}

		chk_err(syscall.IoDestroy(ctx))

	}
}

/*
   Called once upon creation, stays running and reading from channel for directions on what to do
*/
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
