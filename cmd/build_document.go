package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/glacials/twos.dev/cmd/document"
	"github.com/glacials/twos.dev/cmd/transform"
)

var transformations = []document.Transformation{
	// Filename-based transformations
	transform.DiscoverParent,

	// Frontmatter-based transformations
	transform.DiscoverDates,
	transform.DiscoverShortname,
	transform.StripFrontmatter,

	// Markdown-based transformations
	transform.MarkdownToHTML,

	// HTML-based transformations
	transform.DiscoverTitle,
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

	htmlFile, err := os.Create(
		filepath.Join(dst, fmt.Sprintf("%s.html", d.Shortname)),
	)
	if err != nil {
		return fmt.Errorf(
			"can't render HTML to `%s` from template: %w",
			dst,
			err,
		)
	}
	defer htmlFile.Close()

	if _, err := io.Copy(htmlFile, d.Body); err != nil {
		return err
	}

	return nil
}
