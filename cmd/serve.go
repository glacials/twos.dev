package cmd

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	live "github.com/glacials/twos.dev/cmd/livereload"
	"github.com/spf13/cobra"
)

const dst = "dist"

var (
	noBuild *bool
	debug   *bool

	builders = map[string]live.Builder{
		"src/img/*/*/*.[jJ][pP][gG]": photoBuilder,
		"src/cold/*.html.tmpl":       buildDocument,
		"src/cold/*.html":            buildDocument,
		"src/cold/*.md":              buildDocument,
		"src/warm/*.md":              buildDocument,
		"public/*":                   staticFileBuilder("public"),
		"public/*/*":                 staticFileBuilder("public"),
		"public/*/*/*":               staticFileBuilder("public"),
	}

	// globalBuilders must be separate from builders because buildTheWorld depends
	// on builders being populated.
	globalBuilders = map[string]live.Builder{
		"src/templates/*": func(_, _ string) error { return buildTheWorld() },
		"*.css":           func(_, _ string) error { return buildTheWorld() },
	}
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start a local twos.dev server",
	Long: fmt.Sprintf(
		`Start a local twos.dev server by serving files from %s.`,
		dst,
	),
	RunE: func(cmd *cobra.Command, args []string) error {
		port := 8100
		mergedbuilders := map[string]live.Builder{}
		for pattern, builder := range builders {
			mergedbuilders[pattern] = builder
		}
		for pattern, builder := range globalBuilders {
			mergedbuilders[pattern] = builder
		}

		stop := make(chan struct{})
		reloader := live.Reloader{
			Builders: mergedbuilders,
			Ignore:   ignoreDirectories,
		}

		var mux http.ServeMux
		mux.Handle("/", http.FileServer(http.Dir(dst)))
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
			if err := buildTheWorld(); err != nil {
				log.Fatalf("can't build: %s", err.Error())
			}

			err := reloader.Watch()
			if err != nil {
				return err
			}
		}

		log.Printf("Serving %s on http://localhost:%d\n", dst, port)
		<-stop

		return nil
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	noBuild = serveCmd.Flags().
		Bool("no-build", false, "don't continually rebuild while serving")

	debug = serveCmd.Flags().
		Bool("debug", false, "output to dist/debug/ results of each transformation")
}

func listenForCtrlC(
	stop chan struct{},
	server *http.Server,
	reloader *live.Reloader,
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
