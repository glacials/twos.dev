package frontmatter

import (
	"fmt"
	"io"
	"strings"
	"time"

	fm "github.com/adrg/frontmatter"
)

type internalMatter struct {
	TOC       bool      `yaml:"toc"`
	Type      string    `yaml:"type"`
	Filename  string    `yaml:"filename"`
	CreatedAt time.Time `yaml:"date"`
	UpdatedAt time.Time `yaml:"updated"`
}

type Type int

const (
	// DraftType designates a document that should not be linked to from anywhere.
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

func (t Type) IsDraft() bool   { return t == DraftType }
func (t Type) IsPost() bool    { return t == PostType }
func (t Type) IsPage() bool    { return t == PageType }
func (t Type) IsGallery() bool { return t == GalleryType }

type Matter struct {
	// Type is the class of document, specified by the author.
	Type Type

	// Shortname is a human-readable, one-word, all-lowercase tag for the
	// document. The shortname is the part of the filename before the extension. A
	// document's shortname must never change after it is published.
	Shortname string

	// TOC causes a table of contents to be generated if true.
	TOC bool

	CreatedAt time.Time
	UpdatedAt time.Time
}

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
