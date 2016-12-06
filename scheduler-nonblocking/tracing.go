package nonblocking

import (
	"syscall"

	tracer "github.com/spwilson2/cs758-project/tracer"
)

/* Wrapper function to enable dynamic setting of tracing. */
var ioSetupWrapper = syscall.IoSetup
var ioGeteventsWrapper = syscall.IoGetevents
var ioSubmitWrapper = syscall.IoSubmit

var (
	T_SCHEDULER_INIT,
	T_GET_CONTEXT,
	T_GET_EVENTS_SYSCALL,
	T_SUBMIT_SYSCALL,
	T_SETUP_SYSCALL tracer.Event_t
)

var schedulerTraceList tracer.TraceList

/* Print out a list of traces for the trace list local to the scheduler. */
func PrintTrace() {
	schedulerTraceList.PrintLog()
}

func initTracer(enable bool) {
	var err error
	if enable {
		T_SCHEDULER_INIT, err = tracer.NewEventType("InitScheduler")
		chk_err(err)
		T_GET_CONTEXT, err = tracer.NewEventType("GetContext")
		chk_err(err)
		T_GET_EVENTS_SYSCALL, err = tracer.NewEventType("IoGetevents")
		chk_err(err)
		T_SUBMIT_SYSCALL, err = tracer.NewEventType("IoSubmit")
		chk_err(err)
		T_SETUP_SYSCALL, err = tracer.NewEventType("IoSetup")
		chk_err(err)
		ioSetupWrapper = ioSetupTracer
		ioSubmitWrapper = ioSubmitTracer
		ioGeteventsWrapper = ioGeteventsTracer
	}
}

func ioSubmitTracer(ctx_id syscall.AioContext_t, nr int, iocbpp **syscall.Iocb) (err error) {
	trace := tracer.NewTraceEvent(T_SUBMIT_SYSCALL, &schedulerTraceList)
	trace.Start()
	err = syscall.IoSubmit(ctx_id, nr, iocbpp)
	trace.Stop()
	return
}

func ioGeteventsTracer(ctx_id syscall.AioContext_t, nr_min int, nr int, events *syscall.IoEvent, timeout *syscall.Timespec) (n int) {
	trace := tracer.NewTraceEvent(T_GET_EVENTS_SYSCALL, &schedulerTraceList)
	trace.Start()
	n = syscall.IoGetevents(ctx_id, nr_min, nr, events, timeout)
	trace.Stop()
	return
}

func ioSetupTracer(nr_events uint, ctx_idp *syscall.AioContext_t) (err error) {
	trace := tracer.NewTraceEvent(T_SETUP_SYSCALL, &schedulerTraceList)
	trace.Start()
	err = syscall.IoSetup(nr_events, ctx_idp)
	trace.Stop()
	return
}
