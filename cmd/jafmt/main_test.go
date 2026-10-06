package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestRun(t *testing.T) {
	var out bytes.Buffer
	if err := run(strings.NewReader("日本語abc\n"), &out); err != nil {
		t.Fatal(err)
	}
	if got, want := out.String(), "日本語 abc\n"; got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
