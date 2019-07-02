package extio

import (
	"bufio"
	"io"
)

func IterateReader(r io.Reader, sep byte, fn func(line []byte, lineno int) error) (err error) {
	bufr := bufio.NewReader(r)
	var line []byte
	var lineno int
	for {
		if line, err = bufr.ReadBytes(sep); err != nil {
			if err == io.EOF {
				err = nil
				if len(line) > 0 {
					err = fn(line, lineno)
				}
				return
			} else {
				return
			}
		}

		if err = fn(line, lineno); err != nil {
			return
		}

		lineno++
	}
}
