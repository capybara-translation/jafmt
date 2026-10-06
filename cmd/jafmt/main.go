// jafmt は stdin のテキストを整形して stdout に出力する。
package main

import (
	"fmt"
	"io"
	"os"

	"jafmt/formatter"
)

func main() {
	if err := run(os.Stdin, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, "jafmt:", err)
		os.Exit(1)
	}
}

func run(r io.Reader, w io.Writer) error {
	in, err := io.ReadAll(r)
	if err != nil {
		return err
	}
	_, err = io.WriteString(w, formatter.Format(string(in)))
	return err
}
