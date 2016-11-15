package main

import (
	"fmt"
	blocking "github.com/spwilson2/cs758-project/scheduler-blocking"
	nonblocking "github.com/spwilson2/cs758-project/scheduler-nonblocking"
)

func main() {

	// remove these once we use the packages, get compiler errs otherwise.
	_ = blocking.TestExport
	_ = nonblocking.TestExport

	blocking.Write()

	fmt.Println("Hello world!")
}
