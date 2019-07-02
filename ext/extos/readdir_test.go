package extos

import "testing"

func TestReaddirFiles(t *testing.T) {
	err := ReaddirFiles("/tmp/test", ReaddirFilesOptions{
		BeforeFile: func(name string) bool {
			t.Log("before", name)
			return true
		},
		Handle: func(buf []byte, name string) error {
			t.Log("buf", string(buf), name)
			return nil
		},
		AfterFile: func(name string) {
			t.Log("after", name)
		},
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestReaddirLines(t *testing.T) {
	err := ReaddirLines("/tmp/test", ReaddirLinesOptions{
		BeforeFile: func(name string) bool {
			t.Log("before", name)
			return true
		},
		Handle: func(line []byte, name string, lineno int) error {
			t.Log(name, lineno, string(line))
			return nil
		},
		AfterFile: func(name string) {
			t.Log("after", name)
		},
	})
	if err != nil {
		t.Fatal(err)
	}
}
