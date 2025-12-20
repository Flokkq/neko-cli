package log

import "fmt"

const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBold   = "\033[1m"
)

var Verbose = false

func V(msg string, args ...interface{}) {
	if Verbose {
		fmt.Printf(msg+"\n", args...)
	}
}
