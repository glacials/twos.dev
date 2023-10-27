// Package graphic provides utility functions for finding and building images or
// videos on the filesystem.
package graphic

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var (
	ImageExts = map[string]struct{}{
		"png":  {},
		"jpg":  {},
		"jpeg": {},
		"svg":  {},
	}
	VideoExts = map[string]struct{}{
		"mov": {},
		"mp4": {},
	}
)

// SRC is a string that will be rendered into a graphic's src attribute (such as
// in an <img> or <video> tag).
type SRC string

// Alt is a string that will be rendered into a graphic's alt attribute, if one
// is allowed.
type Alt string

// Caption is a string that will be rendered below a graphic or set of graphics
// with a special visual treatment.
type Caption string

// Discover returns the path to the graphic specified.
// If none exists, os.ErrNotExist is returned.
//
// Graphics MUST be built before this is called, because it looks in dist/.
// TODO: Fix this, perhaps by processing graphics just-in-time if a page that references one is built.
func Discover(
	src SRC,
	exts map[string]struct{},
) (SRC, error) {
	for ext := range exts {
		path := fmt.Sprintf("%s.%s", src, ext)
		if _, err := os.Stat(filepath.Join("public", path)); err != nil {
			if errors.Is(err, os.ErrNotExist) {
				continue
			}
			return "", fmt.Errorf("can't discover graphicShortname: %w", err)
		}

		return SRC(path), nil
	}

	// Check again with uppercase extensions
	for ext := range exts {
		path := fmt.Sprintf("%s.%s", src, strings.ToUpper(ext))
		if _, err := os.Stat(filepath.Join("public", path)); err != nil {
			if errors.Is(err, os.ErrNotExist) {
				continue
			}
			return "", fmt.Errorf("can't discover graphic: %w", err)
		}

		return SRC(path), nil
	}

	return "", os.ErrNotExist
}
