package winter

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"path/filepath"
	"strings"
	"text/template/parse"
	"time"

	"twos.dev/winter/graphic"
)

const (
	txtname = "text_document"
	galname = "imgcontainer"
)

type archivesVars []archiveVars

func (a archivesVars) Less(i, j int) bool {
	return a[i].Year > a[j].Year
}

func (a archivesVars) Len() int {
	return len(a)
}

func (a archivesVars) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

type archiveVars struct {
	Year      int
	Documents documents
}

type imgsfunc func(graphic.Caption, ...string) (template.HTML, error)

// imgsPartialVars are the template variables given to
// src/templates/_imgs.html.tmpl to render multiple images inline. At least one
// of its {Light,Dark} fields must be set to an image path.
type imgsPartialVars struct {
	Images  []imageVars
	Caption graphic.Caption
}

// imageVars are the template variables used to render a single image. At least
// one field must be set. If both are set, the rendered image will depend on
// the user agent.
type imageVars struct {
	Alt   graphic.Alt
	Light graphic.SRC
	Dark  graphic.SRC
}

type textDocumentVars struct {
	*document
	*Substructure
	Now time.Time
}

type tocPartialVars struct {
	Items []tocVars
}

type tocVars struct {
	Anchor string
	Items  []tocVars
	Title  string
}

type videosfunc func(graphic.Caption, ...string) (template.HTML, error)

// videosPartialVars are the template variables given to
// src/templates/_videos.html.tmpl to render any number of videos inline.
type videosPartialVars struct {
	Videos  []videoVars
	Caption graphic.Caption
}

// videoVars are the template variables used to render a single video. At least
// one field must be set. If multiple are set, the rendered video will depend on
// the user agent.
type videoVars struct {
	LightMOV graphic.SRC
	DarkMOV  graphic.SRC
	LightMP4 graphic.SRC
	DarkMP4  graphic.SRC
}

// tmplPathToName converts a template path to a template name.
func tmplPathToName(src string) string {
	name := filepath.Base(src)
	name = strings.TrimPrefix(name, "_")
	name, _, _ = strings.Cut(name, ".") // Trim extensions, even e.g. .html.tmpl
	return name
}

// tmplNameToPath converts a template name to a path. If no such template
// exists, returns an error.
func tmplNameToPath(name string) (string, error) {
	paths, err := filepath.Glob(
		filepath.Join("src", "templates", fmt.Sprintf("*%s.*", name)),
	)
	if err != nil {
		return "", fmt.Errorf("can't glob for templates: %w", err)
	}

	if len(paths) > 1 {
		return "", fmt.Errorf("multiple files match template name `%s`", name)
	}
	if len(paths) == 0 {
		return "", fmt.Errorf(
			"no file for template `%s`; expected one at src/templates/[_]%s.html[.tmpl]",
			name,
			name,
		)
	}

	return paths[0], nil
}

func add(a, b int) int {
	return a + b
}
func sub(a, b int) int {
	return a - b
}

// imgs returns a function that can be inserted into a template's FuncMap for
// calling by the template. The returned function takes a caption followed by
// any number of pairs of image shortnames and alt texts; the total number of
// arguments must therefore be odd. The images are rendered next to each other
// with the caption below them. If light and dark mode images are both present,
// the correct ones will be rendered based on the user's mode.
//
// The images must be named in the format SHORTNAME-IMAGE-STYLE.EXTENSION, where
// SHORTNAME is the page shortname (e.g. "apple" for apple.html), IMAGE is any
// arbitrary string that must be passed to the returned function as its imgsrc,
// STYLE is one of "light" or "dark", and EXTENSION is one of "png" or "jpg".
//
// The given shortname must be the page shortname the images will appear on, or
// the rendered images won't point to the right URLs.
func imgs(shortname string) (imgsfunc, error) {
	partial, err := ioutil.ReadFile("src/templates/_imgs.html.tmpl")
	if err != nil {
		return nil, err
	}

	t := template.New("imgs")
	if _, err := t.Parse(string(partial)); err != nil {
		return nil, err
	}

	return func(c graphic.Caption, srcsAndAlts ...string) (template.HTML, error) {
		v := imgsPartialVars{Caption: c}

		for i := 0; i < len(srcsAndAlts); i += 2 {
			light, dark, err := graphic.LightDark(
				shortname,
				graphic.Shortname(srcsAndAlts[i]),
				graphic.ImageExts,
			)
			if err != nil {
				return "", err
			}

			v.Images = append(v.Images, imageVars{
				Alt:   graphic.Alt(srcsAndAlts[i+1]),
				Light: light,
				Dark:  dark,
			})
		}

		var buf bytes.Buffer
		if err := t.Execute(&buf, v); err != nil {
			return "", fmt.Errorf("can't execute imgs template: %w", err)
		}

		return template.HTML(buf.String()), nil
	}, nil
}

