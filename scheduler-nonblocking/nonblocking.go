package nonblocking

// Non-blocking IO Scheduler

import (
	"errors"
	"log"
	"os"
	"syscall"

	aio "github.com/spwilson2/cs758-project/libaio"
)

var TestExport int
var initialized = false
var channel chan Operation

var fail = errors.New("")

const CTX_SIZE = 10
const VALID_STRING string = "package main"
const TESTFILE string = "aio-example.go"

var BAD_CONTEXT_REQUEST = errors.New("Too large of a context was requested.")

const AIO_CONTEXT_MAX = 500000
const AIO_CONTEXT_MIN = 1000
const MIN_LOWER_RATE = 10
const MIN_LOWER_LIMIT = 10
const CONTEXT_REQUEST_MULTIPLIER = 2

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
	// make sure sched is running
	if initialized != true {
		panic("Scheduler not initialized...")
	}

	//log.Println("WRITE requested, sending to scheduler...")

	// create op and send it to scheduler
	var ret_valid = new(bool)
	var ret_n = new(int)
	var ret_err = new(error)

	op := Operation{WRITE, fd, "", p, 0, ret_valid, ret_n, ret_err}
	channel <- op

	for *ret_valid != true {
		// spin until ret_err is valid
	}

	// results now valid, we can return them.
	return *ret_n, *ret_err

	// orig:
	//return f.Write(fd, p)
}

func WriteAt(fd int, off int, p []byte) (n int, err error) {
	// make sure sched is running
	if initialized != true {
		panic("Scheduler not initialized...")
	}

	//log.Println("WRITEAT requested, sending to scheduler...")

	// create op and send it to scheduler
	var ret_valid = new(bool)
	var ret_n = new(int)
	var ret_err = new(error)

	op := Operation{WRITEAT, fd, "", p, int64(off), ret_valid, ret_n, ret_err}
	channel <- op

	for *ret_valid != true {
		// spin until ret_err is valid
	}

	// results now valid, we can return them.
	return *ret_n, *ret_err

	//return f.Write(fd, p)
}

/* Reads from file */
func Read(fd int, p []byte) (n int, err error) {
	// make sure sched is running
	if initialized != true {
		panic("Scheduler not initialized...")
	}

	//log.Println("READ requested, sending to scheduler...")

	// create op and send it to scheduler
	var ret_valid = new(bool)
	var ret_n = new(int)
	var ret_err = new(error)

	op := Operation{READ, fd, "", p, 0, ret_valid, ret_n, ret_err}
	channel <- op

	for *ret_valid != true {
		// spin until ret_err is valid
	}

	// results now valid, we can return them.
	return *ret_n, *ret_err

	//return f.Read(fd, p)
}

/* Read from file, at offset off */
func ReadAt(fd int, off int, p []byte) (n int, err error) {
	// make sure sched is running
	if initialized != true {
		panic("Scheduler not initialized...")
	}

	//log.Println("READAT requested, sending to scheduler...")

	// create op and send it to scheduler
	var ret_valid = new(bool)
	var ret_n = new(int)
	var ret_err = new(error)

	op := Operation{READAT, fd, "", p, int64(off), ret_valid, ret_n, ret_err}
	channel <- op

	for *ret_valid != true {
		// spin until ret_err is valid
	}

	// results now valid, we can return them.
	return *ret_n, *ret_err

	//return f.Read(fd, p)
}

/* Create file */
func Creat(path string, mode uint32) (fd int, err error) {
	return syscall.Creat(path, mode)
}

/* reference counting for context */
type Context struct {
	ctx        syscall.AioContext_t
	references int
	maxsize    int
}

var aio_contexts []*Context
var aio_context_max = AIO_CONTEXT_MAX
var aio_context_min = AIO_CONTEXT_MIN

