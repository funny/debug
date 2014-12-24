package debug

import (
	"fmt"
	//"github.com/funny/goid"
	"os"
)

func Print(v ...interface{}) {
	fmt.Fprintln(os.Stderr, "[DEBUG PRINT]")
	Dump(os.Stderr, DumpStyle{Format: true, Indent: "  "}, v...)
	//fmt.Fprintln(os.Stderr, "by goroutine", goid.Get())
	fmt.Fprint(os.Stderr, StackTrace(2, 0).String(""))
}
