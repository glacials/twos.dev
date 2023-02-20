// Package cmd contains the commands for the winter CLI.
package main

import (
	"os"

	"github.com/spf13/cobra"
)

const (
	dist = "dist"
)

// rootCmd represents the base command when called without any subcommands
var (
	rootCmd = &cobra.Command{
		Use:     "winter",
		Short:   "Build or serve a static website locally",
		Long:    `Build or serve a static website from source.`,
		Version: "0.0.1",
	}
)

func main() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
}
