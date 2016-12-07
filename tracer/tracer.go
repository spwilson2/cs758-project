package tracer

import (
	"errors"
	"fmt"
	"sync"
	atomic "sync/atomic"
	"time"
)

var GlobalTraceList *TraceList

type Event_t int

const (
	T_READ  Event_t = 0
	T_WRITE Event_t = 1
)

var tracerIdCounter int32 = 2

var EventMap map[Event_t]string = map[Event_t]string{
	T_READ:  "Read",
	T_WRITE: "Write",
}

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

func init() {
	GlobalTraceList = new(TraceList)
}

func NewEventType(event string) (Event_t, error) {
	id := atomic.AddInt32(&tracerIdCounter, 1)
	new_event := Event_t(id)
	if val, exists := EventMap[new_event]; exists {
		return 0, errors.New(fmt.Sprintf("Unable to create event %d:%s, already %s exists.", id, event, val))
	}
	EventMap[new_event] = event
	return new_event, nil
}

func NewTraceEvent(traceType Event_t, list *TraceList) *TraceEvent {
	newTrace := new(TraceEvent)
	newTrace.eventType = traceType
	list.addTrace(newTrace)
	newTrace.id = atomic.AddInt32(&list.idCounter, 1)
	return newTrace
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

		var opString string = EventMap[entry.eventType]

		fmt.Printf("Operation: %-15s ", opString)
		fmt.Printf("Length: %-10d ", entry.stopTime.Sub(entry.startTime).Nanoseconds())
		fmt.Printf("ID: %-5d \n", entry.id)
		//fmt.Printf("Start: %-10v ", entry.startTime)
		//fmt.Printf("Stop: %-10v \n", entry.stopTime)
		//fmt.Printf("%v\n", *entry)
	}
}
