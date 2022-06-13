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
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/glacials/twos.dev/cmd/frontmatter"
	"github.com/glacials/twos.dev/cmd/img"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/parser"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type essayPageVars struct {
	pageVars

	Body      template.HTML
	Title     string
	Shortname string

	CreatedAt string
	UpdatedAt string
}

// videoPartialVars are the template variables given to
// src/templates/_video.html to render a video inline. At least one of its
// {Light,Dark}{MOV,MP4} fields must be set to a video path.
type videoPartialVars struct {
	LightMOV string
	LightMP4 string
	DarkMOV  string
	DarkMP4  string
}

type imageVars struct {
	Alt   string
	Light string
	Dark  string
}

// imgPartialVars are the template variables given to src/templates/_img.html to
// render an image inline. At least one of its {Light,Dark} fields must be set
// to an image path.
type imgPartialVars struct {
	imageVars
	Caption string
}

// imgsPartialVars are the template variables given to src/templates/_imgs.html
// to render multiple images inline. At least one of its {Light,Dark} fields
// must be set to an image path.
type imgsPartialVars struct {
	Images  []imageVars
	Caption string
}

type pageVars struct {
	// SourceURL is the GitHub URL to the source code for the page being rendered.
	SourceURL string
}

func htmlBuilder(src, dst string) error {
	f, err := os.Open(src)
	if err != nil {
		// TODO: Clean this up. Prettier autoformatting seems to remove the file
		// and quickly place it back, so for now we ignore the error.
		log.Printf(
			fmt.Errorf(
				"can't open HTML file at `%s` for building: %w",
				src,
				err,
			).Error(),
		)
		return nil
	}
	defer f.Close()

	matter, body, err := frontmatter.Parse(f)
	if err != nil {
		return fmt.Errorf(
			"can't get frontmatter from Markdown file: %w",
			err,
		)
	}

	if matter.Filename == "" {
		matter.Filename = filepath.Base(src)
	}

	if err := buildHTMLStream(
		bytes.NewBuffer(body),
		src,
		filepath.Join(dst, matter.Filename),
		matter,
	); err != nil {
		return fmt.Errorf("can't build HTML stream: %w", err)
	}

	return nil
}

func markdownBuilder(src, dst string) error {
	f, err := os.Open(src)
	if err != nil {
		return fmt.Errorf(
			"can't open Markdown file at `%s` for building: %w",
			src,
			err,
		)
	}
	defer f.Close()

	matter, body, err := frontmatter.Parse(f)
	if err != nil {
		return fmt.Errorf(
			"can't get frontmatter from Markdown file: %w",
			err,
		)
	}

	// Markdown parser cannot be reused :(
	renderedHTML := markdown.ToHTML(body, parser.NewWithExtensions(
		parser.Tables|
			parser.FencedCode|
			parser.Autolink|
			parser.Strikethrough|
			parser.Footnotes|
			parser.HeadingIDs|
			parser.Attributes|
			parser.SuperSubscript,
	), nil)

	if err := buildHTMLStream(
		bytes.NewBuffer(renderedHTML),
		src,
		filepath.Join(dst, matter.Filename),
		matter,
	); err != nil {
		return fmt.Errorf("can't build HTML stream: %w", err)
	}

	return nil
}

