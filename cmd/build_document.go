package cmd

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/glacials/twos.dev/cmd/document"
	"github.com/glacials/twos.dev/cmd/transform"
)

var transformations = []document.Transformation{
	// Filename-based transformations
	transform.DiscoverParent,

	// Frontmatter-based transformations
	transform.DiscoverDates,
	transform.DiscoverType,
	transform.DiscoverShortname,
	transform.StripFrontmatter,

	// Markdown-based transformations
	transform.MarkdownToHTML,

	// HTML-based transformations
	transform.DiscoverTitle,
	transform.UpdateFeeds,
	transform.HighlightSyntax, // Beware, re-renders entire doc

	// English-based transformations
	transform.ReplaceSmartQuotes, // Must come after all HTML re-renders

	// Document-specific peculiarities
	transform.LengthenDashes,

	// Template-based transformations
	transform.PrepareTemplateVars,
	transform.AttachImageTemplateFuncs,
	transform.AttachVideoTemplateFuncs,
	transform.AttachPartials,
	transform.AttachBodyTemplate,
	transform.ExecuteTemplates,

	// Publish-based transformations
}

func buildDocument(src, dst string) error {
	d, err := document.New(src, transformations, *debug)
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
