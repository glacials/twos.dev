package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"gg-scm.io/pkg/git"
	"github.com/spf13/cobra"
	"twos.dev/winter"
)

const (
	KnownURLsPath = "src/urls.txt"
)

func newFreezeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "freeze [shortname...]",
		Short: "Turn a warm file cold",
		Long: wrap(`
		Converts warm documents specified into cold documents,
		and commit to keeping them at those URLs forever.

		"Warm documents" are those that are continually and/or automatically updated by other tools,
		such as a shell script that automatically synchronizes Markdown files from your notes app.
		Warm documents are stored in src/warm.

		"Cold documents" are those that must never touched by automated tools,
		reducing the surface area for bugs in your synchronization process to mangle or delete pages.
		Cold documents are stored in src/cold.

		Conventionally, a document is born warm and remains warm while you work on it;
		when you are done or near done, you run ` + "`winter freeze <document>`" + ` to protect it from your future selves and tools.

		` + "`winter freeze`" + ` also saves ALL public HTML URLs on the generated website to a file,
		which you should commit to your repository.
		` + "`winter test`" + ` will fail if any of these HTML URLs is ever removed from dist (read: not generated).

		To perform this save step without freezing any documents,
		simply run ` + "`winter freeze`" + ` with no arguments.
	`),
		Example: wrap(`
		This command:

		    winter freeze src/warm/hello.md

		moves the file src/warm/hello.md from src/warm to src/cold,
		and adds it to the list of frozen URLs.
	`),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := winter.NewConfig()
			if err != nil {
				return err
			}
			s, err := winter.NewSubstructure(cfg)
			if err != nil {
				return err
			}

			if err := s.SaveNewURIs(dist); err != nil {
				return fmt.Errorf("cannot freeze known URIs: %w", err)
			}

			for _, shortname := range args {
				document, ok := s.DocBySourcePath(shortname)
				if !ok {
					return fmt.Errorf("cannot find document with shortname `%s`", shortname)
				}

				oldpath := document.Metadata().SourcePath
				newpath := filepath.Join("src", "cold", filepath.Base(document.Metadata().SourcePath))

				// The directory to remove warm files from, in addition to src/warm, so
				// that the file doesn't just sync back to src/warm after removal. TODO:
				// Make configurable.
				warmSourceOfTruth := "/Users/glacials/Library/Mobile Documents/27N4MQEA55~pro~writer/Documents/Published/"
				if _, err := os.Stat(warmSourceOfTruth); err != nil {
					// Allow dir to not exist, since it's hardcoded to my dir ;)
					if !os.IsNotExist(err) {
						return err
					}
				} else {
					rel, err := filepath.Rel("src/warm", oldpath)
					if err != nil {
						return err
					}
					relpath := filepath.Join(warmSourceOfTruth, rel)
					if _, err := os.Stat(relpath); err != nil {
						if !os.IsNotExist(err) {
							return err
						}
					} else if err := os.Remove(relpath); err != nil {
						return err
					}
				}

				g, err := git.New(git.Options{})
				if err != nil {
					return err
				}

				ctx := context.TODO()
				if err := os.Rename(oldpath, newpath); err != nil {
					return err
				}
				fmt.Println("git add", oldpath, newpath)
				if err := g.Add(
					ctx,
					[]git.Pathspec{git.Pathspec(oldpath), git.Pathspec(newpath)},
					git.AddOptions{},
				); err != nil {
					return err
				}

				fmt.Printf("git commit -m 'Freeze %s'\n", shortname)
				if err := g.Commit(
					ctx,
					fmt.Sprintf("Freeze %s", shortname),
					git.CommitOptions{},
				); err != nil {
					return err
				}

			}
			return nil
		},
	}
	return cmd
}
