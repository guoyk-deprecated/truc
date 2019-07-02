package extos

import (
	"bytes"
	"github.com/yankeguo/truc/ext/extio"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
)

type ReaddirFilesOptions struct {
	BeforeFile func(name string) bool
	Handle     func(buf []byte, name string) error
	AfterFile  func(name string)
}

func ReaddirFiles(dir string, opt ReaddirFilesOptions) (err error) {
	var d *os.File
	if d, err = os.Open(dir); err != nil {
		return
	}
	var names []string
	if names, err = d.Readdirnames(-1); err != nil {
		return
	}
	sort.Sort(sort.StringSlice(names))
	for _, name := range names {
		if opt.BeforeFile != nil && !opt.BeforeFile(name) {
			continue
		}
		var buf []byte
		if buf, err = ioutil.ReadFile(filepath.Join(dir, name)); err != nil {
			return
		}
		if opt.Handle != nil {
			if err = opt.Handle(buf, name); err != nil {
				return
			}
		}
		if opt.AfterFile != nil {
			opt.AfterFile(name)
		}
	}
	return
}

type ReaddirLinesOptions struct {
	BeforeFile func(name string) bool
	Handle     func(line []byte, name string, lineno int) error
	AfterFile  func(name string)
}

func ReaddirLines(dir string, opt ReaddirLinesOptions) (err error) {
	var d *os.File
	if d, err = os.Open(dir); err != nil {
		return
	}
	var names []string
	if names, err = d.Readdirnames(-1); err != nil {
		return
	}
	sort.Sort(sort.StringSlice(names))
	for _, name := range names {
		if opt.BeforeFile != nil && !opt.BeforeFile(name) {
			continue
		}
		var f *os.File
		if f, err = os.Open(filepath.Join(dir, name)); err != nil {
			return
		}

		if opt.Handle != nil {
			if err = extio.IterateReader(f, '\n', func(line []byte, lineno int) (err error) {
				if err = opt.Handle(bytes.TrimSpace(line), name, lineno); err != nil {
					return
				}
				return
			}); err != nil {
				_ = f.Close()
				return
			}
		}

		_ = f.Close()
		if opt.AfterFile != nil {
			opt.AfterFile(name)
		}
	}
	return
}
