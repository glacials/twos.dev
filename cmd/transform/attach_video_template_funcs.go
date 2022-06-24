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
	videoPartial = "videos"
)

// videoPartialVars are the template variables given to
// src/templates/_video.html.tmpl to render a video inline. At least one of its
// {Light,Dark}{MOV,MP4} fields must be set to a video path.
type videosPartialVars struct {
	Videos  []videoVars
	Caption graphic.Caption
}

type videoVars struct {
	LightMOV graphic.SRC
	LightMP4 graphic.SRC
	DarkMOV  graphic.SRC
	DarkMP4  graphic.SRC
}

// AttachVideoTemplateFuncs makes the video and videos functions available to
// the document's templates. video is an alias for videos.
func AttachVideoTemplateFuncs(d document.Document) (document.Document, error) {
	videos, err := videos(d.Shortname)
	if err != nil {
		return document.Document{}, err
	}

	_ = d.Template.Funcs(
		template.FuncMap{"video": videos, "videos": videos},
	)

	return d, nil
}

func videos(
	pageShortname string,
) (func(graphic.Caption, ...graphic.Shortname) (template.HTML, error), error) {
	partial, err := ioutil.ReadFile(partialPath(videoPartial))
	if err != nil {
		return nil, err
	}

	t := template.New(videoPartial)
	if _, err := t.Parse(string(partial)); err != nil {
		return nil, err
	}

	return func(videoCaption graphic.Caption, videoShortnames ...graphic.Shortname) (template.HTML, error) {
		v := videosPartialVars{
			Caption: videoCaption,
		}

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
