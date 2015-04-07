package debug

import (
	"bytes"
	"fmt"
	"reflect"
)

// Dump style.
type DumpStyle struct {
	HeadLen      int
	Pointer      bool
	Format       bool
	Indent       string
	StructFilter func(string, string) bool
}

type pointerInfo struct {
	prev *pointerInfo
	n    int
	addr uintptr
	pos  int
	used []int
}

// Dump data.
func Dump(style DumpStyle, data ...interface{}) []byte {
	buff := new(bytes.Buffer)
	for _, v := range data {
		var (
			pointers   *pointerInfo
			interfaces = make([]reflect.Value, 0, 10)
		)

		printKeyValue(buff, reflect.ValueOf(v), &pointers, &interfaces, style, 1)

		fmt.Fprintln(buff)

		if style.Pointer && pointers != nil {
			printPointerInfo(buff, style.HeadLen, pointers)
		}
	}
	return buff.Bytes()
}

func isSimpleType(val reflect.Value, kind reflect.Kind, pointers **pointerInfo, interfaces *[]reflect.Value) bool {
	switch kind {
	case reflect.Bool:
		return true
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return true
	case reflect.Uint8, reflect.Uint16, reflect.Uint, reflect.Uint32, reflect.Uint64:
		return true
	case reflect.Float32, reflect.Float64:
		return true
	case reflect.Complex64, reflect.Complex128:
		return true
	case reflect.String:
		return true
	case reflect.Chan:
		return true
	case reflect.Invalid:
		return true
	case reflect.Interface:
		for _, in := range *interfaces {
			if reflect.DeepEqual(in, val) {
				return true
			}
		}
		return false
	case reflect.UnsafePointer:
		if val.IsNil() {
			return true
		}

		var elem = val.Elem()

		if isSimpleType(elem, elem.Kind(), pointers, interfaces) {
			return true
		}

		var addr = val.Elem().UnsafeAddr()

		for p := *pointers; p != nil; p = p.prev {
			if addr == p.addr {
				return true
			}
		}

		return false
	}

	return false
}

