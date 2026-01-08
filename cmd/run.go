package main

import (
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run <path>",
	Short: "runs the app with go run",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path := args[0]

		goArgs := []string{"run"}
		goArgs = append(goArgs, path)

		c := exec.Command("go", goArgs...)
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		c.Stdin = os.Stdin
		c.Env = os.Environ()

		return c.Run()
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
