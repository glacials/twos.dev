package frontmatter

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"time"

	fm "github.com/adrg/frontmatter"
)

type internalMatter struct {
	Filename string `yaml:"filename"`
	// Date is an alias for CreatedAt.
	Date      time.Time `yaml:"date"`
	CreatedAt time.Time `yaml:"created"`
	UpdatedAt time.Time `yaml:"updated"`
}

type Matter struct {
	Shortname string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// ParseDeprecated returns the parsed frontmatter and the remaining non-frontmatter
// content from r.
//
// Deprecated: Use Parse instead.
func ParseDeprecated(r io.Reader) (Matter, []byte, error) {
	var matter internalMatter
	body, err := fm.Parse(r, &matter)
	if err != nil {
		return Matter{}, []byte{}, fmt.Errorf("can't parse frontmatter: %w", err)
	}

	if matter.CreatedAt.IsZero() {
		matter.CreatedAt = matter.Date
	}

	return Matter{
		Shortname: matter.Filename,
		CreatedAt: time.Time(matter.CreatedAt),
		UpdatedAt: time.Time(matter.UpdatedAt),
	}, body, nil
}

// Parse returns the parsed frontmatter and the remaining non-frontmatter
// content from r.
func Parse(r io.Reader) (Matter, io.Reader, error) {
	var matter internalMatter
	body, err := fm.Parse(r, &matter)
	if err != nil {
		return Matter{}, nil, fmt.Errorf("can't parse frontmatter: %w", err)
	}

	if matter.CreatedAt.IsZero() {
		matter.CreatedAt = matter.Date
	}

	filenameParts := strings.Split(matter.Filename, ".")

	return Matter{
		Shortname: filenameParts[0],
		CreatedAt: matter.CreatedAt,
		UpdatedAt: matter.UpdatedAt,
	}, bytes.NewBuffer(body), nil
}
