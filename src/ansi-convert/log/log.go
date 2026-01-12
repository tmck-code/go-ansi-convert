package log

import (
	"fmt"
	"os"
)

var (
	DEBUG bool = os.Getenv("DEBUG") != ""
)

func DebugFprintf(format string, a ...interface{}) {
	if DEBUG {
		fmt.Fprintf(os.Stderr, format, a...)
	}
}

func DebugFprintln(a ...interface{}) {
	if DEBUG {
		fmt.Fprintln(os.Stderr, a...)
	}
}
