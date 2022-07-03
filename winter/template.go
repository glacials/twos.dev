package winter

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"path/filepath"
	"strings"
	"time"

	"twos.dev/winter/graphic"
)

type templateVars struct {
	*document
	*substructure

	Now time.Time
}

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

type tocPartialVars struct {
	Items []tocVars
}

type tocVars struct {
	Anchor string
	Items  []tocVars
	Title  string
}

func loadTemplates(t *template.Template) error {
	if t == nil {
		return fmt.Errorf("nil template")
	}
	partials, err := filepath.Glob("src/templates/*.html.tmpl")
	if err != nil {
		return fmt.Errorf("can't glob for partials: %w", err)
	}

	for _, partial := range partials {
		name := filepath.Base(partial)
		name = strings.TrimPrefix(name, "_")
		name, _, _ = strings.Cut(name, ".") // Trim extensions, even e.g. .html.tmpl

		p, err := ioutil.ReadFile(partial)
		if err != nil {
			return fmt.Errorf(
				"can't read partial `%s`: %w",
				partial,
				err,
			)
		}

		if _, err := t.New(name).Parse(string(p)); err != nil {
			return fmt.Errorf(
				"can't parse partial `%s`: %w",
				partial,
				err,
			)
		}
	}

	return nil
}

type imgsfunc func(graphic.Caption, ...string) (template.HTML, error)

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

type videosfunc func(graphic.Caption, ...string) (template.HTML, error)

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

type postsfunc func() []*document

func posts(s substructure) postsfunc {
	return func() []*document {
		return s.posts()
	}
}
