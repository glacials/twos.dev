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
