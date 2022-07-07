package main

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
	Use:   "freeze shortname",
	Short: "Turn a warm file cold",
	Long: wrap(`
		Convert the warm document specified into a cold document. Run this when a
		document is no longer being actively updated in order to reduce exposure to
		destructive issues caused by the hotbed of src/warm.
	`),
	Example: wrap(`
		The command ` + "`" + `winter freeze autism` + "`" + ` searches for a file
		with the shortname ` + "`" + `autism` + "`" + ` in src/cold (whether
		explicitly in frontmatter or implicitly in filename) and moves it to
		src/cold.
	`),
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		shortname := args[0]
		s, err := winter.NewSubstructure(winter.Config{})
		if err != nil {
			return err
		}

		document := s.DocByShortname(shortname)
		if document == nil {
			return fmt.Errorf("cannot find document with shortname `%s`", shortname)
		}

		oldpath := document.SrcPath
		newpath := filepath.Join("src", "cold", filepath.Base(document.SrcPath))

		if oldpath == newpath {
			fmt.Printf("%s is already frozen (%s)\n", shortname, oldpath)
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

		return nil
	},
}

func init() {
	rootCmd.AddCommand(freezeCmd)
}
