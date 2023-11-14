package cmd

import (
	"fmt"
	"os"

	"github.com/adrg/xdg"
	"github.com/spf13/cobra"
	"twos.dev/winter"
)

func newCleanCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "clean",
		Short: "Purge the Winter cache",
		Long: wrap(`
		Purges the internal Winter cache and the current directory's generated site.

		This should be used instead of manually removing dist/,
		because Winter stores some cache files for large images internally.
	`),
		Args: cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			cache, err := xdg.CacheFile(winter.AppName)
			if err != nil {
				return fmt.Errorf("cannot find cache: %w", err)
			}
			if err := os.RemoveAll(cache); err != nil {
				return fmt.Errorf("cannot purge cache: %w", err)
			}
			if err := os.RemoveAll(dist); err != nil {
				return fmt.Errorf("cannot purge website: %w", err)
			}
			return nil
		},
	}
	return cmd
}
