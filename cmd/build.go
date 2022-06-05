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
	"html/template"

	"github.com/spf13/cobra"
)

const (
	imageContainerTemplatePath = "src/templates/_imgcontainer.html"
	staticAssetsDir            = "public"
	sourceDir                  = "src"
)

var (
	ignoreSrcDirs = map[string]struct{}{
		"src/asciiart": {},
		"src/js":       {},
	}
)

type essayFrontmatter struct {
	Filename string `yaml:"filename"`
	// Date is an alias for CreatedAt.
	Date      string `yaml:"date"`
	CreatedAt string `yaml:"created"`
	UpdatedAt string `yaml:"updated"`
}

type htmlFileVars struct {
	Body      template.HTML
	Title     string
	SourceURL string

	CreatedAt string
	UpdatedAt string
}

type imageContainerVars struct {
	PrevLink  string
	CurImage  string
	NextLink  string
	SourceURL string
}

// buildCmd represents the build command
var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build twos.dev",
	Long:  `Build twos.dev into dist/.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := buildTheWorld(); err != nil {
			return fmt.Errorf("can't build the world: %w", err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)
}
