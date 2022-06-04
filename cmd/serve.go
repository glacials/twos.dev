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
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/glacials/twos.dev/cmd/builders"
	"github.com/spf13/cobra"
)

const dst = "dist"

var (
	noBuild *bool
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

		templater, err := NewTemplateBuilder()
		if err != nil {
			return fmt.Errorf("can't build template builder: %w", err)
		}

		distributer := builders.NewDistributer(map[string]func(src, dst string) error{
			"./src/img/*/*.[jJ][pP][gG]": imageBuilder,
			"./src/*.html":               templater.htmlBuilder,
			"./src/*.md":                 templater.markdownBuilder,
			"./public":                   staticFileBuilder,
		},
		)

		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			return fmt.Errorf("cannot initialize fsnotify watcher: %s", err)
		}
		stop := make(chan struct{})

		if !*noBuild {
			if err := buildCmd.RunE(cmd, args); err != nil {
				log.Fatalf("can't build: %s", err.Error())
			}

			go func() {
				for {
					select {
					case event, ok := <-watcher.Events:
						if !ok {
							return
						}
						log.Println("event:", event)
						if event.Op&fsnotify.Write == fsnotify.Write {
							for path, builder := range distributer.Assignments {
								fmt.Printf("matching %s to %s\n", path, event.Name)

								if ok, err := filepath.Match(path, event.Name); err != nil {
									log.Fatalf("can't match `%s`: %s", path, err)
								} else if ok {
									fmt.Printf("building\n")
									if err := builder(event.Name, dst); err != nil {
										log.Fatalf("can't build `%s`: %s", path, err)
									}
								}
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
			for pattern := range distributer.Assignments {
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
		}

		log.Printf("Serving %s on http://localhost:%d\n", dst, port)
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))

		stop <- struct{}{}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
	noBuild = serveCmd.Flags().Bool("no-build", false, "don't continually rebuild while serving")
}
