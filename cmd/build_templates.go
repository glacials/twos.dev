package cmd

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/glacials/twos.dev/cmd/frontmatter"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
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

	builder := templateBuilder{
		essayTemplate: essay,
	}

	buildHTMLFile := func(src, dst string) error {
		f, err := os.Open(src)
		if err != nil {
			// TODO: Clean this up; Prettier autoformatting seems to remove the file
			// and quickly place it back, so ignore th error.
			log.Printf(
				fmt.Errorf(
					"can't open HTML file at `%s` for building: %w",
					src,
					err,
				).Error(),
			)
			return nil
		}

		matter, body, err := frontmatter.Parse(f)
		if err != nil {
			return fmt.Errorf("can't get frontmatter from Markdown file: %w", err)
		}

		if matter.Filename == "" {
			matter.Filename = filepath.Base(src)
		}

		if err := builder.buildHTMLStream(bytes.NewBuffer(body), src, dst, matter); err != nil {
			return fmt.Errorf("can't build HTML stream: %w", err)
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
				"can't open Markdown file at `%s` for building: %w",
				src,
				err,
			)
		}

		matter, body, err := frontmatter.Parse(f)
		if err != nil {
			return fmt.Errorf("can't get frontmatter from Markdown file: %w", err)
		}

		if err := os.MkdirAll(dst, 0755); err != nil {
			return fmt.Errorf("can't make destination directory `%s`: %w", dst, err)
		}

		if err := builder.buildHTMLStream(bytes.NewBuffer(body), src, dst, matter); err != nil {
			return fmt.Errorf("can't build HTML stream: %w", err)
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

func (builder templateBuilder) buildHTMLStream(r io.Reader, src string, dst string, matter frontmatter.Matter) error {
	if err := os.MkdirAll(dst, 0755); err != nil {
		return fmt.Errorf("can't make destination directory `%s`: %w", dst, err)
	}

	if matter.Filename == "" {
		return fmt.Errorf("file frontmatter has no filename attribute")
	}

	destinationFilePath := filepath.Join(dst, matter.Filename)
	htmlFile, err := os.Create(destinationFilePath)
	if err != nil {
		return fmt.Errorf(
			"can't render HTML to `%s` from template for `%s`: %w",
			destinationFilePath,
			r,
			err,
		)
	}

	var createdAt, updatedAt string
	if !matter.CreatedAt.IsZero() {
		createdAt = matter.CreatedAt.Format("2006 January")
	}
	if !matter.UpdatedAt.IsZero() {
		updatedAt = matter.CreatedAt.Format("2006 January")
	}

	body, err := io.ReadAll(r)
	if err != nil {
		return fmt.Errorf("can't read body from stream: %w", err)
	}

	title, err := titleFromHTML(bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("can't get title from HTML: %w", err)
	}

	v := htmlFileVars{
		Body:  template.HTML(body),
		Title: title,
		SourceURL: fmt.Sprintf(
			"https://github.com/glacials/twos.dev/blob/main/%s",
			src,
		),

		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}

	if err := builder.essayTemplate.Execute(htmlFile, v); err != nil {
		return fmt.Errorf("can't execute essay template: %w", err)
	}

	return nil
}

func titleFromHTML(r io.Reader) (string, error) {
	h, err := html.Parse(r)
	if err != nil {
		return "", fmt.Errorf("can't parse HTML: %w", err)
	}

	if h1, ok := firstH1(h); ok {
		for child := h1.FirstChild; child != nil; child = child.NextSibling {
			if child.Type == html.TextNode {
				return child.Data, nil
			}
		}
	}

	return "twos.dev", nil
}

func firstH1(n *html.Node) (*html.Node, bool) {
	if n.DataAtom == atom.H1 {
		return n, true
	}

	for child := n.FirstChild; child != nil; child = child.NextSibling {
		if h1, ok := firstH1(child); ok {
			return h1, true
		}
	}

	return nil, false
}
