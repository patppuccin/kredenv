package main

import (
	"os"

	"github.com/patppuccin/kredenv/cmd"
)

func main() {
	if err := cmd.KredEnvCmd.Execute(); err != nil {
		os.Stderr.WriteString("Error running the program: " + err.Error())
		os.Exit(1)
	}
}
