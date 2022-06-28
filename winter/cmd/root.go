package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var (
	rootCmd = &cobra.Command{
		Use:   "winter",
		Short: "Build or serve a static website locally",
		Long:  `Build or serve a static website from source.`,
	}
	debug *bool
)

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	debug = rootCmd.PersistentFlags().
		Bool("debug", false, "output to dist/debug/ results of each transformation")
}
