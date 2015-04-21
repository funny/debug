package debug

import "testing"

func Test_StackTrace(t *testing.T) {
	si := StackTrace(0)
	t.Log("\n" + si.String("    "))
}
