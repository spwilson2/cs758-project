package main

import (
	"log"
	"time"
)

func main() {
	performAllBenchmarks()
}

// Benchmark option 1
func performAllBenchmarks() {
	performBlockingBenchmarks()
	performAsyncBenchmarks()
}

// Benchmark option 2
func performBlockingBenchmarks() {
	performSequentialBlockingWriteBenchmarks()
	performSequentialBlockingReadBenchmarks()
	performRandomBlockingWriteBenchmarks()
	performRandomBlockingReadBenchmarks()
}

// Benchmark option 3
func performAsyncBenchmarks() {
	performSequentialAsyncWriteBenchmarks()
	performSequentialAsyncReadBenchmarks()
	performRandomAsyncWriteBenchmarks()
	performRandomAsyncReadBenchmarks()
}

func performSequentialBlockingWriteBenchmarks() {
	defer un(trace("SBW"))
}

func performSequentialBlockingReadBenchmarks() {
	defer un(trace("SBR"))
}

func performSequentialAsyncWriteBenchmarks() {
	defer un(trace("SAW"))
}

func performSequentialAsyncReadBenchmarks() {
	defer un(trace("SAR"))
}

func performRandomBlockingWriteBenchmarks() {
	defer un(trace("RBW"))
}

func performRandomBlockingReadBenchmarks() {
	defer un(trace("RBR"))
}

func performRandomAsyncWriteBenchmarks() {
	defer un(trace("RAW"))
}

func performRandomAsyncReadBenchmarks() {
	defer un(trace("RAR"))
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

// 	// remove these once we use the packages, get compiler errs otherwise.
// 	_ = blocking.TestExport
// 	_ = nonblocking.TestExport

// 	// initialize scheduler
// 	log.Printf("Testing nonblocking")
// 	name := "aheh"
// 	buf := make([]byte, 100)
// 	op := nonblocking.Operation{nonblocking.READ, name, buf, 0}
// 	c := make(chan nonblocking.Operation)
// 	nonblocking.InitScheduler(c)
// 	c <- op // send the read operation
// 	time.Sleep(2 * time.Second)
// 	c <- op // send another read
// 	time.Sleep(2 * time.Second)
// }