func buildHTMLStream(
	r io.Reader,
	src string,
	dst string,
	matter frontmatter.Matter,
) error {
	if matter.Filename == "" {
		return fmt.Errorf("file frontmatter has no filename attribute")
	}

	htmlFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf(
			"can't render HTML to `%s` from template: %w",
			dst,
			err,
		)
	}
	defer htmlFile.Close()

	var createdAt, updatedAt string
	if !matter.CreatedAt.IsZero() {
		createdAt = matter.CreatedAt.Format("2006 January")
	}
	if !matter.UpdatedAt.IsZero() {
		updatedAt = matter.UpdatedAt.Format("2006 January")
	}

	body, err := io.ReadAll(r)
	if err != nil {
		return fmt.Errorf("can't read body from stream: %w", err)
	}

	body = bytes.ReplaceAll(body, []byte("“"), []byte("\""))
	body = bytes.ReplaceAll(body, []byte("”"), []byte("\""))

	title, err := titleFromHTML(bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("can't get title from HTML: %w", err)
	}

	v := essayPageVars{
		Body:      template.HTML(body),
		Title:     title,
		Shortname: strings.TrimSuffix(matter.Filename, ".html"),

		CreatedAt: createdAt,
		UpdatedAt: updatedAt,

		pageVars: pageVars{
			SourceURL: fmt.Sprintf(
				"https://github.com/glacials/twos.dev/blob/main/%s",
				src,
			),
		},
	}

	t, err := template.ParseFiles("src/templates/essay.html")
	if err != nil {
		return fmt.Errorf("can't parse essay template: %w", err)
	}

	t.Funcs(template.FuncMap{
		"videos": func(video string) (videoPartialVars, error) {
			v := videoPartialVars{
				DarkMOV: fmt.Sprintf(
					"img/%s-%s-dark.mov",
					v.Shortname,
					video,
				),
				LightMOV: fmt.Sprintf(
					"img/%s-%s-light.mov",
					v.Shortname,
					video,
				),
				DarkMP4: fmt.Sprintf(
					"img/%s-%s-light.mp4",
					v.Shortname,
					video,
				),
				LightMP4: fmt.Sprintf(
					"img/%s-%s-light.mp4",
					v.Shortname,
					video,
				),
			}

			if _, err := os.Stat(filepath.Join("dist", v.DarkMOV)); err != nil {
				if errors.Is(err, os.ErrNotExist) {
					v.DarkMOV = ""
				} else {
					return videoPartialVars{}, fmt.Errorf(
						"couldn't stat video `%s`: %w",
						v.DarkMOV,
						err,
					)
				}
			}

			if _, err := os.Stat(filepath.Join("dist", v.LightMOV)); err != nil {
				if errors.Is(err, os.ErrNotExist) {
					v.LightMOV = ""
				} else {
					return videoPartialVars{}, fmt.Errorf(
						"couldn't stat video `%s`: %w",
						v.LightMOV,
						err,
					)
				}
			}

			if _, err := os.Stat(filepath.Join("dist", v.DarkMP4)); err != nil {
				if errors.Is(err, os.ErrNotExist) {
					v.DarkMP4 = ""
				} else {
					return videoPartialVars{}, fmt.Errorf(
						"couldn't stat video `%s`: %w",
						v.DarkMP4,
						err,
					)
				}
			}

			if _, err := os.Stat(filepath.Join("dist", v.LightMP4)); err != nil {
				if errors.Is(err, os.ErrNotExist) {
					v.LightMP4 = ""
				} else {
					return videoPartialVars{}, fmt.Errorf(
						"couldn't stat video `%s`: %w",
						v.LightMP4,
						err,
					)
				}
			}

			return v, nil
		},
		"imgs": func(imageA, altA, imageB, altB, caption string) (template.HTML, error) {
			lightA, darkA, err := img.LightDark(v.Shortname, imageA)
			if err != nil {
				return "", fmt.Errorf("can't process image: %w", err)
			}

			lightB, darkB, err := img.LightDark(v.Shortname, imageB)
			if err != nil {
				return "", fmt.Errorf("can't process image: %w", err)
			}

			v := imgsPartialVars{
				Caption: caption,
				Images: []imageVars{
					{
						Alt:   altA,
						Dark:  darkA,
						Light: lightA,
					},
					{
						Alt:   altB,
						Dark:  darkB,
						Light: lightB,
					},
				},
			}

			buf := bytes.NewBuffer([]byte{})
			if err := t.Lookup("imgs").Execute(buf, v); err != nil {
				return "", fmt.Errorf(
					"can't execute imgs template for `%s`/`%s`: %w",
					imageA,
					imageB,
					err,
				)
			}

			return template.HTML(buf.String()), nil
		},
		"img": func(image, caption, alt string) (template.HTML, error) {
			light, dark, err := img.LightDark(v.Shortname, image)
			if err != nil {
				return "", fmt.Errorf("can't process image: %w", err)
			}

			v := imgPartialVars{
				imageVars: imageVars{
					Alt:   alt,
					Dark:  dark,
					Light: light,
				},
				Caption: caption,
			}

			buf := bytes.NewBuffer([]byte{})
			if err := t.Lookup("imgs").Execute(buf, v); err != nil {
				return "", fmt.Errorf(
					"can't execute img template for `%s`: %w",
					image,
					err,
				)
			}

			return template.HTML(buf.String()), nil
		},
	})

	partials, err := filepath.Glob("src/templates/_*.html")
	if err != nil {
		return fmt.Errorf("can't glob for partial templates: %w", err)
	}

	for _, partial := range partials {
		p := t.New(
			strings.TrimPrefix(
				strings.TrimSuffix(filepath.Base(partial), ".html"),
				"_",
			),
		)

		s, err := ioutil.ReadFile(partial)
		if err != nil {
			return fmt.Errorf("can't read partial `%s`: %w", partial, err)
		}

		if _, err := p.Parse(string(s)); err != nil {
			return fmt.Errorf("can't parse partial `%s`: %w", partial, err)
		}
	}

	bodyTemplate := t.New(src)
	_, err = bodyTemplate.Parse(strings.Join(
		[]string{"{{define \"body\"}}", string(body), "{{end}}"},
		"\n",
	),
	)
	if err != nil {
		return fmt.Errorf("can't parse template `%s`: %w", src, err)
	}

	if err := t.Execute(htmlFile, v); err != nil {
		return fmt.Errorf("can't execute essay template `%s`: %w", src, err)
	}

	return nil
}

func buildTemplate(src, dst string) error {
	if err := buildTheWorld(); err != nil {
		return fmt.Errorf("can't \"build\" template `%s`: %w", src, err)
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
