package ext

import "strings"

func SanitizeStrSlice(strs []string) []string {
	var i = 0
	for i < len(strs) {
		strs[i] = strings.TrimSpace(strs[i])
		if len(strs[i]) == 0 {
			strs = append(strs[0:i], strs[i+1:]...)
		} else {
			i++
		}
	}
	return strs
}
