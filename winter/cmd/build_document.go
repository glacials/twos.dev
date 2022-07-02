package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"twos.dev/winter"
	"twos.dev/winter/document"
	"twos.dev/winter/transform"
)

var transformations = []document.Transformation{
	// English-based transformations
	transform.UnescapeHTML,

	// Template-based transformations
	transform.AttachVideoTemplateFunc,
	transform.AttachPostsTemplateFunc,

	// Post-template HTML-based transformations
	transform.HighlightSyntax,    // Beware, re-renders entire doc
	transform.ReplaceSmartQuotes, // Must come after all HTML re-renders

	// Publish-based transformations
	transform.UpdateFeeds,
}

func buildDocument(src, dst string, cfg winter.Config) error {
	d, err := document.New(src, transformations, cfg.Debug)
	if err != nil {
		return err
	}

	d, err = d.Transform()
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(filepath.Join(dst, fmt.Sprintf("%s.html", d.Shortname)), d.Body, 0644); err != nil {
		return err
	}

	return nil
}
