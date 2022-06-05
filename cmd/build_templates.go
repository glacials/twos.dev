package cmd

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/glacials/twos.dev/cmd/frontmatter"
)

type templateBuilder struct {
	essayTemplate *template.Template

	htmlBuilder     func(src, dest string) error
	markdownBuilder func(src, dest string) error
	templateBuilder func(src, dest string) error
}

func NewTemplateBuilder() (templateBuilder, error) {
	templateHTML, err := ioutil.ReadFile("src/templates/_essay.html")
	if err != nil {
		return templateBuilder{}, fmt.Errorf("can't read essay template: %w", err)
	}

	essay, err := template.New("essay").Parse(string(templateHTML))
	if err != nil {
		return templateBuilder{}, fmt.Errorf("can't create essay template: %w", err)
	}

	buildHTMLFile := func(src, dst string) error {
		f, err := os.Open(src)
		if err != nil {
			return fmt.Errorf(
				"template builder can't open HTML file at `%s` for building: %w",
				src,
				err,
			)
		}

		matter, body, err := frontmatter.Parse(f)
		if err != nil {
			return fmt.Errorf("can't get frontmatter from HTML file: %w", err)
		}

		if err := os.MkdirAll(dst, 0755); err != nil {
			return fmt.Errorf("can't make destination directory `%s`: %w", dst, err)
		}

		if matter.Filename == "" {
			if !strings.HasSuffix(filepath.Base(src), ".html") {
				return fmt.Errorf("non-HTML file %s must have a filename attribute in frontmatter", src)
			}
			matter.Filename = filepath.Base(src)
		}

		destinationFilePath := filepath.Join(dst, matter.Filename)
		htmlFile, err := os.Create(destinationFilePath)
		if err != nil {
			return fmt.Errorf(
				"can't render HTML to `%s` from template for `%s`: %w",
				destinationFilePath,
				src,
				err,
			)
		}

		v := htmlFileVars{
			Body:  template.HTML(body),
			Title: filepath.Base(src)[0 : len(filepath.Base(src))-len(".md")],
			SourceURL: fmt.Sprintf(
				"https://github.com/glacials/twos.dev/blob/main/%s",
				src,
			),

			CreatedAt: matter.CreatedAt.Format("2006 January"),
			UpdatedAt: matter.UpdatedAt.Format("2006 January"),
		}

		if err := essay.Execute(htmlFile, v); err != nil {
			return fmt.Errorf("can't execute essay template: %w", err)
		}

		return nil
	}

	buildMarkdownFile := func(src, dst string) error {
		renderedHTML := bytes.NewBuffer([]byte{})
		renderCmd := exec.Command("src/js/build.js", "body", src)
		renderCmd.Stdout = renderedHTML
		renderCmd.Stderr = os.Stderr

		if err := renderCmd.Run(); err != nil {
			return fmt.Errorf("can't run `src/js/build.js body '%s'`: %w", src, err)
		}

		f, err := os.Open(src)
		if err != nil {
			return fmt.Errorf(
				"can't open HTML file at `%s` for building: %w",
				src,
				err,
			)
		}

		matter, _, err := frontmatter.Parse(f)
		if err != nil {
			return fmt.Errorf("can't get frontmatter from Markdown file: %w", err)
		}

		if err := os.MkdirAll(dst, 0755); err != nil {
			return fmt.Errorf("can't make destination directory `%s`: %w", dst, err)
		}

		destinationFilePath := filepath.Join(dst, matter.Filename)
		htmlFile, err := os.Create(destinationFilePath)
		if err != nil {
			return fmt.Errorf(
				"can't render HTML to `%s` from template for `%s`: %w",
				destinationFilePath,
				src,
				err,
			)
		}

		v := htmlFileVars{
			Body:  template.HTML(renderedHTML.String()),
			Title: filepath.Base(src)[0 : len(filepath.Base(src))-len(".md")],
			SourceURL: fmt.Sprintf(
				"https://github.com/glacials/twos.dev/blob/main/%s",
				src,
			),

			CreatedAt: matter.CreatedAt.Format("2006 January"),
			UpdatedAt: matter.UpdatedAt.Format("2006 January"),
		}

		if err := essay.Execute(htmlFile, v); err != nil {
			return fmt.Errorf("can't execute essay template: %w", err)
		}
		return nil
	}

	return templateBuilder{
		essayTemplate: essay,

		htmlBuilder:     buildHTMLFile,
		markdownBuilder: buildMarkdownFile,
		templateBuilder: func(src, dst string) error {
			// TODO: Selectively build files that use this template
			if err := buildTheWorld(); err != nil {
				return fmt.Errorf("can't build the world: %w", err)
			}

			return nil
		},
	}, nil
}
