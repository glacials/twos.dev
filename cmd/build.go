package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const (
	photoGalleryTemplatePath = "src/templates/imgcontainer.html"
	staticAssetsDir          = "public"
	sourceDir                = "src"
)

// buildCmd represents the build command
var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build twos.dev",
	Long:  `Build twos.dev into dist/.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := buildTheWorld(); err != nil {
			return fmt.Errorf("can't build the world: %w", err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)
}
