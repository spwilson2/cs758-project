package main

import (
	"log"
	"syscall"
	"time"

	blocking "github.com/spwilson2/cs758-project/scheduler-blocking"
	nonblocking "github.com/spwilson2/cs758-project/scheduler-nonblocking"
)

func main() {

	//runtime.GOMAXPROCS(2)

	log.Printf("Performing nonblocking IO benchmark tests\n")
	nonblockingChan := make(chan nonblocking.Operation)
	nonblocking.InitScheduler(nonblockingChan)
	performAsyncBenchmarks()

	log.Printf("Performing blocking IO benchmark tests\n")
	performBlockingBenchmarks()
}

func performBlockingBenchmarks() {
	performSequentialBlockingWriteBenchmarks()
	performSequentialBlockingReadBenchmarks()
	performRandomBlockingWriteBenchmarks()
	performRandomBlockingReadBenchmarks()
}

func performAsyncBenchmarks() {
	performSequentialAsyncWriteBenchmarks()
	performSequentialAsyncReadBenchmarks()
	performRandomAsyncWriteBenchmarks()
	performRandomAsyncReadBenchmarks()
}

func performSequentialBlockingWriteBenchmarks() {
	name := "SBW.txt"
	buf := make([]byte, 1000)
	for i := 0; i < 1000; i++ {
		buf[i] = byte(i)
	}

	blocking.Create(name)
	file, err := blocking.Open(name)
	if err != nil {
		log.Printf("error opening file\n")
	}

	defer un(trace("SBW"))
	off := int64(0)
	for i := 0; i < 10; i++ {
		file.WriteAt(buf, off)
		off += 1000
	}
}

func performSequentialBlockingReadBenchmarks() {
	//defer un(trace("SBR"))
}

func performSequentialAsyncWriteBenchmarks() {
	name := "SAW.txt"
	buf := make([]byte, 1000)
	for i := 0; i < 1000; i++ {
		buf[i] = byte(i)
	}

	nonblocking.Create(name)
	fd, err := nonblocking.Open(name, syscall.O_RDWR, 0)
	if err != nil {
		log.Printf("error opening file\n")
	}

	defer un(trace("SAW"))

	off := 0
	for i := 0; i < 10; i++ {
		nonblocking.WriteAt(fd, off, buf)
		off += 100
	}
}

func performSequentialAsyncReadBenchmarks() {
	//defer un(trace("SAR"))
}

func performRandomBlockingWriteBenchmarks() {
	//defer un(trace("RBW"))
}

func performRandomBlockingReadBenchmarks() {
	//defer un(trace("RBR"))
}

func performRandomAsyncWriteBenchmarks() {
	//defer un(trace("RAW"))
}

func performRandomAsyncReadBenchmarks() {
	//defer un(trace("RAR"))
}

/*
	These functions can be used together to benchmark a functions. Simply provide
	a unique identifier that is used for logging when calling the functions using
	the defer keyword:
	e.g. defer un(trace("testId"))
	Basically, the nested function trace() is called immediately and that grabs
	the start time which is passed to un(); however, un() is not executed until
	the end of the function. This Is pretty clean, but you'll need a seperate
	function that must start with this and continue no code beyond what you
	want to benchmark.
*/

func trace(id string) (string, time.Time) {
	log.Printf("Benchmarking: %s", id)
	start := time.Now()
	return id, start
}

func un(id string, start time.Time) time.Duration {
	elapsed := time.Since(start)
	log.Printf("%s completed in %d nanoseconds\n\n", id, elapsed.Nanoseconds())
	return elapsed
}

// func main() {

// send various operations to scheduler
// name := "hello.txt"
// buf := make([]byte, 100)
// off := 0
// fd, _ := nonblocking.Open(name, syscall.O_RDWR, 0)

// // READ
// nonblocking.Read(fd, buf)
// time.Sleep(2 * time.Second)

// // READAT
// nonblocking.ReadAt(fd, off, buf)
// time.Sleep(2 * time.Second)

// // WRITE
// buf = []byte("World Hello")
// nonblocking.Write(fd, buf)
// time.Sleep(2 * time.Second)

// // WRITEAT
// buf = []byte("World")
// off = 5
// nonblocking.WriteAt(fd, off, buf)
// time.Sleep(2 * time.Second)

//}
