package tracer

import (
	"fmt"
	"sync"
	atomic "sync/atomic"
	"time"
)

var GlobalTraceList TraceList

type Event_t int

const (
	T_READ         Event_t = 0
	T_WRITE        Event_t = 1
	T_READ_STRING          = "Read"
	T_WRITE_STRING         = "Write"
)

type TraceEvent struct {
	startTime time.Time
	stopTime  time.Time
	id        int32
	eventType Event_t
}

type TraceList struct {
	list      []*TraceEvent
	lock      sync.Mutex
	idCounter int32
}

func NewTraceEvent(traceType Event_t, list *TraceList) *TraceEvent {
	var newTrace TraceEvent
	newTrace.eventType = traceType
	list.addTrace(&newTrace)
	newTrace.id = atomic.AddInt32(&list.idCounter, 1)
	return &newTrace
}

func (event *TraceEvent) Start() {
	event.startTime = time.Now()
}

func (event *TraceEvent) Stop() {
	event.stopTime = time.Now()
}

func (list *TraceList) addTrace(event *TraceEvent) {
	list.lock.Lock()
	list.list = append(list.list, event)
	list.lock.Unlock()
}

func (log *TraceList) PrintLog() {
	for _, entry := range log.list {

		var opString string

		switch entry.eventType {
		case T_READ:
			opString = T_READ_STRING
		case T_WRITE:
			opString = T_WRITE_STRING
		}

		fmt.Printf("Operation: %-8s ", opString)
		fmt.Printf("Length: %-10d", entry.stopTime.Sub(entry.startTime).Nanoseconds())
		fmt.Printf("ID: %-5d\n", entry.id)
		//fmt.Printf("%v\n", *entry)
	}
}
