package debug

import "testing"

func Test_Print(t *testing.T) {
	Print("abc", 123)
}

func Test_Pause(t *testing.T) {
	Pause(true)
}
