package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"syscall"
	"time"

	//blocking "github.com/spwilson2/cs758-project/scheduler-blocking"
	//nonblocking "github.com/spwilson2/cs758-project/scheduler-nonblocking"
	sut "github.com/spwilson2/cs758-project/scheduler-nonblocking"
)

const GEN_FILE_BASENAME = "testfile-"
const GEN_FILE_SUFFIX = ".gen"

/*Vars set by flags.*/
var f_threads int
var f_readSize int
var f_writeSize int
var f_numWrites int
var f_numReads int
var f_readOffset int
var f_writeOffset int
var f_numFiles int
var f_files []string
var f_blocking bool

func main() {

	parseArgs()
	initSchedulers()

}

func parseArgs() {
	//var Usage = func() {
	//	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	//	PrintDefaults()
	//}

	f_blocking = *flag.Bool("blocking", true, "Use blocking interface false uses async")
	f_threads = *flag.Int("t", 0, "Number of threads to test with")
	f_readSize = *flag.Int("rsize", 0, "Size of reads to execute")
	f_writeSize = *flag.Int("wsize", 0, "Size of writes to execute")
	f_numWrites = *flag.Int("nwrites", 0, "Number of writes to execute")
	f_numReads = *flag.Int("nreads", 0, "Number of reads to execute")
	f_readOffset = *flag.Int("roff", 0, "Offset for each additional read")
	f_writeOffset = *flag.Int("woff", 0, "Offset for each additional write")
	f_numFiles = *flag.Int("nfiles", 0, "Number of different files to dispatch r/w's to")

	switch {
	case f_threads < 0:
		fallthrough
	case f_readSize < 0:
		fallthrough
	case f_writeSize < 0:
		fallthrough
	case f_numWrites < 0:
		fallthrough
	case f_numReads < 0:
		fallthrough
	case f_readOffset < 0:
		fallthrough
	case f_writeOffset < 0:
		fallthrough
	case f_numFiles < 0:
		flag.PrintDefaults()
	default:
	}

	file_list := *flag.String("files", "", "Comma separated list of files to dispatch r/w's, overrids nfiles")

	if file_list != "" {
		f_files = strings.Split(file_list, ",")
	}
}

func initSchedulers() {
	nonblockingChan := make(chan nonblocking.Operation)
	nonblocking.InitScheduler(nonblockingChan)

	blockingChan := make(chan blocking.Operation)
	blocking.InitScheduler(blockingChan)
}

func runTest() {
	var use_threads bool
	var do_read bool
	var do_write bool
	var do_mixed bool
	var file_list []string
	var file_list_handles []int

	var reads_per_file int
	var writes_per_file int

	setupTests := func() {
		use_threads = (f_threads >= 1)
		if len(f_files) >= 1 {
			file_list = f_files
		} else {
			file_list = genFiles(f_numFiles)
		}

		reads_per_file = f_numReads / len(file_list)
		writes_per_file = f_numWrites / len(file_list)
		do_read = (reads_per_file >= 1)
		do_write = (writes_per_file >= 1)

		// Open the files first.
		for _, file_name := range file_list {
			handle, err := sut.OpenFile(file_name)
			panic_chk(err)
			Append(file_list_handles, handle)
		}
	}

	if do_mixed {

	} else {
		for _, handle := range file_list_handles {
			for op_num := range reads_per_file {

			}
		}
		for _, handle := range file_list_handles {
			for op_num := range writes_per_file {

			}
		}
	}

}

func mixedTests() {
}

// Generate a list of files (making sure they exist)
func genFiles(num int) []string {
	var file_list []string
	for val := range num {
		file_name := GEN_FILE_BASENAME + string(val) + GEN_FILE_SUFFIX
		file_p, err := os.OpenFile(file_name, os.O_CREATE, 0)
		panic_chk(err)
		Append(file_list)
		err = file_p.Close()
		panic_chk(err)
	}
	return file_list
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

func panic_chk(err error) {
	if err != nil {
		fmt.Println(print(err))
		panic(err)
	}
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
