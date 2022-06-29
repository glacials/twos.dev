package cmd

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"twos.dev/winter/document"
	"twos.dev/winter/transform"
)

var transformations = []document.Transformation{
	// Filename-based transformations
	transform.DiscoverParent,

	// Frontmatter-based transformations
	transform.ExtractFrontmatter,

	// Markup-based transformations
	transform.RenderMarkdown,
	transform.RenderLaTeX,

	// Pre-tempalte HTML-based transformations
	transform.DiscoverTitle,
	transform.RenderTOC, // Beware, re-renders entire doc

	// English-based transformations
	transform.UnescapeHTML,

	// Document-specific peculiarities
	transform.LengthenDashes,

	// Template-based transformations
	transform.PrepareTemplateVars,
	transform.AttachImageTemplateFuncs,
	transform.AttachVideoTemplateFuncs,
	transform.AttachPartials,
	transform.AttachBodyTemplate,
	transform.ExecuteTemplates,

	// Post-template HTML-based transformations
	transform.HighlightSyntax,    // Beware, re-renders entire doc
	transform.ReplaceSmartQuotes, // Must come after all HTML re-renders

	// Publish-based transformations
	transform.UpdateFeeds,
}

func buildDocument(src, dst string, cfg Config) error {
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
