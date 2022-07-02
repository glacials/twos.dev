package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	"twos.dev/winter"
)

var (
	noBuild *bool

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
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Build and serve the static website",
	Long: fmt.Sprintf(
		`Start a local Winter server by continually building and serving files.`,
	),
	RunE: func(cmd *cobra.Command, args []string) error {
		port := 8100
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

		go func() {
			if err := server.ListenAndServe(); err != nil {
				if errors.Is(err, http.ErrServerClosed) {
					// Expected when we call server.Shutdown()
					return
				}
				log.Fatal(err)
			}
		}()

		if !*noBuild {
			if err := buildAll(dist, builders, cfg); err != nil {
				log.Fatalf("can't build: %s", err.Error())
			}

			err := reloader.Watch()
			if err != nil {
				return err
			}
		}

		log.Printf("Serving %s on http://localhost:%d\n", dist, port)
		<-stop

		return nil
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	noBuild = serveCmd.Flags().
		Bool("no-build", false, "don't continually rebuild while serving")
}

func listenForCtrlC(
	stop chan struct{},
	server *http.Server,
	reloader *Reloader,
) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	log.Println("Ctrl+C detected, stopping...")
	if err := server.Shutdown(context.TODO()); err != nil {
		log.Fatal(err)
	}
	stop <- struct{}{}
}
