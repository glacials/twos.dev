package cmd

import (
	"embed"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

//go:embed all:defaults
var defaults embed.FS

func newInitCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize a new Winter project",
		Long: wrap(`
		Create a new Winter project in the current directory.


		This creates the necessary directory structure, configuration, and
		a simplistic set of templates for minimal operation.
	`),
		Args: cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			return fs.WalkDir(
				defaults,
				".",
				func(path string, d fs.DirEntry, err error) error {
					if err != nil {
						return err
					}
					if d.IsDir() {
						return nil
					}

					relpath, err := filepath.Rel("defaults", path)
					if err != nil {
						return err
					}

					source, err := defaults.Open(path)
					if err != nil {
						return err
					}
					defer source.Close()

					if err := os.MkdirAll(filepath.Dir(relpath), 0o755); err != nil {
						return err
					}

					dest, err := os.Create(relpath)
					if err != nil {
						return err
					}
					defer dest.Close()

					_, err = io.Copy(dest, source)
					if err != nil {
						return err
					}

					return nil
				},
			)
		},
	}
	return cmd
}
