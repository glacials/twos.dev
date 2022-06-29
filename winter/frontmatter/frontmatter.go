package frontmatter

import (
	"fmt"
	"io"
	"strings"
	"time"

	fm "github.com/adrg/frontmatter"
)

// internalMatter is a struct that is used to unmarshal frontmatter. It is
// separated from Matter so that any backwards-compatible idiosyncrasies aren't
// exposed in the package's outward API.
type internalMatter struct {
	TOC       bool      `yaml:"toc"`
	Type      string    `yaml:"type"`
	Filename  string    `yaml:"filename"`
	CreatedAt time.Time `yaml:"date"`
	UpdatedAt time.Time `yaml:"updated"`
}

// Matter is the frontmatter of a document.
type Matter struct {
	// Type is the class of document, specified by the author.
	Type Type

	// Shortname is a human-readable, one-word, all-lowercase tag for the
	// document. The shortname is the part of the filename before the extension. A
	// document's shortname must never change after it is published.
	Shortname string

	// TOC causes a table of contents to be generated from headings if true.
	TOC bool

	// CreatedAt is the publish date of the document as specified by the date
	// frontmatter field.
	CreatedAt time.Time

	// UpdatedAt is the last time the document was updated as specified by the
	// updated frontmatter field.
	UpdatedAt time.Time
}

// Type is the style of document being handled, as specified by its frontmatter.
// A zero value indicates that the document is of type DraftType.
type Type int

const (
	// DraftType designates a document that should not be linked to from anywhere,
	// but should still be published.
	DraftType Type = iota

	// PostType designates a document with a front-facing post date and listed by
	// recency in a blog-style page.
	PostType

	// PageType designates a document that is manually linked to from somewhere on
	// the site. It has a publish date, but it is less pronounced in display than
	// posts.
	PageType

	// GaleryType designates a document whose sole purpose is to display a large
	// image.
	GalleryType
)

// IsDraft returns whether the document type is DraftType. This function exists
// to be used by templates.
func (t Type) IsDraft() bool { return t == DraftType }

// IsPost returns whether the document type is PostType. This function exists to
// be used by templates.
func (t Type) IsPost() bool { return t == PostType }

// IsPage returns whether the document type is PageType. This function exists to
// be used by templates.
func (t Type) IsPage() bool { return t == PageType }

// IsGallery returns whether the document type is GalleryType. This function
// exists to be used by templates.
func (t Type) IsGallery() bool { return t == GalleryType }

// Parse returns the parsed frontmatter and the remaining non-frontmatter
// content from r.
func Parse(r io.Reader) (Matter, []byte, error) {
	var matter internalMatter

	body, err := fm.Parse(r, &matter)
	if err != nil {
		return Matter{}, nil, fmt.Errorf("can't parse frontmatter: %w", err)
	}

	var t Type
	switch matter.Type {
	case "", "draft":
		t = DraftType
	case "post":
		t = PostType
	case "page":
		t = PageType
	case "gallery":
		t = GalleryType
	default:
		return Matter{}, nil, fmt.Errorf("invalid document type %s", matter.Type)
	}

	filenameParts := strings.Split(matter.Filename, ".")

	return Matter{
		TOC:       matter.TOC,
		Type:      t,
		Shortname: filenameParts[0],
		CreatedAt: matter.CreatedAt,
		UpdatedAt: matter.UpdatedAt,
	}, body, nil
}
