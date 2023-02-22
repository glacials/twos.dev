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

var freezeCmd = &cobra.Command{
	Use:   "freeze shortname...",
	Short: "Turn a warm file cold",
	Long: wrap(`
		Convert the warm document specified into a cold document. Run this when a
		document is no longer being actively updated in order to reduce exposure to
		destructive issues caused by the hotbed of src/warm.
	`),
	Example: wrap(`
		The command ` + "`" + `winter freeze DOCUMENT` + "`" + ` searches for a file
		with the shortname ` + "`" + `DOCUMENT` + "`" + ` in src/cold (whether
		explicitly in frontmatter or implicitly in filename) and moves it to
		src/cold.
	`),
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// TODO: Allow an argument to also render .md to .html
		sources, err := cmd.Flags().GetStringArray("source")
		if err != nil {
			return err
		}

		for _, shortname := range args {
			s, err := winter.NewSubstructure(winter.Config{SourceDirectories: sources})
			if err != nil {
				return err
			}

			document := s.DocByShortname(shortname)
			if document == nil || document.Document == nil {
				return fmt.Errorf("cannot find document with shortname `%s`", shortname)
			}

			oldpath := document.Source
			newpath := filepath.Join("src", "cold", filepath.Base(document.Source))

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

func init() {
	rootCmd.AddCommand(freezeCmd)
}
