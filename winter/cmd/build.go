package main

import (
	"net/url"
	"time"

	"github.com/spf13/cobra"
	"twos.dev/winter"
)

const (
	photoGalleryTemplatePath = "src/templates/imgcontainer.html.tmpl"
	staticAssetsDir          = "public"
	sourceDir                = "src"
)

var (
	cfg winter.Config
	// buildCmd represents the build command
	buildCmd = &cobra.Command{
		Use:   "build",
		Short: "Build the website",
		Long:  `Build the website into dist/.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return buildAll(dist, builders, cfg)
		},
	}
)

func init() {
	f := buildCmd.PersistentFlags()

	cfg = winter.Config{
		AuthorEmail: *f.StringP("author-email", "e", "", "site author email"),
		AuthorName:  *f.StringP("author-name", "a", "", "site author name"),
		Desc:        *f.StringP("desc", "d", "", "site description"),
		Domain: url.URL{
			Scheme: "https",
			Host:   *f.StringP("domain", "m", "", "site root domain (e.g. twos.dev)"),
		},
		Name: *f.StringP("name", "n", "", "site name"),
		Since: *f.IntP(
			"since",
			"y",
			time.Now().Year(),
			"site year of creation (for copyright)",
		),
	}

	rootCmd.AddCommand(buildCmd)
}
