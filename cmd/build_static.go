package cmd

import (
	"fmt"

	"github.com/otiai10/copy"
)

func staticFileBuilder(staticFile, destinationDir string) error {
	if err := copy.Copy(staticFile, destinationDir); err != nil {
		return fmt.Errorf("can't build static asset %s: %w", staticFile, err)
	}

	return nil
}

func buildStaticDir(staticDir, destinationDir string) error {
	if err := copy.Copy(staticDir, destinationDir); err != nil {
		return fmt.Errorf("can't build static assets: %w", err)
	}

	return nil
}
