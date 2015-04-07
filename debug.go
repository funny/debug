package debug

import (
	"bytes"
	"fmt"
	"os"
	"runtime"
)

func Print(v ...interface{}) {
	fmt.Fprintf(os.Stderr, "[DEBUG PRINT]\n%s", Dump(DumpStyle{Format: true, Indent: "  "}, v...))
	fmt.Fprintln(os.Stderr, "by goroutine ", GoroutineID())
	fmt.Fprint(os.Stderr, StackTrace(2, 0).String(""))
}

func GoroutineID() string {
	buf := make([]byte, 15)
	buf = buf[:runtime.Stack(buf, false)]
	return string(bytes.Split(buf, []byte(" "))[1])
}
