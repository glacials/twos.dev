package main

import (
	"github.com/spf13/cobra"
)

const (
	photoGalleryTemplatePath = "src/templates/imgcontainer.html.tmpl"
	staticAssetsDir          = "public"
	sourceDir                = "src"
)

// buildCmd represents the build command
var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build the website",
	Long:  `Build the website into dist/.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return buildAll(dst, builders, cfg)
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)
}
