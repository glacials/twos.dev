package cmd

import (
	"fmt"

	"github.com/otiai10/copy"
)

func buildStatic(staticDir, destinationDir string) error {
	if err := copy.Copy(staticDir, destinationDir); err != nil {
		return fmt.Errorf("can't build static assets: %w", err)
	}

	return nil
}
