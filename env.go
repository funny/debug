package debug

import (
	"os"
	"strings"
)

var _GODEBUG_ITEMS_ = make(map[string]string)

func init() {
	if v := os.Getenv("GODEBUG"); v != "" {
		items := strings.Split(v, ",")
		for _, item := range items {
			pair := strings.Split(item, "=")
			if len(pair) != 2 {
				continue
			}
			_GODEBUG_ITEMS_[pair[0]] = pair[1]
		}
	}
}

// Get setting value in GODEBUG variable by name.
// If setting not exists, returns defaultValue.
// See GODEBUG environment variable description in http://golang.org/pkg/runtime/
func GODEBUG(name, defaultValue string) string {
	if v, ok := _GODEBUG_ITEMS_[name]; ok {
		return v
	}
	return defaultValue
}
