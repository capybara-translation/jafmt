// jafmt は stdin のテキストを整形して stdout に出力する。
package main

import (
	"fmt"
	"io"
	"os"

	"github.com/capybara-translation/jafmt/formatter"
)

const usageText = `Usage:
  jafmt < input.txt     stdin のテキストを整形して stdout に出力する
  jafmt version         バージョンを表示する
`

func main() {
	os.Exit(run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr))
}

// run は引数に応じて処理し、終了コードを返す。テストから呼べるよう os.Exit は呼ばない。
func run(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	if len(args) > 0 {
		switch args[0] {
		case "help", "-h", "--help":
			fmt.Fprint(stdout, usageText)
			return 0
		case "version", "--version":
			fmt.Fprintf(stdout, "jafmt %s\n", currentVersion())
			return 0
		default:
			fmt.Fprintf(stderr, "jafmt: unknown argument %q\n", args[0])
			fmt.Fprint(stderr, usageText)
			return 2
		}
	}
	if err := format(stdin, stdout); err != nil {
		fmt.Fprintln(stderr, "jafmt:", err)
		return 1
	}
	return 0
}

func format(r io.Reader, w io.Writer) error {
	in, err := io.ReadAll(r)
	if err != nil {
		return err
	}
	_, err = io.WriteString(w, formatter.Format(string(in)))
	return err
}
