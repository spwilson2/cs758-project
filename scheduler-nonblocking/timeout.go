package nonblocking

import (
	"time"
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

func (t *timeoutObject) CountDown() {
	for {
		time.Sleep(t.timer)
		t.signal <- true
	}
}
