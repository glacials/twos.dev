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
	"fmt"
	"strings"
	"time"
)

type yearMonth time.Time

func (t *yearMonth) UnmarshalYAML(b []byte) error {
	s := strings.Trim(string(b), `"`)

	nt, err := time.Parse("2006-01", s)
	if err != nil {
		return fmt.Errorf("can't parse `%s` as time: %w", s, err)
	}

	*t = yearMonth(nt)
	return nil
}

func (t *yearMonth) MarshalYAML() ([]byte, error) {
	return []byte(t.String()), nil
}

func (t *yearMonth) Format(layout string) string {
	return time.Time(*t).Format(layout)
}

func (t *yearMonth) IsZero() bool {
	return time.Time(*t).IsZero()
}

func (t *yearMonth) String() string {
	return time.Time(*t).Format("2006 January")
}
