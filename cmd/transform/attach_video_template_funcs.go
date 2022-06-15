package transform

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"

	"github.com/glacials/twos.dev/cmd/document"
	"github.com/glacials/twos.dev/cmd/graphic"
)

const (
	videoPartial = "video"
)

// videoPartialVars are the template variables given to
// src/templates/_video.html.tmpl to render a video inline. At least one of its
// {Light,Dark}{MOV,MP4} fields must be set to a video path.
type videoPartialVars struct {
	LightMOV graphic.SRC
	LightMP4 graphic.SRC
	DarkMOV  graphic.SRC
	DarkMP4  graphic.SRC
}

// AttachVideoTemplateFuncs makes the img and imgs functions available to the
// document's templates.
func AttachVideoTemplateFuncs(d document.Document) (document.Document, error) {
	video, err := video(d.Shortname)
	if err != nil {
		return document.Document{}, err
	}

	_ = d.Template.Funcs(
		template.FuncMap{"video": video},
	)

	return d, nil
}

func video(
	pageShortname string,
) (func(graphic.Shortname, graphic.Caption) (template.HTML, error), error) {
	partial, err := ioutil.ReadFile(partialPath(videoPartial))
	if err != nil {
		return nil, err
	}

	t := template.New(videoPartial)
	if _, err := t.Parse(string(partial)); err != nil {
		return nil, err
	}

	return func(videoShortname graphic.Shortname, caption graphic.Caption) (template.HTML, error) {
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

		v := videoPartialVars{
			DarkMOV:  darkMOV,
			DarkMP4:  darkMP4,
			LightMOV: lightMOV,
			LightMP4: lightMP4,
		}

		var buf bytes.Buffer
		if err := t.Execute(&buf, v); err != nil {
			return "", fmt.Errorf(
				"can't execute video template for `%s`/`%s`: %w",
				pageShortname,
				videoShortname,
				err,
			)
		}

		return template.HTML(buf.String()), nil
	}, nil
}
