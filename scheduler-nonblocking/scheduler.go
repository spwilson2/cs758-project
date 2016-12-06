package nonblocking

import (
	"log"
	_ "runtime"
	"sync"
	"syscall"
	"unsafe"

	aio "github.com/spwilson2/cs758-project/libaio"
	tracer "github.com/spwilson2/cs758-project/tracer"
)

const IO_EVENTS_FAIL = Error("Error on IoGetevents call. Expected 1 or more return events.")
const UNSUPPORTED_EVENT = Error("Unsupported operation provided!\n")

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

type operation struct {
	Op   int    // the operation to do (as per consts above)
	Fd   int    // file descriptor, if doing rd/wr, need this instead
	Name string // name argument for certain ops (open, openfile, create)
	Buf  []byte // input (read, readat) or output (write, writeat) buf.
	Off  int64  // offset, used for READAT, WRITEAT

	Ret_Valid *bool  // true when operation is done, false otherwise.
	Ret_N     *int   // return value, you must specify pointer for it to write to.
	Ret_Err   *error // ^
}

type queuedOp struct {
	iocbp   *syscall.Iocb
	context *Context
}

type s struct {
	init           sync.Once
	initialized    bool
	channel        chan operation
	inflight_ops   map[*syscall.Iocb]operation
	inflight_ctx   []*Context
	queued_ops     map[*syscall.Iocb]operation
	executeTimeout *timeoutObject
	//geteventsTimeout *timeoutObject
}

var scheduler s

/*
   Called once upon creation, stays running and reading from channel for directions on what to do
*/
func InitScheduler(enableTracing bool) {

	do_init := func() {
		initTracer(enableTracing)

		var trace *tracer.TraceEvent //Ensure trace is scoped to end of function.
		if enableTracing {
			trace = tracer.NewTraceEvent(T_SCHEDULER_INIT, &schedulerTraceList)
			trace.Start()
			defer trace.Stop()
		}

		// set up goroutine for scheduler to run on another routine
		scheduler.channel = make(chan operation)
		scheduler.inflight_ops = make(map[*syscall.Iocb]operation)
		scheduler.inflight_ctx = make([]*Context, 100, 1000)

		//scheduler.geteventsTimeout = newTimeoutObject(GETEVENTS_TIMEOUT)
		scheduler.executeTimeout = newTimeoutObject(EXECUTE_TIMEOUT)

		//go scheduler.geteventsTimeout.begin()
		go scheduler.executeTimeout.begin()

		go runScheduler()
		scheduler.initialized = true
	}

	scheduler.init.Do(do_init)
}

/*
* Scheduler main function.
*
* The scheduler Has three jobs - Jobs are numbered by their priority, but
* listed in a way that makes most sense to a reader:
*
* 3) Must read and collect op requests from the channel
*
* 2) Must execute submit these ops every executeOpTimeout
*
* 1) Must check for completion of ops every checkCompletionTimeout
*
* Due to time constraints we haven't implemented a progress garuntee, so we
* have decided to keep the order the statements are in, not the labled order.
*
 */
func runScheduler() {

	// Put us on our own thread to avoid thrashing from syscalls.
	//runtime.LockOSThread()

	for {
		select {
		case op := <-scheduler.channel:
			// Queue the Op to be executed.
			queueOp(op)

		case <-scheduler.executeTimeout.signal:
			// If there are any ops waiting to be submitted, submit
			// them now.
			if len(scheduler.queued_ops) >= 1 {
				submitQueue()
			}

		default:
			// Check if any ops have completed, if so return their
			// results.
			getEvents()
		}
	}
}

func submitQueue() {
}

func queueOp(op operation) {
}

func getEvents() {
	for _, context := range scheduler.inflight_ctx {
		//runtime.Gosched()

		if context == nil {
			continue
		}

		var event syscall.IoEvent
		var timeout syscall.Timespec
		ctx := context.ctx

		//events := syscall.IoGetevents(ctx, 1, 1, &event, &timeout)
		events := ioGeteventsWrapper(ctx, 1, 1, &event, &timeout)

		if events <= 0 {
			continue // not done yet
		}

		// set return vals, obtained from a map
		op, ok := scheduler.inflight_ops[(*syscall.Iocb)(unsafe.Pointer(uintptr(event.Obj)))]
		// make sure this event existed. If not, ???
		if ok == false {
			log.Println("event did not exist ??")
			chk_err(IO_EVENTS_FAIL)
		}

		//@TODO: FIX THIS
		//N_ptr := (*int)(unsafe.Pointer(uintptr(event.Obj)))
		//log.Println("Setting return vals: ", *event.Obj)
		*(op.Ret_N) = int(event.Res)
		*(op.Ret_Err) = nil
		*(op.Ret_Valid) = true

		//@TODO: Verify.
		ungetCtx(context, 1)

	}
}

func newQueuedOp(op operation) (newOp *queuedOp) {
	var ctx syscall.AioContext_t
	context, err := getCtx(1)
	chk_err(err)
	ctx = context.ctx
	_ = ctx

	var iocb syscall.Iocb
	var iocbp = &iocb

	newOp.iocbp = iocbp
	newOp.context = context

	switch {
	case op.Op == READ:
		op.Off = 0
		fallthrough
	case op.Op == READAT:
		// begin read
		aio.PrepPread(iocbp, op.Fd, op.Buf, uint64(len(op.Buf)), op.Off)

	case op.Op == WRITE:
		op.Off = 0
		fallthrough
	case op.Op == WRITEAT:
		aio.PrepPwrite(iocbp, op.Fd, op.Buf, uint64(len(op.Buf)), op.Off)
	default:
		chk_err(UNSUPPORTED_EVENT)
	}
	return
}

func enqueueOp(op operation) {
}

//func execOp(op operation) {
//	//log.Printf("Received operation: %v \n", op.Op)
//
//	// set up AIO
//	var ctx syscall.AioContext_t
//	context, err := getCtx(1)
//	chk_err(err)
//	ctx = context.ctx
//
//	var iocb syscall.Iocb
//	var iocbp = &iocb
//
//	var offset = false
//
//	switch {
//	case op.Op == READAT:
//		offset = true
//		fallthrough
//	case op.Op == READ:
//
//		if offset == false {
//			op.Off = 0 // not using offset
//		}
//
//		// begin read
//		aio.PrepPread(iocbp, op.Fd, op.Buf, uint64(len(op.Buf)), op.Off)
//		//chk_err(syscall.IoSubmit(ctx, 1, &iocbp))
//		chk_err(ioSubmitWrapper(ctx, 1, &iocbp))
//
//		// save this op as inflight
//		scheduler.inflight_ctx = append(scheduler.inflight_ctx, context)
//		scheduler.inflight_ops[iocbp] = op
//
//	case op.Op == WRITEAT:
//		offset = true
//		fallthrough
//	case op.Op == WRITE:
//		if offset == false {
//			op.Off = 0 //do not use offset
//		}
//
//		// begin read
//		aio.PrepPwrite(iocbp, op.Fd, op.Buf, uint64(len(op.Buf)), op.Off)
//		//chk_err(syscall.IoSubmit(ctx, 1, &iocbp))
//		chk_err(ioSubmitWrapper(ctx, 1, &iocbp))
//
//		// save op & ctx as inflight
//		scheduler.inflight_ctx = append(scheduler.inflight_ctx, context)
//		scheduler.inflight_ops[iocbp] = op
//	}
//}
