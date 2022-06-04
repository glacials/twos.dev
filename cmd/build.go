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
	"os"
	"path/filepath"
	"strings"

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
		templater, err := NewTemplateBuilder()
		if err != nil {
			return fmt.Errorf("can't make template builder: %w", err)
		}

		if err := filepath.WalkDir("src", func(src string, d os.DirEntry, err error) error {
			if err != nil {
				return fmt.Errorf("can't walk src for build `%s`: %w", src, err)
			}
			if d.IsDir() {
				return nil
			}
			if d.Name() == ".DS_Store" {
				return nil
			}
			for path := range ignoreSrcDirs {
				if strings.HasPrefix(src, path) {
					return filepath.SkipDir
				}
			}

			if ok, err := filepath.Match("*.md", d.Name()); err != nil {
				return fmt.Errorf("can't match `%s` to *.md: %w", src, err)
			} else if ok {
				if err := templater.markdownBuilder(src, dst); err != nil {
					return fmt.Errorf("can't build Markdown in `%s`: %w", src, err)
				}

				return nil
			}

			if ok, err := filepath.Match("*.html", d.Name()); err != nil {
				return fmt.Errorf("can't match `%s` to *.html: %w", src, err)
			} else if ok {
				if err := templater.htmlBuilder(src, dst); err != nil {
					return fmt.Errorf("can't build HTML in `%s`: %w", src, err)
				}

				return nil
			}

			if ok, err := filepath.Match("*.[jJ][pP][gG]", d.Name()); err != nil {
				return fmt.Errorf("can't match `%s` to *.jpg: %w", src, err)
			} else if ok {
				if err := imageBuilder(src, dst); err != nil {
					return fmt.Errorf("can't build JPEG `%s`: %w", src, err)
				}

				return nil
			}

			return fmt.Errorf("don't know what to do with file `%s`", src)
		}); err != nil {
			return fmt.Errorf("can't watch file: %w", err)
		}

		if err := buildStaticDir(staticAssetsDir, dst); err != nil {
			return fmt.Errorf("can't build static assets: %w", err)
		}

		if err := buildFormatting(dst); err != nil {
			return fmt.Errorf("can't post-process build dir: %w", err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)
}
