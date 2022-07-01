// Winter is a static website generator that allows for easy extensibility. The
// generator itself is a lightweight framework that only knows how to process
// transformations.
//
// A transformation is a Go function that receives a document as input, applies
// some modification to it, and returns the result. A Winter configuration is
// defined as a list of transformations.
//
// Some examples of transformations are converting Markdown to HTML, scraping a
// piece of frontmatter from the document, or executing the document as a
// template.
package main

import (
	"net/url"
	"os"
	"time"

	"github.com/spf13/cobra"
	"twos.dev/winter"
)

// rootCmd represents the base command when called without any subcommands
var (
	rootCmd = &cobra.Command{
		Use:   "winter",
		Short: "Build or serve a static website locally",
		Long:  `Build or serve a static website from source.`,
	}

	cfg winter.Config

	authorEmail string
	authorName  string
	debug       bool
	desc        string
	domain      string
	name        string
	since       string
)

func main() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	f := *buildCmd.PersistentFlags()

	d := f.Bool("debug", false, "output results of transformations to dist/debug")
	if d != nil {
		debug = *d
	}

	cfg = winter.Config{
		AuthorEmail: *f.StringP("author-email", "e", "", "site author email"),
		AuthorName:  *f.StringP("author-name", "a", "", "site author name"),
		Desc:        *f.StringP("desc", "d", "", "site description"),
		Domain: url.URL{
			Scheme: "https",
			Host:   *f.StringP("domain", "h", "", "site root domain (e.g. twos.dev)"),
		},
		Name: *f.StringP("name", "n", "", "site name"),
		Since: *f.IntP(
			"since",
			"y",
			time.Now().Year(),
			"site year of creation (for copyright)",
		),
	}
}
