package debug

import (
	"testing"
)

type MyData struct {
	IntField   int
	FloatField float64
	StrField   string
	MapField   map[int]string
	SliceField []int
	PointField *MyData
}

func Test_Dump(t *testing.T) {
	data := &MyData{
		1234,
		77.88,
		"xyz",
		map[int]string{
			1: "abc",
			2: "def",
			3: "ghi",
		},
		[]int{
			3,
			7,
			11,
			13,
			17,
		},
		nil,
	}
	data.PointField = data

	t.Logf("\n%s", Dump(DumpStyle{Pointer: true, Indent: "    "}, data))
	t.Logf("\n%s", Dump(DumpStyle{Format: true, Indent: "    "}, data))
}
