package cmd

import "github.com/spf13/cobra"

func newServeCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "serve",
		Short: "Alias of build --serve",
		Long: wrap(`
			Build, serve, and continually rebuild the website.

			Alias of winter build --serve.
			Builds the website, starts a web server,
			and monitors the filesystem for changes.
			When a change occurs, the website is selectively rebuilt
			according to the scope of the change.

			For example, if a .go file is changed, the web server is taken down,
			rebuilt, and booted back up;
			while if a .html.tmpl file is changed, that template is re-executed
			along with any related page that may have changed as a result.
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			buildCmd := newBuildCmd()
			if err := buildCmd.Flag(serveFlag).Value.Set("true"); err != nil {
				return err
			}
			return buildCmd.RunE(cmd, args)
		},
	}
}
