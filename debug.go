package debug

import (
	"bytes"
	"fmt"
	"os"
	"runtime"
)

func Print(v ...interface{}) {
	fmt.Fprintf(os.Stderr, "[DEBUG PRINT]\nby goroutine %s\n%s\n", GoroutineID(), StackTrace(2).Bytes(""))
	fmt.Fprintf(os.Stderr, string(Dump(DumpStyle{Format: true, Indent: ""}, v...)))
}

func GoroutineID() string {
	buf := make([]byte, 15)
	buf = buf[:runtime.Stack(buf, false)]
	return string(bytes.Split(buf, []byte(" "))[1])
}

func Pause(condition bool) {
	if condition {
		fmt.Fprintf(os.Stderr, "[DEBUG PAUSE]\nby goroutine %s\n%s\n", GoroutineID(), StackTrace(2).Bytes(""))
		fmt.Fprint(os.Stderr, "press ENTER to continue\n")
		fmt.Scanln()
	}
}
