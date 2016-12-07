package main

/* Note: This file WILL NOT compile. You must use the run-benchmarks script to
* generate files with the approprate imports! */

import (
	"flag"
	"fmt"
	tracer "github.com/spwilson2/cs758-project/tracer"
	"math"
	rand "math/rand"
	"os"
	"strconv"
	"strings"
	"syscall"
	"time"
	//SCHEDULER_UNDER_TEST
)

const GEN_FILE_BASENAME = "testfile-"
const GEN_FILE_SUFFIX = ".gen"
const BLOCKSIZE int = 512

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
var f_mixOps bool
var f_path string

func main() {

	getArgs()
	sut.InitScheduler(true)
	runTest()
	sut.EndScheduler()
	tracer.GlobalTraceList.PrintLog()
	sut.PrintTrace()

}

func getArgs() {

	flag.IntVar(&f_threads, "threads", 0, "Number of threads to test with")
	flag.IntVar(&f_readSize, "rsize", 0, "Size of reads to execute")
	flag.IntVar(&f_writeSize, "wsize", 0, "Size of writes to execute")
	flag.IntVar(&f_numWrites, "nwrites", 0, "Number of writes to execute")
	flag.IntVar(&f_numReads, "nreads", 0, "Number of reads to execute")
	flag.IntVar(&f_readOffset, "roff", 0, "Offset for each additional read")
	flag.IntVar(&f_writeOffset, "woff", 0, "Offset for each additional write")
	flag.IntVar(&f_numFiles, "nfiles", 0, "Number of different files to dispatch r/w's to")
	flag.BoolVar(&f_mixOps, "mix", false, "Mix both reads and writes at the same time")

	flag.StringVar(&f_path, "path", "", "The realpath of this file.")
	f_path += f_path + GEN_FILE_BASENAME

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

	// Change sizes to block granularity.
	f_writeSize = (f_writeSize / BLOCKSIZE) * BLOCKSIZE
	f_readSize = (f_readSize / BLOCKSIZE) * BLOCKSIZE

	// TODO
	/* Assert that we have set at least one testable configuration. */

	if file_list != "" {
		f_files = strings.Split(file_list, ",")
	}
}

func runTest() {

	var use_threads bool
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

	// Open the files to be tested.
	for _, file_name := range file_list {
		//handle, err := sut.OpenFile(file_name)
		handle, err := syscall.Open(file_name, syscall.O_RDONLY|syscall.O_DIRECT, 0)
		panic_chk(err)
		file_list_handles = append(file_list_handles, handle)
	}

	// Generate an ordering of the operations to use - not going to
	// truely limit reads and writes to their count, but will help.
	if f_mixOps {
		split := (float64(f_numReads-f_numWrites) / float64(math.MaxInt64)) * math.MaxInt64
		for i := 0; i < f_numReads+f_numWrites; i++ {
			val := rand.NormFloat64()
			ret := val >= float64(split)
			writeOrder = append(writeOrder, ret)
		}

		// Approximate the offset we will use in the file.
		// Divide first to avoid overflow.
		combinedOps := f_numReads + f_numWrites

		averageOffset = int(float64(f_readOffset)/float64(combinedOps)*float64(f_numReads) + (float64(f_writeOffset)/float64(combinedOps))*float64(f_numWrites))
	}

	//TODO: For now we use the same read/write buffer. Will need to change
	// for tests.

	wbuf := make([]byte, f_writeSize)
	rbuf := make([]byte, f_readSize)

	/* Place enough tokens for each of our active threads to take. */
	thread_collector := make(chan bool, f_threads)
	for thread := 0; thread < f_threads; thread++ {
		thread_collector <- true
	}

	/* Execute the tests. */
	if f_mixOps {

		ops_per_file := reads_per_file + writes_per_file

		for iter, handle := range file_list_handles {
			for op := 0; op < ops_per_file; op++ {
				op_index := (iter * ops_per_file) + op
				offset := (op) * averageOffset

				readNotWrite := writeOrder[op_index]
				var buffer []byte

				if readNotWrite {
					buffer = rbuf
				} else {
					buffer = wbuf
				}
				scheduleOp(handle, offset, buffer, readNotWrite, use_threads, thread_collector)
			}
		}
	} else {
		for _, handle := range file_list_handles {
			for op := 0; op < reads_per_file; op++ {
				offset := op * f_readOffset

				scheduleOp(handle, offset, rbuf, true, use_threads, thread_collector)
			}
		}
		for _, handle := range file_list_handles {
			for op := 0; op < writes_per_file; op++ {
				offset := op * f_writeOffset

				scheduleOp(handle, offset, wbuf, false, use_threads, thread_collector)
			}
		}
	}

	/* Finish the multithreaded test. */
	for thread := 0; thread < f_threads; thread++ {
		_ = <-thread_collector
	}

	/* Cleanup from the tests. */

}

func scheduleOp(file, offset int, buffer []byte, readNotWrite, threaded bool, collector chan bool) {

	execOp := func() {
		var eventType tracer.Event_t
		var op func(int, int, []byte) (int, error)

		if readNotWrite {
			op = sut.ReadAt
			eventType = tracer.T_READ
		} else {
			op = sut.WriteAt
			eventType = tracer.T_WRITE
		}

		trace := tracer.NewTraceEvent(eventType, tracer.GlobalTraceList)
		trace.Start()
		ret, err := op(file, offset, buffer)
		trace.Stop()

		panic_chk(err)
		assertF(ret == len(buffer), "ret: %v len: %v\n", ret, len(buffer))
		thread_exit(threaded, collector)
	}

	if threaded {
		thread_enter(threaded, collector)
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
		file_name := f_path + strconv.Itoa(val) + GEN_FILE_SUFFIX
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

func panic_chk(err error) {
	if err != nil {
		fmt.Println(err)
		panic("Panic!")
	}
}

func assertF(val bool, format string, message ...interface{}) {
	if !val {
		fmt.Printf(format, message)
		panic(nil)
	}
}

func assert(val bool) {
	if !val {
		panic(nil)
	}
}

func trace(id string, elapsed *int64) (string, time.Time, *int64) {
	return id, time.Now(), elapsed
}

func un(id string, start time.Time, elapsed *int64) {
	*elapsed = time.Since(start).Nanoseconds()
}

func thread_enter(threaded bool, limiter chan bool) {
	if threaded {
		<-limiter
	}
}
func thread_exit(threaded bool, limiter chan bool) {
	if threaded {
		limiter <- true
	}
}
