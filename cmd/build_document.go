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
	transform.ReplaceSmartQuotes,

	// Template-based transformations
	transform.PrepareTemplateVars,
	transform.AttachImageTemplateFuncs,
	transform.AttachVideoTemplateFuncs,
	transform.AttachPartials,
	transform.AttachBodyTemplate,
	transform.ExecuteTemplates,
}

func buildDocument(src, dst string) error {
	d, err := document.New(src, transformations)
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
