// Package cmd contains the commands for the winter CLI.
package cmd

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
		Version: "0.1.0",
	}
)

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	f := rootCmd.PersistentFlags()

	_ = *f.StringArrayP("source", "i", []string{}, "supplemental source file or directory to build (can be specified multiple times)")
}