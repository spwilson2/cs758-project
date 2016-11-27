package main

import (
	"log"
	"time"
	"os"
	"bufio"
	"strconv"

	blocking "github.com/spwilson2/cs758-project/scheduler-blocking"
	nonblocking "github.com/spwilson2/cs758-project/scheduler-nonblocking"
)

// We test 1 KB, 10 KB, 100 KB, 1MB, and 10 MB (5 steps of a log scale) 
sizesPerTest := [5]{1000, 10000, 100000, 1000000, 10000000}

func main() {
	if len(os.Args) == 2 && os.Args[1] == "--default" {
		performAllBenchmarks()
		return
	}

	usagePrompt()
}

func usagePrompt() (int) {
	fmt.Println("***************************************************************************")
	fmt.Println("** Golang Blocking and Asynchonous Scheduling Simulation Benmarking tool **")
	fmt.Println("***************************************************************************\n")

	fmt.Println("Testing Options")
	fmt.Println("1) All Benchmarks")
	fmt.Println("2) BIO Benchmarks")
	fmt.Println("3) AIO Benchmarks")
	
	fmt.Print("Benchmark option: ")

	return getTestOption()
}

func getTestOption() (int) {
	reader := bufio.NewReader(os.Stdin)
	option, _ := reader.ReadString('\n')
	
	optionInt, _ := strconv.Atoi(option)
	return optionInt
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

}

func performSequentialBlockingReadBenchmarks() {

}

func performSequentialAsyncWriteBenchmarks() {

}

func performSequentialAsyncReadBenchmarks() {

}

func performRandomBlockingWriteBenchmarks() {

}

func performRandomBlockingReadBenchmarks() {

}

func performRandomAsyncWriteBenchmarks() {

}

func performRandomAsyncReadBenchmarks() {

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

func trace(id string) (string, Time) {
	log.Printf("Benchmarking: %s", id)
	start := time.Now()
	return (id, start)
}

func un(id string, Time start) (Duration) {
	elapsed := time.Since(time)
	log.Printf("%s completed in %d", id, elapsed.Nanoseconds())
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
