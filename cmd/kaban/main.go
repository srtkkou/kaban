package main

import (
	"fmt"
	"os"

	app "github.com/srtkkou/kaban"
)

var (
	version string // バージョン情報
)

func main() {
	if err := app.AgingMain(version, os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(-1)
	}
}
