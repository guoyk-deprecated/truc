package cmdutil

import (
	"os"
	"strings"
)

func EnvBool(out *bool, key string) {
	v := strings.TrimSpace(strings.ToLower(os.Getenv(key)))
	if len(v) == 0 {
		return
	}
	if v == "0" || strings.HasPrefix(v, "f") || strings.HasPrefix(v, "n") || v == "off" {
		*out = false
	}
	if v == "1" || strings.HasPrefix(v, "t") || strings.HasPrefix(v, "y") || v == "on" {
		*out = true
	}
}

func EnvStr(out *string, key string) {
	v := strings.TrimSpace(os.Getenv(key))
	if len(v) == 0 {
		return
	}
	*out = v
}