// videos returns a function that can be inserted into a template's FuncMap for
// calling by the template. The returned function takes a caption followed by
// any number of video shortnames. The videos are rendered next to each other
// with the caption below them. If light and dark mode videos are both present,
// the correct ones will be rendered based on the user's mode.
//
// The videos must be named in the format SHORTNAME-VIDEO-STYLE.EXTENSION, where
// SHORTNAME is the page shortname (e.g. "apple" for apple.html), VIDEO is any
// arbitrary string that must be passed to the returned function as its
// videosrc., STYLE is one of "light" or "dark". EXTENSION is "mp4" and/or
// "mov"; if both are present, the user agent will choose which to load.
//
// The given shortname must be the page shortname the videos will appear on, or
// the rendered videos won't point to the right URLs.
func videos(
	pageShortname string,
) (func(graphic.Caption, ...graphic.Shortname) (template.HTML, error), error) {
	partial, err := ioutil.ReadFile("src/templates/_videos.html.tmpl")
	if err != nil {
		return nil, err
	}

	t := template.New("videos")
	if _, err := t.Parse(string(partial)); err != nil {
		return nil, err
	}

	return func(c graphic.Caption, videoShortnames ...graphic.Shortname) (template.HTML, error) {
		v := videosPartialVars{Caption: c}

		for _, videoShortname := range videoShortnames {
			lightMP4, darkMP4, err := graphic.LightDark(
				pageShortname,
				videoShortname,
				map[string]struct{}{"mp4": {}},
			)
			if err != nil {
				return "", fmt.Errorf("can't process video: %w", err)
			}

			lightMOV, darkMOV, err := graphic.LightDark(
				pageShortname,
				videoShortname,
				map[string]struct{}{"mp4": {}},
			)
			if err != nil {
				return "", fmt.Errorf("can't process video: %w", err)
			}

			v.Videos = append(v.Videos, videoVars{
				DarkMOV:  darkMOV,
				DarkMP4:  darkMP4,
				LightMOV: lightMOV,
				LightMP4: lightMP4,
			},
			)
		}

		var buf bytes.Buffer
		if err := t.Execute(&buf, v); err != nil {
			return "", fmt.Errorf(
				"can't execute video template for `%s`/`%s`: %w",
				pageShortname,
				videoShortnames,
				err,
			)
		}

		return template.HTML(buf.String()), nil
	}, nil
}

// loadDependencies loads and parses the given template from disk if needed,
// then searches its parse tree for templates it references, and loads and
// parses those templates.
func loadDependencies(t *template.Template) error {
	if t.Tree == nil {
		name, err := tmplNameToPath(t.Name())
		if err != nil {
			return err
		}
		b, err := ioutil.ReadFile(name)
		if err != nil {
			return err
		}
		if _, err := t.Parse(string(b)); err != nil {
			return err
		}
	}
	for _, node := range t.Tree.Root.Nodes {
		if node.Type() == parse.NodeTemplate {
			name := node.(*parse.TemplateNode).Name
			if t.Lookup(name) == nil {
				path, err := tmplNameToPath(name)
				if err != nil {
					return err
				}
				b, err := ioutil.ReadFile(path)
				if err != nil {
					return err
				}
				if _, err := t.New(name).Parse(string(b)); err != nil {
					return err
				}
			}
		}
	}
	return nil
}
