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
	imgPartial = "imgs"
)

// imgsPartialVars are the template variables given to
// src/templates/_imgs.html.tmpl to render multiple images inline. At least one
// of its {Light,Dark} fields must be set to an image path.
type imgsPartialVars struct {
	Images  []imageVars
	Caption graphic.Caption
}

type imageVars struct {
	Alt   graphic.Alt
	Light graphic.SRC
	Dark  graphic.SRC
}

// AttachImageTemplateFuncs makes the img and imgs functions available to the
// document's templates.
func AttachImageTemplateFuncs(d document.Document) (document.Document, error) {
	img, err := img(d.Shortname)
	if err != nil {
		return document.Document{}, err
	}

	imgs, err := imgs(d.Shortname)
	if err != nil {
		return document.Document{}, err
	}

	_ = d.Template.Funcs(
		template.FuncMap{"img": img, "imgs": imgs},
	)

	return d, nil
}

// img returns a function that can be inserted into a template's FuncMap for
// calling by the template itself. The returned function takes an image
// shortname, alt text, and caption; the image is rendered with the caption
// below it. If light and dark mode images are both present, the correct one
// will be rendered.
//
// The image must be named in the format SHORTNAME-IMAGE-STYLE.EXTENSION, where
// SHORTNAME is the page shortname (e.g. "apple" for apple.html), IMAGE is any
// arbitrary string that must be passed to the returned function as its imgsrc,
// STYLE is one of "light" or "dark", and EXTENSION is one of "png" or "jpg".
//
// The given shortname must be the page shortname the image will appear on, or
// the rendered image won't point to the right URL.
func img(
	pageShortname string,
) (func(
	graphic.Shortname,
	graphic.Alt,
	graphic.Caption,
) (template.HTML, error), error) {
	partial, err := ioutil.ReadFile(partialPath(imgPartial))
	if err != nil {
		return nil, err
	}

	t := template.New(imgPartial)
	if _, err := t.Parse(string(partial)); err != nil {
		return nil, err
	}

	return func(
		imageShortname graphic.Shortname,
		alt graphic.Alt,
		caption graphic.Caption,
	) (template.HTML, error) {
		light, dark, err := graphic.LightDark(
			pageShortname,
			imageShortname,
			graphic.ImageExts,
		)
		if err != nil {
			return "", fmt.Errorf("can't process image: %w", err)
		}

		v := imgsPartialVars{
			Caption: caption,
			Images: []imageVars{
				{
					Alt:   alt,
					Dark:  dark,
					Light: light,
				},
			},
		}

		var buf bytes.Buffer
		if err := t.Lookup("imgs").Execute(&buf, v); err != nil {
			return "", fmt.Errorf(
				"can't execute img template for `%s`: %w",
				imageShortname,
				err,
			)
		}

		return template.HTML(buf.String()), nil
	}, nil
}

// imgs returns a function that can be inserted into a template's FuncMap for
// calling by the template itself. The returned function takes a pair of image
// shortnames, a pair of alt texts, and one caption; the images are rendered
// next to each other with the caption below them. If light and dark mode images
// are both present, the correct ones will be rendered.
//
// The images must be named in the format SHORTNAME-IMAGE-STYLE.EXTENSION, where
// SHORTNAME is the page shortname (e.g. "apple" for apple.html), IMAGE is any
// arbitrary string that must be passed to the returned function as its imgsrc,
// STYLE is one of "light" or "dark", and EXTENSION is one of "png" or "jpg".
//
// The given shortname must be the page shortname the images will appear on, or
// the rendered images won't point to the right URLs.
func imgs(
	shortname string,
) (func(
	graphic.Shortname,
	graphic.Alt,
	graphic.Shortname,
	graphic.Alt,
	graphic.Caption,
) (template.HTML, error), error) {
	partial, err := ioutil.ReadFile(partialPath(imgPartial))
	if err != nil {
		return nil, err
	}

	t := template.New(imgPartial)
	if _, err := t.Parse(string(partial)); err != nil {
		return nil, err
	}

	return func(
		imageA graphic.Shortname,
		altA graphic.Alt,
		imageB graphic.Shortname,
		altB graphic.Alt,
		caption graphic.Caption,
	) (template.HTML, error) {
		lightA, darkA, err := graphic.LightDark(
			shortname,
			imageA,
			graphic.ImageExts,
		)
		if err != nil {
			return "", fmt.Errorf("can't process image: %w", err)
		}

		lightB, darkB, err := graphic.LightDark(
			shortname,
			imageB,
			graphic.ImageExts,
		)
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

		var buf bytes.Buffer
		if err := t.Execute(&buf, v); err != nil {
			return "", fmt.Errorf(
				"can't execute imgs template for `%s`/`%s`: %w",
				imageA,
				imageB,
				err,
			)
		}

		return template.HTML(buf.String()), nil
	}, nil
}
