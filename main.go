package main

import (
	"os"

	"github.com/andrewvc/stockbeat/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
