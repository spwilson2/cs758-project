package nonblocking

import (
	"time"
)

const (
	EXECUTE_TIMEOUT   = 100 * time.Nanosecond
	GETEVENTS_TIMEOUT = 1000 * time.Nanosecond
)

type timeoutObject struct {
	close  chan bool //TODO: Add closing of the timeout.
	signal chan bool
	timer  time.Duration
}

func newTimeoutObject(delay time.Duration) *timeoutObject {
	new_timeout := new(timeoutObject)
	new_timeout.timer = delay
	return new_timeout
}

func (t *timeoutObject) begin() {
	for {
		time.Sleep(t.timer)
		t.signal <- true
	}
}