func printKeyValue(buf *bytes.Buffer, val reflect.Value, pointers **pointerInfo, interfaces *[]reflect.Value, style DumpStyle, level int) {
	var t = val.Kind()

	switch t {
	case reflect.Bool:
		fmt.Fprint(buf, val.Bool())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		fmt.Fprint(buf, val.Int())
	case reflect.Uint8, reflect.Uint16, reflect.Uint, reflect.Uint32, reflect.Uint64:
		fmt.Fprint(buf, val.Uint())
	case reflect.Float32, reflect.Float64:
		fmt.Fprint(buf, val.Float())
	case reflect.Complex64, reflect.Complex128:
		fmt.Fprint(buf, val.Complex())
	case reflect.UnsafePointer:
		fmt.Fprintf(buf, "unsafe.Pointer(0x%X)", val.Pointer())
	case reflect.Ptr:
		if val.IsNil() {
			fmt.Fprint(buf, "nil")
			return
		}

		var addr = val.Elem().UnsafeAddr()

		for p := *pointers; p != nil; p = p.prev {
			if addr == p.addr {
				p.used = append(p.used, buf.Len())
				fmt.Fprintf(buf, "0x%X", addr)
				return
			}
		}

		*pointers = &pointerInfo{
			prev: *pointers,
			addr: addr,
			pos:  buf.Len(),
			used: make([]int, 0),
		}

		fmt.Fprint(buf, "&")

		printKeyValue(buf, val.Elem(), pointers, interfaces, style, level)
	case reflect.String:
		fmt.Fprint(buf, "\"", val.String(), "\"")
	case reflect.Interface:
		var value = val.Elem()

		if !value.IsValid() {
			fmt.Fprint(buf, "nil")
		} else {
			for _, in := range *interfaces {
				if reflect.DeepEqual(in, val) {
					fmt.Fprint(buf, "repeat")
					return
				}
			}

			*interfaces = append(*interfaces, val)

			printKeyValue(buf, value, pointers, interfaces, style, level+1)
		}
	case reflect.Struct:
		var t = val.Type()

		fmt.Fprint(buf, t)
		fmt.Fprint(buf, "{")

		for i := 0; i < val.NumField(); i++ {
			if style.Format {
				fmt.Fprintln(buf)
			} else {
				fmt.Fprint(buf, " ")
			}

			var name = t.Field(i).Name

			if style.Format {
				for ind := 0; ind < level; ind++ {
					fmt.Fprint(buf, style.Indent)
				}
			}

			fmt.Fprint(buf, name)
			fmt.Fprint(buf, ": ")

			if style.StructFilter != nil && style.StructFilter(t.String(), name) {
				fmt.Fprint(buf, "ignore")
			} else {
				printKeyValue(buf, val.Field(i), pointers, interfaces, style, level+1)
			}

			fmt.Fprint(buf, ",")
		}

		if style.Format {
			fmt.Fprintln(buf)

			for ind := 0; ind < level-1; ind++ {
				fmt.Fprint(buf, style.Indent)
			}
		} else {
			fmt.Fprint(buf, " ")
		}

		fmt.Fprint(buf, "}")
	case reflect.Array, reflect.Slice:
		fmt.Fprint(buf, val.Type())
		fmt.Fprint(buf, "{")

		var allSimple = true

		for i := 0; i < val.Len(); i++ {
			var elem = val.Index(i)

			var isSimple = isSimpleType(elem, elem.Kind(), pointers, interfaces)

			if !isSimple {
				allSimple = false
			}

			if style.Format && !isSimple {
				fmt.Fprintln(buf)
			} else {
				fmt.Fprint(buf, " ")
			}

			if style.Format && !isSimple {
				for ind := 0; ind < level; ind++ {
					fmt.Fprint(buf, style.Indent)
				}
			}

			printKeyValue(buf, elem, pointers, interfaces, style, level+1)

			if i != val.Len()-1 || !allSimple {
				fmt.Fprint(buf, ",")
			}
		}

		if style.Format && !allSimple {
			fmt.Fprintln(buf)

			for ind := 0; ind < level-1; ind++ {
				fmt.Fprint(buf, style.Indent)
			}
		} else {
			fmt.Fprint(buf, " ")
		}

		fmt.Fprint(buf, "}")
	case reflect.Map:
		var t = val.Type()
		var keys = val.MapKeys()

		fmt.Fprint(buf, t)
		fmt.Fprint(buf, "{")

		var allSimple = true

		for i := 0; i < len(keys); i++ {
			var elem = val.MapIndex(keys[i])

			var isSimple = isSimpleType(elem, elem.Kind(), pointers, interfaces)

			if !isSimple {
				allSimple = false
			}

			if style.Format && !isSimple {
				fmt.Fprintln(buf)
			} else {
				fmt.Fprint(buf, " ")
			}

			if style.Format && !isSimple {
				for ind := 0; ind <= level; ind++ {
					fmt.Fprint(buf, style.Indent)
				}
			}

			printKeyValue(buf, keys[i], pointers, interfaces, style, level+1)
			fmt.Fprint(buf, ": ")
			printKeyValue(buf, elem, pointers, interfaces, style, level+1)

			if i != val.Len()-1 || !allSimple {
				fmt.Fprint(buf, ",")
			}
		}

		if style.Format && !allSimple {
			fmt.Fprintln(buf)

			for ind := 0; ind < level-1; ind++ {
				fmt.Fprint(buf, style.Indent)
			}
		} else {
			fmt.Fprint(buf, " ")
		}

		fmt.Fprint(buf, "}")
	case reflect.Chan:
		fmt.Fprint(buf, val.Type())
	case reflect.Invalid:
		fmt.Fprint(buf, "invalid")
	default:
		fmt.Fprint(buf, "unknow")
	}
}

func printPointerInfo(buf *bytes.Buffer, headlen int, pointers *pointerInfo) {
	var anyused = false
	var pointerNum = 0

	for p := pointers; p != nil; p = p.prev {
		if len(p.used) > 0 {
			anyused = true
		}
		pointerNum += 1
		p.n = pointerNum
	}

	if anyused {
		var pointerBufs = make([][]rune, pointerNum+1)

		for i := 0; i < len(pointerBufs); i++ {
			var pointerBuf = make([]rune, buf.Len()+headlen)

			for j := 0; j < len(pointerBuf); j++ {
				pointerBuf[j] = ' '
			}

			pointerBufs[i] = pointerBuf
		}

		for pn := 0; pn <= pointerNum; pn++ {
			for p := pointers; p != nil; p = p.prev {
				if len(p.used) > 0 && p.n >= pn {
					if pn == p.n {
						pointerBufs[pn][p.pos+headlen] = '└'

						var maxpos = 0

						for i, pos := range p.used {
							if i < len(p.used)-1 {
								pointerBufs[pn][pos+headlen] = '┴'
							} else {
								pointerBufs[pn][pos+headlen] = '┘'
							}

							maxpos = pos
						}

						for i := 0; i < maxpos-p.pos-1; i++ {
							if pointerBufs[pn][i+p.pos+headlen+1] == ' ' {
								pointerBufs[pn][i+p.pos+headlen+1] = '─'
							}
						}
					} else {
						pointerBufs[pn][p.pos+headlen] = '│'

						for _, pos := range p.used {
							if pointerBufs[pn][pos+headlen] == ' ' {
								pointerBufs[pn][pos+headlen] = '│'
							} else {
								pointerBufs[pn][pos+headlen] = '┼'
							}
						}
					}
				}
			}

			buf.WriteString(string(pointerBufs[pn]) + "\n")
		}
	}
}
