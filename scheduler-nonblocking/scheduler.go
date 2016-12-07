package nonblocking

import (
	"runtime"
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

type queueOp struct {
	iocbp   *syscall.Iocb
	context *Context
	op      *operation
}

type s struct {
	init             sync.Once
	initialized      bool
	disable          chan bool
	channel          chan operation
	inflight_ops     map[*syscall.Iocb]operation
	inflight_context map[*Context][]*queueOp
	queued_ops       map[*Context][]*queueOp
	executeTimeout   *timeoutObject
	//geteventsTimeout *timeoutObject
}

var scheduler s

/*
   Called once upon creation, the scheduler stays running and reads from channel
   for directions on what to do
*/
func InitScheduler(enableTracing bool) {

	do_init := func() {
		initTracer(enableTracing)

		var trace *tracer.TraceEvent //Ensure trace is scoped to end of function.
		if enableTracing {
			trace = tracer.NewTraceEvent(T_SCHEDULER_INIT, schedulerTraceList)
			trace.Start()
			defer trace.Stop()
		}

		// set up goroutine for scheduler to run on another routine
		scheduler.channel = make(chan operation)
		scheduler.disable = make(chan bool)
		scheduler.inflight_ops = make(map[*syscall.Iocb]operation)
		scheduler.inflight_context = make(map[*Context][]*queueOp)
		scheduler.queued_ops = make(map[*Context][]*queueOp)

		//scheduler.geteventsTimeout = newTimeoutObject(GETEVENTS_TIMEOUT)
		scheduler.executeTimeout = newTimeoutObject(EXECUTE_TIMEOUT)

		//go scheduler.geteventsTimeout.begin()
		go scheduler.executeTimeout.begin()

		go runScheduler()
		scheduler.initialized = true
	}

	scheduler.init.Do(do_init)
}

func EndScheduler() {
	end := func() {
		scheduler.disable <- true
	}
	scheduler.init.Do(end)
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
	runtime.LockOSThread()

	for {
		select {
		case <-scheduler.disable:
			return
		case op := <-scheduler.channel:
			// Queue the Op to be executed.
			qop_p := newQueueOp(op)
			enqueueOp(qop_p)

		case <-scheduler.executeTimeout.signal:
			// If there are any ops waiting to be submitted, submit
			// them now.
			submitQueue()

		default:
			// Check if any ops have completed, if so return their
			// results.
			getEvents()
		}
	}
}

/* submitQueue submits all queueOps saved by the scheduler */
func submitQueue() {
	if len(scheduler.queued_ops) == 0 {
		return
	}

	for context, op_queue := range scheduler.queued_ops {
		// Combine the ops for the same context into
		// a single submission.
		submitOps(context, op_queue)

		// Clear entries once they are submitted.
		delete(scheduler.queued_ops, context)
	}
}

/* enqueueOp equeues a queueOp for later submission with syscall.IoSubmit */
func enqueueOp(qop_p *queueOp) {
	context := qop_p.context
	op_list, exists := scheduler.queued_ops[context]
	if !exists {
		op_list = make([]*queueOp, 0)
	}
	op_list = append(op_list, qop_p)
	scheduler.queued_ops[context] = op_list
}

/* Check all inflight events for completion. */
func getEvents() {
	if len(scheduler.inflight_context) == 0 {
		runtime.Gosched()
		return
	}
	for context, qop_p_list := range scheduler.inflight_context {

		var timeout syscall.Timespec

		// TODO: Change max number events (context.maxsize) to the
		// current number of inflight events for the context.
		events, err := ioGeteventsWrapper(context.ctx, 0, len(qop_p_list), &context.event_list[0], &timeout)
		chk_err(err)

		if events < 0 {
			//chk_err(IO_EVENTS_FAIL)
			chk_err(err)
			chk_err(FAIL)
			continue // not done yet
		}

		if events == 0 {
			//chk_err(IO_EVENTS_FAIL)
			continue // not done yet
		}

		// For each returned event, return the vals and finish the op.
		for i := 0; i < events; i++ {

			event := context.event_list[i]

			// set return vals, obtained from a map
			op, _ := scheduler.inflight_ops[(*syscall.Iocb)(unsafe.Pointer(uintptr(event.Obj)))]

			//@TODO: FIX THIS
			//N_ptr := (*int)(unsafe.Pointer(uintptr(event.Obj)))
			//log.Println("Setting return vals: ", *event.Obj)

			*(op.Ret_N) = int(event.Res)
			*(op.Ret_Err) = nil
			*(op.Ret_Valid) = true

		}

		ungetCtx(context, events)
		delete(scheduler.inflight_context, context)
	}
}

/*
* newQueueOp creates and initializes a new queueOp, must be thread safe to be
* able to allow op callers to call this.
 */
func newQueueOp(op operation) (newOp *queueOp) {
	var ctx syscall.AioContext_t
	newOp = new(queueOp)

	context, err := getCtx(1) // Needs to be thread safe call.

	chk_err(err)
	ctx = context.ctx
	_ = ctx

	var iocbp = new(syscall.Iocb)

	newOp.iocbp = iocbp
	newOp.context = context
	newOp.op = &op

	switch {
	case op.Op == READ:
		op.Off = 0
		fallthrough
	case op.Op == READAT:
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

/*
* submitOps submits all give queueOps in the queue_list using the given
* context.
 */
func submitOps(context *Context, queue_list []*queueOp) {

	iocbp_list := make([]*syscall.Iocb, len(queue_list), len(queue_list))

	for i, qop_p := range queue_list {

		// Create the list of iocbp's
		iocbp_list[i] = qop_p.iocbp

		// save op & ctx as inflight
		qop_p_list, exists := scheduler.inflight_context[qop_p.context]
		if !exists {
			qop_p_list = make([]*queueOp, 0)
		}
		scheduler.inflight_context[qop_p.context] = append(qop_p_list, qop_p)
		scheduler.inflight_ops[qop_p.iocbp] = *qop_p.op
	}

	submitted, err := ioSubmitWrapper(context.ctx, len(queue_list), &iocbp_list[0])

	chk_err(err)
	if submitted != len(queue_list) {
		chk_err(FAIL)
	}
}
