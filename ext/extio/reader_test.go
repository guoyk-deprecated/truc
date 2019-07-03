package extio

import (
	"bytes"
	"github.com/pkg/errors"
	"testing"
)

func TestIterateReader(t *testing.T) {
	var err error
	r := bytes.NewBufferString("line 1\nline 2\n\nline3\n")
	s := []string{"line 1\n", "line 2\n", "\n", "line3\n"}
	if err = IterateReader(r, '\n', func(line []byte, lineno int) error {
		if string(line) != s[lineno] {
			t.Fatalf("line %d", lineno)
		}
		if lineno > 3 {
			t.Fatal("extra line")
		}
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	if err = IterateReader(bytes.NewBufferString(" "), '\n', func(line []byte, lineno int) error {
		if string(line) != " " {
			t.Fatal("not equal")
		}
		return errors.New("dummy error")
	}); err == nil {
		t.Fatal("missing err")
	}
}
