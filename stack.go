package debug

import (
	"bytes"
	"fmt"
	"runtime"
)

var (
	dunno     = []byte("???")
	dot       = []byte(".")
	centerDot = []byte("·")
	slash     = []byte("/")
)

// See runtime/debug.Stack()
func StackTrace(begin, end int) StackInfo {
	si := StackInfo(make([]uintptr, 0, 5))
	pc := make([]uintptr, 5)
	for i := begin + 1; end == 0 || i < end; {
		n := runtime.Callers(i, pc)
		if n == 0 {
			break
		}
		i += n
		si = append(si, pc[0:n]...)
	}
	return si
}

type StackInfo []uintptr

func (si StackInfo) Bytes(indent string) []byte {
	var buf = new(bytes.Buffer)
	for i := 0; i < len(si); i++ {
		buf.WriteString(indent)
		name, file, line := function(si[i])
		fmt.Fprintf(buf, "at %s() [%s:%d]\n", name, file, line)
	}
	return buf.Bytes()
}

func (si StackInfo) String(indent string) string {
	return string(si.Bytes(indent))
}

// function returns, if possible, the name of the function containing the PC.
func function(pc uintptr) (name []byte, file string, line int) {
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return dunno, "???", 0
	}
	file, line = fn.FileLine(pc)
	name = []byte(fn.Name())
	// The name includes the path name to the package, which is unnecessary
	// since the file name is already included.  Plus, it has center dots.
	// That is, we see
	//	runtime/debug.*T·ptrmethod
	// and want
	//	*T.ptrmethod
	// Since the package path might contains dots (e.g. code.google.com/...),
	// we first remove the path prefix if there is one.
	if lastslash := bytes.LastIndex(name, slash); lastslash >= 0 {
		name = name[lastslash+1:]
	}
	if period := bytes.Index(name, dot); period >= 0 {
		name = name[period+1:]
	}
	name = bytes.Replace(name, centerDot, dot, -1)
	return
}
