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
	nonblockingChan := make(chan nonblocking.Operation)
	nonblocking.InitScheduler(nonblockingChan)

	blockingChan := make(chan blocking.Operation)
	blocking.InitScheduler(blockingChan)
}

func performSequentialBlockingWriteBenchmarks(opSize int) {
	name := "SB.txt"
	buf := make([]byte, opSize)
	for i := 0; i < opSize; i++ {
		buf[i] = byte(i)
	}

	blocking.Creat(name, syscall.S_IRUSR|syscall.S_IWUSR)
	fd, err := blocking.Open(name, syscall.O_RDWR, 0)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error opening file\n")
	}

	executionTime := int64(0)
	for off := 0; off < opSize; off += 100 {
		executionTime += scheduleBlockingWriteAt(fmt.Sprint("SBW @ ", off), fd, off, buf)
	}
	fmt.Println(executionTime)
}

func performSequentialBlockingReadBenchmarks(opSize int) {
	name := "SB.txt"
	buf := make([]byte, opSize)

	//blocking.Creat(name, syscall.S_IRUSR|syscall.S_IWUSR)
	fd, err := blocking.Open(name, syscall.O_RDWR, 0)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error opening file\n")
	}

	executionTime := int64(0)
	for off := 0; off < opSize; off += 100 {
		executionTime += scheduleBlockingReadAt(fmt.Sprint("SBW @ ", off), fd, off, buf)
	}
	fmt.Println(executionTime)
}

func performSequentialAsyncWriteBenchmarks(opSize int) {
	name := "SA.txt"
	buf := make([]byte, opSize)
	for i := 0; i < opSize; i++ {
		buf[i] = byte(i)
	}

	nonblocking.Creat(name, syscall.S_IRUSR|syscall.S_IWUSR)
	fd, err := nonblocking.Open(name, syscall.O_RDWR, 0)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error opening file\n")
	}
	fmt.Println(executionTime)

	executionTime := int64(0)
	for off := 0; off < opSize; off += 100 {
		executionTime += scheduleNonblockingWriteAt(fmt.Sprint("SAW @ ", off), fd, off, buf)
	}
	fmt.Println(executionTime)
}

func performSequentialAsyncReadBenchmarks(opSize int) {
	name := "SA.txt"
	buf := make([]byte, opSize)

	//nonblocking.Creat(name, syscall.S_IRUSR|syscall.S_IWUSR)
	fd, err := nonblocking.Open(name, syscall.O_RDWR, 0)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error opening file\n")
	}

	executionTime := int64(0)
	for off := 0; off < opSize; off += 100 {
		executionTime += scheduleNonblockingReadAt(fmt.Sprint("SAW @ ", off), fd, off, buf)
	}
	fmt.Println(executionTime)
}

func performRandomBlockingWriteBenchmarks(opSize int) {
	elapsed := new(int64)
	defer un(trace("RBW", elapsed))
	fmt.Println(*elapsed)
}

func performRandomBlockingReadBenchmarks(opSize int) {
	elapsed := new(int64)
	defer un(trace("RBR", elapsed))
	fmt.Println(*elapsed)
}

func performRandomAsyncWriteBenchmarks(opSize int) {
	elapsed := new(int64)
	defer un(trace("RAW", elapsed))
	fmt.Println(*elapsed)
}

func performRandomAsyncReadBenchmarks(opSize int) {
	elapsed := new(int64)
	defer un(trace("RAR", elapsed))
	fmt.Println(*elapsed)
}

/*
	Helper functions for scheduling blocking and non blocking operations
*/

func scheduleBlockingWrite(id string, fd int, buf []byte) int64 {
	elapsed := new(int64)
	defer un(trace(id, elapsed))
	blocking.Write(fd, buf)
	return *elapsed
}

func scheduleBlockingWriteAt(id string, fd int, off int, buf []byte) int64 {
	elapsed := new(int64)
	defer un(trace(id, elapsed))
	blocking.WriteAt(fd, off, buf)
	return *elapsed
}

func scheduleBlockingRead(id string, fd int, buf []byte) int64 {
	elapsed := new(int64)
	defer un(trace(id, elapsed))
	blocking.Read(fd, buf)
	return *elapsed
}

func scheduleBlockingReadAt(id string, fd int, off int, buf []byte) int64 {
	elapsed := new(int64)
	defer un(trace(id, elapsed))
	blocking.ReadAt(fd, off, buf)
	return *elapsed
}

func scheduleNonblockingWrite(id string, fd int, buf []byte) int64 {
	elapsed := new(int64)
	defer un(trace(id, elapsed))
	nonblocking.Write(fd, buf)
	return *elapsed
}

func scheduleNonblockingWriteAt(id string, fd int, off int, buf []byte) int64 {
	elapsed := new(int64)
	defer un(trace(id, elapsed))
	nonblocking.WriteAt(fd, off, buf)
	return *elapsed
}

func scheduleNonblockingRead(id string, fd int, buf []byte) int64 {
	elapsed := new(int64)
	defer un(trace(id, elapsed))
	nonblocking.Read(fd, buf)
	return *elapsed
}

func scheduleNonblockingReadAt(id string, fd int, off int, buf []byte) int64 {
	elapsed := new(int64)
	defer un(trace(id, elapsed))
	nonblocking.ReadAt(fd, off, buf)
	return *elapsed
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

func trace(id string, elapsed *int64) (string, time.Time, *int64) {
	start := time.Now()
	return id, start, elapsed
}

func un(id string, start time.Time, elapsed *int64) {
	*elapsed = time.Since(start).Nanoseconds()
}
