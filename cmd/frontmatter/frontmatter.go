package frontmatter

import (
	"fmt"
	"io"
	"strings"
	"time"

	fm "github.com/adrg/frontmatter"
	"github.com/glacials/twos.dev/cmd/document"
)

type internalMatter struct {
	Type     string
	Filename string `yaml:"filename"`
	// Date is an alias for CreatedAt.
	Date      time.Time `yaml:"date"`
	CreatedAt time.Time `yaml:"created"`
	UpdatedAt time.Time `yaml:"updated"`
}

type Matter struct {
	Type        document.Type
	Shortname   string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	PublishedAt time.Time
}

// Parse returns the parsed frontmatter and the remaining non-frontmatter
// content from r.
func Parse(r io.Reader) (Matter, []byte, error) {
	var matter internalMatter

	body, err := fm.Parse(r, &matter)
	if err != nil {
		return Matter{}, nil, fmt.Errorf("can't parse frontmatter: %w", err)
	}

	var t document.Type
	switch matter.Type {
	case "", "draft":
		t = document.DraftType
	case "post":
		t = document.PostType
	case "page":
		t = document.PageType
	default:
		return Matter{}, nil, fmt.Errorf("invalid document type %s", matter.Type)
	}

	if matter.CreatedAt.IsZero() {
		matter.CreatedAt = matter.Date
	}

	filenameParts := strings.Split(matter.Filename, ".")

	return Matter{
		Type:      t,
		Shortname: filenameParts[0],
		CreatedAt: matter.CreatedAt,
		UpdatedAt: matter.UpdatedAt,
	}, body, nil
}
