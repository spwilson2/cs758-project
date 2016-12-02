package main

import (
	"flag"
	"fmt"
	"math"
	rand "math/rand"
	"os"
	"strconv"
	"strings"
	"sync"
	atomic "sync/atomic"
	"time"

	sut "github.com/spwilson2/cs758-project/scheduler-blocking"
	//sut "github.com/spwilson2/cs758-project/scheduler-nonblocking"
)

const GEN_FILE_BASENAME = "generated/testfile-"
const GEN_FILE_SUFFIX = ".gen"

type Event_t int

const (
	T_READ  Event_t = 0
	T_WRITE Event_t = 1
)

type TraceEvent struct {
	startTime time.Time
	stopTime  time.Time
	id        int32
	traceType int32
}

type TraceList struct {
	list      []*TraceEvent
	lock      sync.Mutex
	idCounter int32
}

var traces TraceList

func NewTraceEvent(traceType Event_t) *TraceEvent {
	var newTrace TraceEvent
	traces.addTrace(&newTrace)
	newTrace.id = atomic.AddInt32(&traces.idCounter, 1)
	return &newTrace
}

func (event *TraceEvent) start() {
	event.startTime = time.Now()
}

func (event *TraceEvent) stop() {
	event.stopTime = time.Now()
}

func (list *TraceList) addTrace(event *TraceEvent) {
	list.lock.Lock()
	list.list = append(list.list, event)
	list.lock.Unlock()
}

func (log *TraceList) printLog() {
	for _, entry := range log.list {
		fmt.Printf("%v\n", *entry)
	}
}

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

	getArgs()
	initScheduler()
	runTest()
	traces.printLog()

}

func getArgs() {

	flag.BoolVar(&f_blocking, "blocking", true, "Use blocking interface false uses async")
	flag.IntVar(&f_threads, "t", 0, "Number of threads to test with")
	flag.IntVar(&f_readSize, "rsize", 0, "Size of reads to execute")
	flag.IntVar(&f_writeSize, "wsize", 0, "Size of writes to execute")
	flag.IntVar(&f_numWrites, "nwrites", 0, "Number of writes to execute")
	flag.IntVar(&f_numReads, "nreads", 0, "Number of reads to execute")
	flag.IntVar(&f_readOffset, "roff", 0, "Offset for each additional read")
	flag.IntVar(&f_writeOffset, "woff", 0, "Offset for each additional write")
	flag.IntVar(&f_numFiles, "nfiles", 0, "Number of different files to dispatch r/w's to")

	var file_list string
	flag.StringVar(&file_list, "files", "", "Comma separated list of files to dispatch r/w's, overrids nfiles")

	flag.Parse()

	/* Check for invalid arguments. */
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

	// TODO
	/* Assert that we have set at least one testable configuration. */

	if file_list != "" {
		f_files = strings.Split(file_list, ",")
	}
}

func initScheduler() {
	c := make(chan sut.Operation)
	sut.InitScheduler(c)
}

