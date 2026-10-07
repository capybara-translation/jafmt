package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestRun_Format(t *testing.T) {
	var stdout, stderr bytes.Buffer
	if code := run(nil, strings.NewReader("日本語abc\n"), &stdout, &stderr); code != 0 {
		t.Fatalf("exit code = %d, want 0 (stderr: %q)", code, stderr.String())
	}
	if got, want := stdout.String(), "日本語 abc\n"; got != want {
		t.Errorf("stdout = %q, want %q", got, want)
	}
}

func TestRun_Version(t *testing.T) {
	for _, arg := range []string{"version", "--version"} {
		var stdout, stderr bytes.Buffer
		if code := run([]string{arg}, strings.NewReader(""), &stdout, &stderr); code != 0 {
			t.Errorf("%s: exit code = %d, want 0", arg, code)
		}
		if !strings.HasPrefix(stdout.String(), "jafmt ") {
			t.Errorf("%s: stdout = %q, want \"jafmt <version>\"", arg, stdout.String())
		}
	}
}

func TestRun_Help(t *testing.T) {
	for _, arg := range []string{"help", "-h", "--help"} {
		var stdout, stderr bytes.Buffer
		if code := run([]string{arg}, strings.NewReader(""), &stdout, &stderr); code != 0 {
			t.Errorf("%s: exit code = %d, want 0", arg, code)
		}
		if !strings.HasPrefix(stdout.String(), "Usage:") {
			t.Errorf("%s: stdout = %q, want usage", arg, stdout.String())
		}
	}
}

// 未知の引数では stdin を読まずに終了する。
// 誤った呼び出しで入力を消費したり、整形結果を出力したりしないようにする。
func TestRun_UnknownArgument(t *testing.T) {
	var stdout, stderr bytes.Buffer
	if code := run([]string{"file.txt"}, strings.NewReader("日本語abc"), &stdout, &stderr); code != 2 {
		t.Errorf("exit code = %d, want 2", code)
	}
	if stdout.Len() != 0 {
		t.Errorf("stdout = %q, want empty", stdout.String())
	}
	if !strings.Contains(stderr.String(), "Usage:") {
		t.Errorf("stderr = %q, want usage", stderr.String())
	}
}
