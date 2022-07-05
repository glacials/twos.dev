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

var (
	authorPattern = regexp.MustCompile(`^(.*) <(.*)>$`)
	cfg           winter.Config
	buildCmd      = &cobra.Command{
		Use:   "build",
		Short: "Build the website",
		Long:  `Build the website into dist/.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := buildAll(dist, builders, cfg); err != nil {
				return err
			}

			serve, err := cmd.Flags().GetBool(serveFlag)
			if err != nil {
				return err
			}

			if serve {
				mergedbuilders := map[string]Builder{}
				for pattern, builder := range builders {
					mergedbuilders[pattern] = builder
				}
				for pattern, builder := range globalBuilders {
					mergedbuilders[pattern] = builder
				}

				stop := make(chan struct{})
				reloader := Reloader{
					Builders: mergedbuilders,
					Ignore:   ignoreDirectories,
				}

				var mux http.ServeMux
				mux.Handle("/", http.FileServer(http.Dir(dist)))
				mux.Handle("/ws", reloader.Handler())

				server := http.Server{
					Addr:    fmt.Sprintf(":%d", port),
					Handler: &mux,
				}
				server.RegisterOnShutdown(reloader.ShutdownFunc())

				go listenForCtrlC(stop, &server, &reloader)
				go startFileServer(&server)

				if err := reloader.Watch(); err != nil {
					return err
				}

				log.Printf("Serving %s on http://localhost:%d\n", dist, port)
				<-stop
				return nil
			}

			return nil
		},
	}
	builders = map[string]Builder{
		"src/img/*/*/*.[jJ][pP][gG]": winter.BuildPhoto,
		"src/favicon/*":              buildStaticFile("src/favicon"),
		"public/*":                   buildStaticFile("public"),
		"public/*/*":                 buildStaticFile("public"),
		"public/*/*/*":               buildStaticFile("public"),
	}
	// globalBuilders must be separate from builders because buildTheWorld depends
	// on builders being populated.
	globalBuilders = map[string]Builder{
		"src/templates/*": func(_, _ string, cfg winter.Config) error { return buildAll(dist, builders, cfg) },
		"*.css":           func(_, _ string, cfg winter.Config) error { return buildAll(dist, builders, cfg) },
	}
	serve bool
)

func init() {
	f := buildCmd.PersistentFlags()

	author := *f.StringP("author", "a", "", "author (e.g. Benjamin Carlsson <ben@twos.dev>)")
	if author == "" {
		author = "Anonymous <unspecified>"
	}
	authorParts := authorPattern.FindStringSubmatch(author)
	if len(authorParts) != 3 {
		fmt.Println("invalid --author format: must be \"NAME <EMAIL>\"")
		os.Exit(1)
	}

	cfg = winter.Config{
		AuthorEmail: authorParts[2],
		AuthorName:  authorParts[1],
		Desc:        *f.StringP("desc", "d", "", "site description (e.g. misc thoughts)"),
		Domain: url.URL{
			Scheme: "https",
			Host:   *f.StringP("domain", "m", "", "site root domain (e.g. twos.dev)"),
		},
		Name: *f.StringP("name", "n", "", "site name (e.g. twos.dev)"),
		Since: *f.IntP(
			"since",
			"y",
			0,
			"site year of creation (e.g. 2021)",
		),
	}

	serve = *f.BoolP(serveFlag, "s", false, "start a webserver and rebuild on file changes")

	rootCmd.AddCommand(buildCmd)
}

func listenForCtrlC(stop chan struct{}, srvr *http.Server, reloader *Reloader) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	log.Println("Ctrl+C detected, stopping...")
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
		log.Fatal(err)
	}
}
