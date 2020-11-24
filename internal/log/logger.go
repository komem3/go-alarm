package log

import (
	"fmt"
	"os"
)

var verbose bool

func SetVerbose(v bool) {
	verbose = v
}

func Printf(format string, args ...interface{}) {
	if verbose {
		fmt.Fprintf(os.Stderr, format, args...)
	}
}
