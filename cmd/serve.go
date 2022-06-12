/*
Copyright © 2022 Benjamin Carlsson

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
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

	builders = map[string]func(src, dst string) error{
		"src/img/*/*/*.[jJ][pP][gG]": photoBuilder,
		"src/cold/*.html":            htmlBuilder,
		"src/cold/*.md":              markdownBuilder,
		"src/warm/*.md":              markdownBuilder,
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

		c := make(chan os.Signal)
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
										log.Fatalf("can't build `%s`: %s", pattern, err)
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

			http.Handle("/ws", websocket.Handler(func(conn *websocket.Conn) {
				// TODO: This works sometimes because multiple ws connecions are open
				// (not timed out yet, or just multiple pages) and so these consumers
				// fight to receive the message.
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
}
