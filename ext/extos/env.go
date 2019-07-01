package extos

import (
	"os"
	"strings"
)

func EnvBool(out *bool, keys ...string) {
	for _, key := range keys {
		v := strings.TrimSpace(strings.ToLower(os.Getenv(key)))
		if v == "0" || strings.HasPrefix(v, "f") || strings.HasPrefix(v, "n") || v == "off" {
			*out = false
			return
		}
		if v == "1" || strings.HasPrefix(v, "t") || strings.HasPrefix(v, "y") || v == "on" {
			*out = true
			return
		}
	}
}

func EnvStr(out *string, keys ...string) {
	for _, key := range keys {
		if v := strings.TrimSpace(os.Getenv(key)); len(v) > 0 {
			*out = v
			return
		}
	}
}
