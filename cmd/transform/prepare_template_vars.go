package transform

import (
	"html/template"

	"github.com/glacials/twos.dev/cmd/document"
)

// PrepareTemplateVars massages and converts various variables on the document
// into formats usable by the template engine, and stores the converted results
// as a field on the document.
//
// PrepareTemplateVars implements document.Transformation.
func PrepareTemplateVars(d document.Document) (document.Document, error) {
	d.TemplateVars.Body = template.HTML(string(d.Body))
	d.TemplateVars.Parent = d.Parent
	d.TemplateVars.SourceURL = d.SourceURL
	d.TemplateVars.Shortname = d.Shortname
	d.TemplateVars.Title = d.Title

	d.TemplateVars.CreatedAt = d.CreatedAt
	d.TemplateVars.UpdatedAt = d.UpdatedAt

	return d, nil
}
