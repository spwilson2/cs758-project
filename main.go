package main

import (
	"flag"
	"fmt"
	"os"
	"syscall"
	"time"

	blocking "github.com/spwilson2/cs758-project/scheduler-blocking"
	nonblocking "github.com/spwilson2/cs758-project/scheduler-nonblocking"
)

func main() {

	initSchedulers()

	//runtime.GOMAXPROCS(2)
	flags, opSize := parseArgs()

	// Run Sequential Blocking Write Benchmark
	if *flags[0] {
		fmt.Printf("Running Sequential Blocking Write Benchmark\n")
		performSequentialBlockingWriteBenchmarks(*opSize)
	}

	// Run Sequential Blocking Read Benchmark
	if *flags[1] {
		fmt.Printf("Running Sequential Blocking Read Benchmark\n")
		performSequentialBlockingReadBenchmarks(*opSize)
	}

	// Run Random Blocking Write Benchmark
	if *flags[2] {
		fmt.Printf("Running Random Blocking Write Benchmark\n")
		performRandomBlockingWriteBenchmarks(*opSize)
	}

	// Run Random Blocking Read Benchmark
	if *flags[3] {
		fmt.Printf("Running Random Blocking Read Benchmark\n")
		performRandomBlockingReadBenchmarks(*opSize)
	}

	// Run Sequential Nonblocking Write Benchmark
	if *flags[4] {
		fmt.Printf("Running Sequential Nonblocking Write Benchmark\n")
		performSequentialAsyncWriteBenchmarks(*opSize)
	}
	// Run Sequential Nonblocking Read Benchmark
	if *flags[5] {
		fmt.Printf("Running Sequential Nonblocking Read Benchmark\n")
		performSequentialAsyncReadBenchmarks(*opSize)
	}

	// Run Random Nonblocking Write Benchmark
	if *flags[6] {
		fmt.Printf("Running Random Nonblocking Write Benchmark\n")
		performRandomAsyncWriteBenchmarks(*opSize)
	}

	// Run Random Nonblocking Write Benchmark
	if *flags[7] {
		fmt.Printf("Running Random Nonblocking Read Benchmark\n")
		performRandomAsyncReadBenchmarks(*opSize)
	}
}

func parseArgs() ([]*bool, *int) {
	flags := make([]*bool, 8)

	flags[0] = flag.Bool("SBW", false, "Should run the sequential blocking write benchmark")
	flags[1] = flag.Bool("SBR", false, "Should run the sequential blocking read benchmark")
	flags[2] = flag.Bool("RBW", false, "Should run the random blocking write benchmark")
	flags[3] = flag.Bool("RBR", false, "Should run the random blocking read benchmark")

	flags[4] = flag.Bool("SAW", false, "Should run the sequential nonblocking write benchmark")
	flags[5] = flag.Bool("SAR", false, "Should run the sequential nonblocking read benchmark")
	flags[6] = flag.Bool("RAW", false, "Should run the random nonblocking write benchmark")
	flags[7] = flag.Bool("RAR", false, "Should run the random nonblocking read benchmark")

	var opSize = flag.Int("size", 1000, "The size of the reads and/or writes to be performed")

	flag.Parse()

	return flags, opSize
}

func initSchedulers() {
	fmt.Printf("Initializing nonblocking IO scheduler\n")
	nonblockingChan := make(chan nonblocking.Operation)
	nonblocking.InitScheduler(nonblockingChan)

	fmt.Printf("Initializing blocking IO scheduler\n")
	blockingChan := make(chan blocking.Operation)
	blocking.InitScheduler(blockingChan)
}

func performSequentialBlockingWriteBenchmarks(opSize int) {
	name := "SBW.txt"
	buf := make([]byte, opSize)
	for i := 0; i < opSize; i++ {
		buf[i] = byte(i)
	}

	blocking.Creat(name, syscall.S_IRUSR|syscall.S_IWUSR)
	fd, err := blocking.Open(name, syscall.O_RDWR, 0)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error opening file\n")
	}

	for off := 0; off < opSize; off += 100 {
		scheduleBlockingWriteAt(fmt.Sprint("SBW @ ", off), fd, off, buf)
	}
}

func performSequentialBlockingReadBenchmarks(opSize int) {
	defer un(trace("SBR"))
}

func performSequentialAsyncWriteBenchmarks(opSize int) {
	name := "SAW.txt"
	buf := make([]byte, opSize)
	for i := 0; i < opSize; i++ {
		buf[i] = byte(i)
	}

	nonblocking.Creat(name, syscall.S_IRUSR|syscall.S_IWUSR)
	fd, err := nonblocking.Open(name, syscall.O_RDWR, 0)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error opening file\n")
	}

	for off := 0; off < opSize; off += 100 {
		scheduleNonblockingWriteAt(fmt.Sprint("SAW @ ", off), fd, off, buf)
	}
}

func performSequentialAsyncReadBenchmarks(opSize int) {
	defer un(trace("SAR"))
}

func performRandomBlockingWriteBenchmarks(opSize int) {
	defer un(trace("RBW"))
}

func performRandomBlockingReadBenchmarks(opSize int) {
	defer un(trace("RBR"))
}

func performRandomAsyncWriteBenchmarks(opSize int) {
	defer un(trace("RAW"))
}

func performRandomAsyncReadBenchmarks(opSize int) {
	defer un(trace("RAR"))
}

/*
	Helper functions for scheduling blocking and non blocking operations
*/

func scheduleBlockingWrite(id string, fd int, buf []byte) {
	defer un(trace(id))
	blocking.Write(fd, buf)
}

func scheduleBlockingWriteAt(id string, fd int, off int, buf []byte) {
	defer un(trace(id))
	blocking.WriteAt(fd, off, buf)
}

func scheduleBlockingRead(id string, fd int, buf []byte) {
	defer un(trace(id))
	blocking.Read(fd, buf)
}

func scheduleBlockingReadAt(id string, fd int, off int, buf []byte) {
	defer un(trace(id))
	blocking.ReadAt(fd, off, buf)
}

func scheduleNonblockingWrite(id string, fd int, buf []byte) {
	defer un(trace(id))
	nonblocking.Write(fd, buf)
}

func scheduleNonblockingWriteAt(id string, fd int, off int, buf []byte) {
	defer un(trace(id))
	nonblocking.WriteAt(fd, off, buf)
}

func scheduleNonblockingRead(id string, fd int, buf []byte) {
	defer un(trace(id))
	nonblocking.Read(fd, buf)
}

func scheduleNonblockingReadAt(id string, fd int, off int, buf []byte) {
	defer un(trace(id))
	nonblocking.ReadAt(fd, off, buf)
}

/*
	These functions can be used together to benchmark a functions. Simply provide
	a unique identifier that is used for fmtging when calling the functions using
	the defer keyword:
	e.g. defer un(trace("testId"))
	Basically, the nested function trace() is called immediately and that grabs
	the start time which is passed to un(); however, un() is not executed until
	the end of the function. This Is pretty clean, but you'll need a seperate
	function that must start with this and continue no code beyond what you
	want to benchmark.
*/

func trace(id string) (string, time.Time) {
	//fmt.Printf("%s Benchmark running\n", id)
	start := time.Now()
	return id, start
}

func un(id string, start time.Time) time.Duration {
	elapsed := time.Since(start)
	//fmt.Printf("%s Benchmark complete\n", id)
	fmt.Printf("%d nanoseconds\n\n", elapsed.Nanoseconds())
	return elapsed
}
