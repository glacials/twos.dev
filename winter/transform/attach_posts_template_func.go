package transform

import (
	"html/template"

	"twos.dev/winter/document"
)

// AttachPostsTemplateFunc makes the posts function available to the document's
// templates. It returns a slice of all documents with type document.PostType.
//
// AttachPostsTemplateFunc implements document.Transformation.
func AttachPostsTemplateFunc(d document.Document) (document.Document, error) {
	_ = d.Template.Funcs(template.FuncMap{"posts": posts})
	return d, nil
}

func posts() ([]document.Document, error) {
	return []document.Document{}, nil
}
