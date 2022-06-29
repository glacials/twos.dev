package main

import (
	"fmt"

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
		if err := buildAll(dst, builders, Config{Debug: *debug}); err != nil {
			return fmt.Errorf("can't build the world: %w", err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)
}
