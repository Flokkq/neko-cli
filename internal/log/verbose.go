package log

import "fmt"

var Verbose = false

func V(msg string, args ...interface{}) {
	if Verbose {
		fmt.Printf(msg+"\n", args...)
	}
}
