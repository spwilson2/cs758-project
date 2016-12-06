package nonblocking

import (
	"log"
	"os"
)

// Hack to get constant errors.
type Error string

func (e Error) Error() string { return string(e) }

/* easily check for errors and panic */
func chk_err(err error) {
	if err != nil {
		log.Printf("(FAILED) %s\n", os.Args[0])
		log.Printf("%v \n", err)
		panic(err) //os.Exit(-1)
	}
}
