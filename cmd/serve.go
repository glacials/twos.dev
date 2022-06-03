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
	"time"

	"github.com/spf13/cobra"
)

var noBuild *bool

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start a local twos.dev server",
	Long:  `Start a local twos.dev server by serving files from dist/.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		port := 8100
		http.Handle("/", http.FileServer(http.Dir("dist/")))

		stop := make(chan struct{})

		if !*noBuild {
			go func() {
				for {
					if err := buildCmd.RunE(cmd, args); err != nil {
						log.Fatalf(fmt.Errorf("cannot build: %w", err).Error())
					}
					select {
					case <-time.After(1 * time.Second):
					case <-stop:
						return
					}
				}
			}()
		}

		log.Printf("Serving dist/ on http://localhost:%d\n", port)
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))

		stop <- struct{}{}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
	noBuild = serveCmd.Flags().Bool("no-build", false, "don't continually rebuild while serving")
}
