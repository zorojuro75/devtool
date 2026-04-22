package main

import "github.com/zorojuro75/devtool/cmd"

var (
	version   = "dev"
	buildDate = "unknown"
	commit    = "unknown"
)

func main() {
	cmd.Execute(version, buildDate, commit)
}