/* setup context for aio, manages as reference counter */
func getCtx(num int) (*Context, error) {

	// TODO: Order the list by remaining references
	for _, context := range aio_contexts {
		if new_references := context.references + num; new_references <= context.maxsize {
			//log.Printf("GetCtx, have %d references so far\n", context.references)
			context.references = new_references
			return context, nil
		}

	}

	// Unable to find a context with remaining space, let's create a new
	// one.

	var num_context_request int

	switch {
	case num < aio_context_min:
		num_context_request = aio_context_min
	case (num > aio_context_min) && (num < aio_context_max):
		num_context_request = (CONTEXT_REQUEST_MULTIPLIER * num)
		if num_context_request > aio_context_max {
			num_context_request = aio_context_max
		}
	default:
		//log.Printf("Bad number of contexts requested %d.\n", num)
		return nil, BAD_CONTEXT_REQUEST
	}

	var ctx syscall.AioContext_t
	err := syscall.IoSetup(uint(num_context_request), &ctx)

	var new_context *Context = new(Context)

	if err == nil {
		new_context.ctx = ctx
		new_context.references = num
		new_context.maxsize = num_context_request
		aio_contexts = append(aio_contexts, new_context)
		//log.Printf("Appending the new context %p to list\n", new_context)
	} else {
		// Assume the request failed due to not enough resources
		// TODO: should check the error condition...
		// Due to there not being enough resources for a request of
		// this size, assume the Max(cur_max, (num requested - 1)) is
		// actual Max.
		//log.Printf("Failed to request %d from IoSetup\n", num_context_request)

		if aio_context_max < num_context_request {
			aio_context_max = num_context_request
		}

		// If the request couldn't complete and we're requesting the
		// bare minimun, move our min lower (limited to MIN_LOWER_LIMIT)
		if aio_context_min == num_context_request {
			aio_context_min -= (aio_context_min / MIN_LOWER_RATE)
			if aio_context_min <= (MIN_LOWER_LIMIT - 1) {
				aio_context_min = MIN_LOWER_LIMIT
			}
		}

		new_context = nil
	}

	//TODO: Need to add cleanup.

	return new_context, err
}

func ungetCtx(context *Context, num int) error {
	context.references -= num
	return nil
}

/* Scheduler function, reads op from channel and does it */
func scheduler(c chan Operation) {
	for {
		op := <-c
		//log.Println("Received operation: ", op.Op)

		// @TODO: Handle operations.

		// set up AIO
		var ctx syscall.AioContext_t
		context, err := getCtx(1)
		chk_err(err)
		ctx = context.ctx

		var iocb syscall.Iocb
		var iocbp = &iocb

		var offset = false

		switch {
		case op.Op == READAT:
			//log.Println("READAT: ", op)
			offset = true
			fallthrough
		case op.Op == READ:
			//log.Println("READ: ", op)

			if offset == false {
				op.Off = 0 // not using offset
			}

			// begin read
			aio.PrepPread(iocbp, op.Fd, op.Buf, uint64(len(op.Buf)), op.Off)
			chk_err(syscall.IoSubmit(ctx, 1, &iocbp))
			//log.Println("Read submitted, waiting...")

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

			// set return vals
			*(op.Ret_N) = events
			*(op.Ret_Err) = nil
			*(op.Ret_Valid) = true

			//log.Println("Read succeeded... waiting for next op. ")

		case op.Op == WRITEAT:
			//log.Println("WRITEAT: ", op)
			offset = true
			fallthrough
		case op.Op == WRITE:
			if offset == false {
				op.Off = 0 //do not use offset
			}

			//log.Println("WRITE: ", op)

			// begin read
			aio.PrepPwrite(iocbp, op.Fd, op.Buf, uint64(len(op.Buf)), op.Off)
			chk_err(syscall.IoSubmit(ctx, 1, &iocbp))
			//log.Println("Write submitted...")

			// check to see if we actually got valid results back
			var event syscall.IoEvent
			var timeout syscall.Timespec
			events := syscall.IoGetevents(ctx, 1, 1, &event, &timeout)
			if events <= 0 {
				chk_err(fail)
			}

			// set return vals
			*(op.Ret_N) = events
			*(op.Ret_Err) = nil
			*(op.Ret_Valid) = true

			//log.Println("Write successful in scheduler")

		default:
			//log.Println("Op not found ", op.Op)
		}

		ungetCtx(context, 1)
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
