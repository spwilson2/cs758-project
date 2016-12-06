package nonblocking

import (
	"log"
	"runtime"
	"sync"
	"syscall"
	"unsafe"

	aio "github.com/spwilson2/cs758-project/libaio"
	tracer "github.com/spwilson2/cs758-project/tracer"
)

const IO_EVENTS_FAIL = Error("Error on IoGetevents call. Expected 1 or more return events.")

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
type s struct {
	init        sync.Once
	initialized bool
	channel     chan operation
}

var scheduler s

/*
   Called once upon creation, stays running and reading from channel for directions on what to do
*/
func InitScheduler(enableTracing bool) {

	do_init := func() {
		initTracer(enableTracing)

		var trace *tracer.TraceEvent
		if enableTracing {
			trace = tracer.NewTraceEvent(T_SCHEDULER_INIT, &schedulerTraceList)
			trace.Start()
			defer trace.Stop()
		}

		// set up goroutine for scheduler to run on another routine
		scheduler.channel = make(chan operation)
		go scheduler.run()
		scheduler.initialized = true
	}

	scheduler.init.Do(do_init)
}

/* Scheduler function, reads op from channel and does it */
func (*s) run() {
	// map to hold inflight aio ops
	inflight := make(map[*syscall.Iocb]operation)
	inflight_ctx := make([]*Context, 100, 1000)

	//var context *Context // context so we can use it inside of the diff select cases

	for {
		select {
		case op := <-scheduler.channel:
			log.Printf("Received operation: %v \n", op.Op)

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
				offset = true
				fallthrough
			case op.Op == READ:

				if offset == false {
					op.Off = 0 // not using offset
				}

				// begin read
				aio.PrepPread(iocbp, op.Fd, op.Buf, uint64(len(op.Buf)), op.Off)
				//chk_err(syscall.IoSubmit(ctx, 1, &iocbp))
				chk_err(ioSubmitWrapper(ctx, 1, &iocbp))

				// save this op as inflight
				inflight_ctx = append(inflight_ctx, context)
				inflight[iocbp] = op

			case op.Op == WRITEAT:
				offset = true
				fallthrough
			case op.Op == WRITE:
				if offset == false {
					op.Off = 0 //do not use offset
				}

				// begin read
				aio.PrepPwrite(iocbp, op.Fd, op.Buf, uint64(len(op.Buf)), op.Off)
				//chk_err(syscall.IoSubmit(ctx, 1, &iocbp))
				chk_err(ioSubmitWrapper(ctx, 1, &iocbp))

				// save op & ctx as inflight
				inflight_ctx = append(inflight_ctx, context)
				inflight[iocbp] = op

			// switch
			default:
			}

		// select: no ops yet, check for any inflight ops to be done
		default:
			//log.Println("No new ops queued, checking in-flight...")

			runtime.Gosched()
			// check all in-flight contexts
			for _, context := range inflight_ctx {
				runtime.Gosched()

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
				op, ok := inflight[(*syscall.Iocb)(unsafe.Pointer(uintptr(event.Obj)))]
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

		} // end of select

	}
}
