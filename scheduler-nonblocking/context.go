package nonblocking

import (
	"syscall"

	tracer "github.com/spwilson2/cs758-project/tracer"
)

const BAD_CONTEXT_REQUEST = Error("Too large of a context was requested.")

const AIO_CONTEXT_MAX = 500000
const AIO_CONTEXT_MIN = 1000
const MIN_LOWER_RATE = 10
const MIN_LOWER_LIMIT = 10
const CONTEXT_REQUEST_MULTIPLIER = 2

var aio_contexts []*Context
var aio_context_max = AIO_CONTEXT_MAX
var aio_context_min = AIO_CONTEXT_MIN

/* reference counting for context */
type Context struct {
	ctx        syscall.AioContext_t
	references int
	maxsize    int
}

/* setup context for aio, manages as reference counter */
func getCtx(num int) (*Context, error) {

	trace := tracer.NewTraceEvent(T_GET_CONTEXT, &schedulerTraceList)
	trace.Start()
	defer trace.Stop()

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
	//err := syscall.IoSetup(uint(num_context_request), &ctx)
	err := ioSetupWrapper(uint(num_context_request), &ctx)

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
