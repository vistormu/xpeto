package main

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{}
var headless bool

func main() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
