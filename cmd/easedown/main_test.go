package main

import "testing"

func TestSeparator(t *testing.T) {
	t.Log(separator.Split("abcdefg@hotmail.com----", -1))
}