func runTest() {

	var use_threads bool
	var do_read bool
	var do_write bool
	var do_mixed bool
	var file_list []string
	var file_list_handles []int

	var writeOrder []bool

	var reads_per_file int
	var writes_per_file int

	// Only used in mixed case
	var averageOffset int

	/* Setup test variables */
	rand.Seed(time.Now().UnixNano())

	use_threads = (f_threads >= 1)

	if len(f_files) >= 1 {
		file_list = f_files
	} else {
		//Approximate the max file size we will need to read.
		maxSize := int(((float64(f_numReads) / float64(f_numFiles)) * float64(f_readOffset))) + f_readSize
		file_list = genFiles(f_numFiles, int64(maxSize))
	}

	reads_per_file = f_numReads / len(file_list)
	writes_per_file = f_numWrites / len(file_list)
	do_read = (reads_per_file >= 1)
	do_write = (writes_per_file >= 1)

	// Open the files to be tested.
	for _, file_name := range file_list {
		handle, err := sut.OpenFile(file_name)
		panic_chk(err)
		file_list_handles = append(file_list_handles, handle)
	}

	// Generate an ordering of the operations to use - not going to
	// truely limit reads and writes to their count, but will help.
	if do_mixed {
		split := int((float64(f_numReads-f_numWrites) / float64(math.MaxInt64)) * math.MaxInt64)
		for i := 0; i < f_numReads+f_numWrites; i++ {
			writeOrder = append(writeOrder, rand.Int() >= split)
		}

		// Approximate the offset we will use in the file.
		// Divide first to avoid overflow.
		combinedOps := f_numReads + f_numWrites

		averageOffset = int(float64(f_readOffset)/float64(combinedOps)*float64(f_numReads) + (float64(f_writeOffset)/float64(combinedOps))*float64(f_numWrites))
	}

	//TODO: For now we use the same read/write buffer. Will need to change
	// for tests.
	buffer := make([]byte, int(math.Max(float64(f_writeSize), float64(f_readSize))))

	thread_collector := make(chan bool)
	thread_count := 0

	/* Execute the tests. */
	if do_mixed {

		ops_per_file := reads_per_file + writes_per_file

		for iter, handle := range file_list_handles {
			for op := 0; op < ops_per_file; op++ {
				op_index := (iter * ops_per_file) + op
				offset := (op) * averageOffset

				thread_count++
				scheduleOp(handle, offset, buffer, writeOrder[op_index], use_threads, thread_collector)
			}
		}
	} else {
		if do_read {
			for _, handle := range file_list_handles {
				for op := 0; op < reads_per_file; op++ {
					offset := op * f_readOffset

					thread_count++
					scheduleOp(handle, offset, buffer, true, use_threads, thread_collector)
				}
			}
		}
		if do_write {
			for _, handle := range file_list_handles {
				for op := 0; op < writes_per_file; op++ {
					offset := op * f_writeOffset

					thread_count++
					scheduleOp(handle, offset, buffer, false, use_threads, thread_collector)
				}
			}
		}
	}

	if use_threads {
		/* Finish the multithreaded test. */
		for thread := 0; thread < thread_count; thread++ {
			_ = <-thread_collector
		}
	}

	/* Cleanup from the tests. */

}

func scheduleOp(file, offset int, buffer []byte, readNotWrite, threaded bool, collector chan bool) {

	var eventType Event_t
	var op func(int, int, []byte) (int, error)

	if readNotWrite {
		op = sut.ReadAt
		eventType = T_READ
	} else {
		op = sut.WriteAt
		eventType = T_WRITE
	}

	execOp := func() {
		trace := NewTraceEvent(eventType)
		trace.start()
		ret, err := op(file, offset, buffer)
		trace.stop()

		panic_chk(err)
		assert(ret == len(buffer))
		if threaded {
			collector <- true
		}
	}

	if threaded {
		go execOp()
	} else {
		execOp()
	}
}

// Generate a list of files (making sure they exist)
func genFiles(num int, size int64) []string {
	var file_list []string

	var buf = []byte{0}

	for val := 0; val < num; val++ {
		file_name := GEN_FILE_BASENAME + strconv.Itoa(val) + GEN_FILE_SUFFIX
		file_p, err := os.OpenFile(file_name, os.O_CREATE|os.O_WRONLY, 0777)
		panic_chk(err)
		file_list = append(file_list, file_name)
		_, err = file_p.WriteAt(buf, size)
		panic_chk(err)
		err = file_p.Close()
		panic_chk(err)
	}
	return file_list
}

//func performSequentialAsyncReadBenchmarks(opSize int) {
//
//	name := "SA.txt"
//	buf := make([]byte, opSize)
//
//	//nonblocking.Creat(name, syscall.S_IRUSR|syscall.S_IWUSR)
//	fd, err := nonblocking.Open(name, syscall.O_RDWR, 0)
//	if err != nil {
//		fmt.Fprintf(os.Stderr, "error opening file\n")
//	}
//
//	executionTime := int64(0)
//	elapsed := new(int64)
//	for off := 0; off < opSize; off += 100 {
//		scheduleNonblockingReadAt("SAW", fd, off, buf, elapsed)
//		executionTime += *elapsed
//	}
//	fmt.Println(executionTime)
//}
//
//func scheduleNonblockingReadAt(id string, fd int, off int, buf []byte, elapsed *int64) {
//	defer un(trace(id, elapsed))
//	nonblocking.ReadAt(fd, off, buf)
//}

func panic_chk(err error) {
	if err != nil {
		fmt.Println(err)
		panic("Panic!")
	}
}

func assert(val bool, message ...interface{}) {
	if !val {
		fmt.Printf("%v\n", message)
		panic(nil)
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
