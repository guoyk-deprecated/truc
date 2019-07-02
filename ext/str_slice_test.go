package ext

import (
	"strings"
	"testing"
)

func TestSanitizeStrSlice(t *testing.T) {
	if "b,c" != strings.Join(SanitizeStrSlice([]string{"   ", "b  ", " c"}), ",") {
		t.Fatal("failed")
	}
	if "" != strings.Join(SanitizeStrSlice([]string{"   "}), ",") {
		t.Fatal("failed")
	}
}
