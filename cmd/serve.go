package cmd

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
	"golang.org/x/net/websocket"
)

const dst = "dist"

var (
	noBuild *bool
	debug   *bool

	builders = map[string]func(src, dst string) error{
		"src/img/*/*/*.[jJ][pP][gG]": photoBuilder,
		"src/cold/*.html.tmpl":       buildDocument,
		"src/cold/*.html":            buildDocument,
		"src/cold/*.md":              buildDocument,
		"src/warm/*.md":              buildDocument,
		"public/*":                   staticFileBuilder("public"),
		"public/*/*":                 staticFileBuilder("public"),
		"public/*/*/*":               staticFileBuilder("public"),
	}

	buildTheWorldTriggers = map[string]struct{}{
		"src/templates/*.html": {},
		"*.css":                {},
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
		http.Handle("/", http.FileServer(http.Dir("dist/")))
		var refreshListeners []chan struct{}

		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			return fmt.Errorf("cannot initialize fsnotify watcher: %s", err)
		}
		stop := make(chan struct{})

		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		go func() {
			select {
			case <-c:
				stop <- struct{}{}
			case <-stop:
			}
			os.Exit(0)
		}()

		if !*noBuild {
			go func() {
				if err := buildTheWorld(); err != nil {
					log.Fatalf("can't build: %s", err.Error())
				}
				for _, ch := range refreshListeners {
					ch <- struct{}{}
				}
			}()

			go func() {
				for {
					select {
					case event, ok := <-watcher.Events:
						if !ok {
							return
						}
						if event.Op&(fsnotify.Write|fsnotify.Rename|fsnotify.Create) > 0 {
							log.Println("changed:", event.Name)
							for pattern, builder := range builders {
								if ok, err := filepath.Match(pattern, event.Name); err != nil {
									log.Fatalf("can't match `%s`: %s", pattern, err)
								} else if ok {
									if err := builder(event.Name, dst); err != nil {
										log.Printf("can't build `%s`: %s", pattern, err)
									}
								}
							}
							for pattern := range buildTheWorldTriggers {
								if ok, err := filepath.Match(pattern, event.Name); err != nil {
									log.Fatalf("can't match `%s`: %s", pattern, err)
								} else if ok {
									if err := buildTheWorld(); err != nil {
										log.Fatalf("can't build the world: %s", err)
									}
								}
							}
							for _, ch := range refreshListeners {
								ch <- struct{}{}
							}
						}
					case err, ok := <-watcher.Errors:
						if !ok {
							return
						}
						log.Println("error:", err)
					case <-stop:
						return
					}
				}
			}()

			watched := map[string]struct{}{}
			for pattern := range builders {
				paths, err := filepath.Glob(pattern)
				if err != nil {
					return fmt.Errorf("can't glob `%s`: %w", pattern, err)
				}

				for _, path := range paths {
					for p := path; p != "."; p = filepath.Dir(p) {
						if _, ok := watched[p]; !ok {
							if err := watcher.Add(p); err != nil {
								return fmt.Errorf("cannot watch file %s: %w", path, err)
							}
						}
						watched[p] = struct{}{}
					}
				}
			}
			for pattern := range buildTheWorldTriggers {
				paths, err := filepath.Glob(pattern)
				if err != nil {
					return fmt.Errorf("can't glob `%s`: %w", pattern, err)
				}

				for _, path := range paths {
					for p := path; p != "."; p = filepath.Dir(p) {
						if _, ok := watched[p]; !ok {
							if err := watcher.Add(p); err != nil {
								return fmt.Errorf("cannot watch file %s: %w", path, err)
							}
						}
						watched[p] = struct{}{}
					}
				}
			}

			// Auto refresh the page when a change is made
			http.Handle("/ws", websocket.Handler(func(conn *websocket.Conn) {
				c := make(chan struct{})
				refreshListeners = append(refreshListeners, c)
				for {
					<-c
					conn.Write([]byte("refresh"))
				}
			}))
		}

		log.Printf("Serving %s on http://localhost:%d\n", dst, port)
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))

		stop <- struct{}{}

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
