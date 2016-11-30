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
		performSequentialBlockingWriteBenchmarks(*opSize)
	}

	// Run Sequential Blocking Read Benchmark
	if *flags[1] {
		performSequentialBlockingReadBenchmarks(*opSize)
	}

	// Run Random Blocking Write Benchmark
	if *flags[2] {
		performRandomBlockingWriteBenchmarks(*opSize)
	}

	// Run Random Blocking Read Benchmark
	if *flags[3] {
		performRandomBlockingReadBenchmarks(*opSize)
	}

	// Run Sequential Nonblocking Write Benchmark
	if *flags[4] {
		performSequentialAsyncWriteBenchmarks(*opSize)
	}
	// Run Sequential Nonblocking Read Benchmark
	if *flags[5] {
		performSequentialAsyncReadBenchmarks(*opSize)
	}

	// Run Random Nonblocking Write Benchmark
	if *flags[6] {
		performRandomAsyncWriteBenchmarks(*opSize)
	}

	// Run Random Nonblocking Write Benchmark
	if *flags[7] {
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
	elapsed := new(int64)
	for off := 0; off < opSize; off += 100 {
		scheduleBlockingWriteAt("SBW", fd, off, buf, elapsed)
		executionTime += *elapsed
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
	elapsed := new(int64)
	for off := 0; off < opSize; off += 100 {
		scheduleBlockingReadAt("SBW", fd, off, buf, elapsed)
		executionTime += *elapsed
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

	executionTime := int64(0)
	elapsed := new(int64)
	for off := 0; off < opSize; off += 100 {
		scheduleNonblockingWriteAt("SAW", fd, off, buf, elapsed)
		executionTime += *elapsed
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
	elapsed := new(int64)
	for off := 0; off < opSize; off += 100 {
		scheduleNonblockingReadAt("SAW", fd, off, buf, elapsed)
		executionTime += *elapsed
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

func scheduleBlockingWrite(id string, fd int, buf []byte, elapsed *int64) {
	defer un(trace(id, elapsed))
	blocking.Write(fd, buf)
}

func scheduleBlockingWriteAt(id string, fd int, off int, buf []byte, elapsed *int64) {
	defer un(trace(id, elapsed))
	blocking.WriteAt(fd, off, buf)
}

func scheduleBlockingRead(id string, fd int, buf []byte, elapsed *int64) {
	defer un(trace(id, elapsed))
	blocking.Read(fd, buf)
}

func scheduleBlockingReadAt(id string, fd int, off int, buf []byte, elapsed *int64) {
	defer un(trace(id, elapsed))
	blocking.ReadAt(fd, off, buf)
}

func scheduleNonblockingWrite(id string, fd int, buf []byte, elapsed *int64) {
	defer un(trace(id, elapsed))
	nonblocking.Write(fd, buf)
}

func scheduleNonblockingWriteAt(id string, fd int, off int, buf []byte, elapsed *int64) {
	defer un(trace(id, elapsed))
	nonblocking.WriteAt(fd, off, buf)
}

func scheduleNonblockingRead(id string, fd int, buf []byte, elapsed *int64) {
	defer un(trace(id, elapsed))
	nonblocking.Read(fd, buf)
}

func scheduleNonblockingReadAt(id string, fd int, off int, buf []byte, elapsed *int64) {
	defer un(trace(id, elapsed))
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

func trace(id string, elapsed *int64) (string, time.Time, *int64) {
	return id, time.Now(), elapsed
}

func un(id string, start time.Time, elapsed *int64) {
	*elapsed = time.Since(start).Nanoseconds()
}
