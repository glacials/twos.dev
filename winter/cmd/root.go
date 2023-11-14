// Package cmd contains the commands for the winter CLI.
package cmd

import (
	"context"
	"os"

	"github.com/spf13/cobra"
)

const (
	dist = "dist"
)

// Execute sets up the root command and all attached subcommands,
// then runs them according to the CLI arguments supplied.
func Execute() {
	rootCmd := &cobra.Command{
		Use:     "winter",
		Short:   "Build or serve a static website locally",
		Long:    `Build or serve a static website from source.`,
		Version: "0.1.0",
	}
	f := rootCmd.PersistentFlags()
	_ = *f.StringArrayP("source", "i", []string{}, "supplemental source file or directory to build (can be specified multiple times)")
	rootCmd.AddCommand(newBuildCmd())
	rootCmd.AddCommand(newCleanCmd())
	rootCmd.AddCommand(newConfigCmd())
	rootCmd.AddCommand(newFreezeCmd())
	rootCmd.AddCommand(newInitCmd())
	rootCmd.AddCommand(newServeCmd())
	rootCmd.AddCommand(newTestCmd())
	err := rootCmd.ExecuteContext(context.Background())
	if err != nil {
		os.Exit(1)
	}
}

func init() {
}
