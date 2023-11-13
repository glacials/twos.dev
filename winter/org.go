package winter

import (
	"io"
	"strings"
	"time"

	"github.com/niklasfasching/go-org/org"
)

// OrgDocument represents a source file written in Org,
// with optional Go template syntax embedded in it.
//
// OrgDocument implements [Document].
//
// The OrgDocument is transitory;
// its only purpose is to create an [HTMLDocument].
type OrgDocument struct {
	// SourcePath is the path on disk to the file this Org is read from or generated from.
	// The path is relative to the working directory.
	SourcePath string

	deps map[string]struct{}
	meta *Metadata
	// next is the HTML document generated from this Org document.
	next Document
}

// NewOrgDocument creates a new document whose original source is at path src.
//
// Nothing is read from disk; src is metadata.
// To read and parse Org, call [Load].
func NewOrgDocument(src string, meta *Metadata, next Document) *OrgDocument {
	return &OrgDocument{
		SourcePath: src,

		deps: map[string]struct{}{
			"public/style.css": {},
		},
		meta: meta,
		next: next,
	}
}

func (doc *OrgDocument) DependsOn(src string) bool {
	if _, ok := doc.deps[src]; ok {
		return true
	}
	return doc.next.DependsOn(src)
}

// Load reads Org from r and loads it into doc.
//
// If called more than once, the last call wins.
func (d *OrgDocument) Load(r io.Reader) error {
	orgparser := org.New()
	orgparser.DefaultSettings["OPTIONS"] = strings.Replace(orgparser.DefaultSettings["OPTIONS"], "toc:t", "toc:nil", 1)
	orgdoc := orgparser.Parse(r, d.SourcePath)
	orgdoc.BufferSettings["OPTIONS"] = strings.Replace(
		orgdoc.BufferSettings["OPTIONS"],
		"toc:t",
		"toc:nil",
		1,
	)

	orgwriter := org.NewHTMLWriter()
	orgwriter.TopLevelHLevel = 1

	var err error
	for k, v := range orgdoc.BufferSettings {
		switch strings.ToLower(k) {
		case "category":
			d.meta.Category = v
		case "date":
			d.meta.CreatedAt, err = time.Parse("2006-01-02", v)
			if err != nil {
				return err
			}
		case "type":
			var err error
			d.meta.Kind, err = parseKind(v)
			if err != nil {
				return err
			}
		case "filename":
			d.meta.WebPath = v
		case "title":
			d.meta.Title = v
		case "toc":
			d.meta.TOC = (v == "t" || v == "true")
		case "updated":
			d.meta.UpdatedAt, err = time.Parse("2006-01-02", v)
			if err != nil {
				return err
			}
		}
	}

	htm, err := orgdoc.Write(orgwriter)
	if err != nil {
		return err
	}

	return d.next.Load(strings.NewReader(htm))
}

func (doc *OrgDocument) Metadata() *Metadata {
	return doc.meta
}

func (doc *OrgDocument) Render(w io.Writer) error {
	return doc.next.Render(w)
}
