package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"regexp"
	"syscall"

	"github.com/spf13/cobra"
	"twos.dev/winter"
)

const (
	photoGalleryTemplatePath = "src/templates/imgcontainer.html.tmpl"
	port                     = 8100
	serveFlag                = "serve"
	sourceDir                = "src"
	staticAssetsDir          = "public"
)

// Builder is a function that builds a source file src into a destination
// directory dst.
type Builder func(src, dst string, cfg winter.Config) error

var (
	authorPattern = regexp.MustCompile(`^(.*) <(.*)>$`)
	buildCmd      = &cobra.Command{
		Use:   "build",
		Short: "Build the website",
		Long:  `Build the website into dist/.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := cmdConfig(cmd)
			if err != nil {
				return err
			}

			s, err := winter.NewSubstructure(cfg)
			if err != nil {
				return err
			}

			if err := s.ExecuteAll(dist, cfg); err != nil {
				return err
			}

			serve, err := cmd.Flags().GetBool(serveFlag)
			if err != nil {
				return err
			}

			if serve {
				stop := make(chan struct{})
				reloader := Reloader{
					Ignore:       ignoreDirectories,
					Substructure: s,
				}

				var mux http.ServeMux
				mux.Handle("/", http.FileServer(http.Dir(dist)))
				mux.Handle("/ws", reloader.Handler())

				server := http.Server{
					Addr:    fmt.Sprintf(":%d", port),
					Handler: &mux,
				}
				server.RegisterOnShutdown(reloader.Shutdown)

				go listenForCtrlC(stop, &server, &reloader)
				go startFileServer(&server)

				if err := reloader.Watch(append(cfg.SourceDirectories, ".")); err != nil {
					return err
				}

				log.Printf("Serving %s on http://localhost:%d\n", dist, port)
				<-stop
				return nil
			}

			return nil
		},
	}
	ignoreFiles = map[string]struct{}{
		"README.md": {},
		".DS_Store": {},
	}
	ignoreDirectories = map[string]struct{}{
		".git":         {},
		".github":      {},
		dist:           {},
		"node_modules": {},
	}
	serve bool
)

func init() {
	// TODO: Allow all these options to be in a config file
	f := buildCmd.PersistentFlags()

	_ = *f.BoolP(serveFlag, "s", false, "start a webserver and rebuild on file changes")
	_ = *f.IntP("since", "y", 0, "site year of creation (e.g. 2021)")
	_ = *f.StringP("author", "a", "", "author (e.g. Benjamin Carlsson <ben@twos.dev>)")
	_ = *f.StringP("desc", "d", "", "site description (e.g. misc thoughts)")
	_ = *f.StringP("domain", "m", "", "site root domain (e.g. twos.dev)")
	_ = *f.StringP("name", "n", "", "site name (e.g. twos.dev)")

	rootCmd.AddCommand(buildCmd)
}

func listenForCtrlC(stop chan struct{}, srvr *http.Server, reloader *Reloader) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	log.Println("Stopping")
	if err := srvr.Shutdown(context.TODO()); err != nil {
		log.Fatal(err)
	}
	stop <- struct{}{}
}

func startFileServer(server *http.Server) {
	if err := server.ListenAndServe(); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			// Expected when we call server.Shutdown()
			return
		}
		log.Fatal(fmt.Errorf("can't listen and serve: %w", err))
	}
}

func cmdConfig(cmd *cobra.Command) (winter.Config, error) {
	author := cmd.Flag("author").Value.String()
	if author == "" {
		author = "Anonymous <unspecified>"
	}
	authorParts := authorPattern.FindStringSubmatch(author)
	if len(authorParts) != 3 {
		fmt.Println("invalid --author format: must be \"NAME <EMAIL>\"")
		os.Exit(1)
	}

	since, err := cmd.Flags().GetInt("since")
	if err != nil {
		return winter.Config{}, err
	}

	sources, err := cmd.Flags().GetStringArray("source")
	if err != nil {
		return winter.Config{}, err
	}

	return winter.Config{
		AuthorName:  authorParts[1],
		AuthorEmail: authorParts[2],
		Desc:        cmd.Flag("desc").Value.String(),
		Dist:        "dist",
		Domain: url.URL{
			Scheme: "https",
			Host:   cmd.Flag("domain").Value.String(),
		},
		Name:              cmd.Flag("name").Value.String(),
		Since:             since,
		SourceDirectories: sources,
	}, nil
}
