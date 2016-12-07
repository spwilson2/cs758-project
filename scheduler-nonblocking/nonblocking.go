package nonblocking

// Non-blocking IO Scheduler

import (
	"runtime"
	"syscall"
)

/* Create file */
var Creat = syscall.Creat
var Open = syscall.Open

/* opens file w/ read params */
func OpenFile(path string) (fd int, err error) {
	mode := syscall.O_RDWR
	perm := uint32(0)
	return syscall.Open(path, mode, perm)
}

/* Writes p to file opened w/  fd */
func Write(fd int, p []byte) (n int, err error) {
	// make sure sched is running
	if scheduler.initialized != true {
		panic("Scheduler not initialized...")
	}

	// create op and send it to scheduler
	var ret_valid = new(bool)
	var ret_n = new(int)
	var ret_err = new(error)

	op := operation{WRITE, fd, "", p, 0, ret_valid, ret_n, ret_err}
	scheduler.channel <- op

	for *ret_valid != true {
		// spin until op complete
		runtime.Gosched() // Allow preemption
	}

	// results now valid, we can return them.
	return *ret_n, *ret_err
}

func WriteAt(fd int, off int, p []byte) (n int, err error) {
	// make sure sched is running
	if scheduler.initialized != true {
		panic("Scheduler not initialized...")
	}

	// create op and send it to scheduler
	var ret_valid = new(bool)
	var ret_n = new(int)
	var ret_err = new(error)

	op := operation{WRITEAT, fd, "", p, int64(off), ret_valid, ret_n, ret_err}
	scheduler.channel <- op

	for *ret_valid != true {
		// spin until op complete
		runtime.Gosched() // Allow preemption
	}

	// results now valid, we can return them.
	return *ret_n, *ret_err
}

/* Reads from file */
func Read(fd int, p []byte) (n int, err error) {
	// make sure sched is running
	if scheduler.initialized != true {
		panic("Scheduler not initialized...")
	}

	// create op and send it to scheduler
	var ret_valid = new(bool)
	var ret_n = new(int)
	var ret_err = new(error)

	op := operation{READ, fd, "", p, 0, ret_valid, ret_n, ret_err}
	scheduler.channel <- op

	for *ret_valid != true {
		// spin until op complete
		runtime.Gosched() // Allow preemption
	}

	// results now valid, we can return them.
	return *ret_n, *ret_err
}

/* Read from file, at offset off */
func ReadAt(fd int, off int, p []byte) (n int, err error) {
	// make sure sched is running
	if scheduler.initialized != true {
		panic("Scheduler not initialized...")
	}

	// create op and send it to scheduler
	var ret_valid = new(bool)
	var ret_n = new(int)
	var ret_err = new(error)

	op := operation{READAT, fd, "", p, int64(off), ret_valid, ret_n, ret_err}
	scheduler.channel <- op

	for *ret_valid != true {
		// spin until op complete
		runtime.Gosched() // Allow preemption
	}

	// results now valid, we can return them.
	return *ret_n, *ret_err
}
