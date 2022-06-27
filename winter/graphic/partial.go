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

// Shortname is a short, page-scoped, human-readable identifier for a graphic,
// such as "subscriptions" in the filename apple-subscriptions-dark.png.
type Shortname string

// SRC is a string that will be rendered into a graphic's src attribute (such as
// in an <img> or <video> tag).
type SRC string

// Alt is a string that will be rendered into a graphic's alt attribute, if one
// is allowed.
type Alt string

// Caption is a string that will be rendered below a graphic or set of graphics
// with a special visual treatment.
type Caption string

// LightDark returns the paths to the light and dark versions of the graphic
// specified. Any nonempty returned paths are guaranteed to be on disk and can
// have any extension in exts. If err is non-nil, at least one path is
// guaranteed nonempty.
//
// If neither a light nor a dark version is found but a neutral (unspecified)
// version is, the neutral version is returned as an unspecified one of the
// light or dark graphics; the other is empty.
//
// Graphics MUST be built before this is called, because it looks in dist/ for
// light and dark copies. TODO: Fix this, perhaps by processing graphics
// just-in-time if a page that references one is built.
func LightDark(
	pageShortname string,
	graphicShortname Shortname,
	exts map[string]struct{},
) (SRC, SRC, error) {
	if strings.Contains(string(graphicShortname), ".") {
		return "", "", fmt.Errorf(
			"graphic name %s must not have any dots",
			graphicShortname,
		)
	}
	light, err := discover(pageShortname, graphicShortname, "light", exts)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			light = ""
		} else {
			return "", "", fmt.Errorf("can't discover light graphic: %w", err)
		}
	}

	dark, err := discover(pageShortname, graphicShortname, "dark", exts)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			dark = ""
		} else {
			return "", "", fmt.Errorf("can't discover dark graphic: %w", err)
		}
	}

	if light != "" || dark != "" {
		return light, dark, nil
	}

	neutral, err := discover(pageShortname, graphicShortname, "", exts)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return "", "", fmt.Errorf(
				"found neither light, dark, nor neutral version of %s-%s",
				pageShortname,
				graphicShortname,
			)
		}
		return "", "", fmt.Errorf("can't discover neutral graphic: %w", err)
	}

	return neutral, "", nil
}

var (
	ImageExts = map[string]struct{}{
		"png":  {},
		"jpg":  {},
		"jpeg": {},
	}
	VideoExts = map[string]struct{}{
		"mov": {},
		"mp4": {},
	}
)

// Discover tries to find a graphic on the filesystem with the given arguments
// embedded in its filename, under any of the given extensions.
func discover(
	page string,
	graphic Shortname,
	suffix string,
	exts map[string]struct{},
) (SRC, error) {
	// Allow suffix to be blank without needing a trailing hyphen
	if suffix != "" {
		suffix = fmt.Sprintf("-%s", suffix)
	}

	for ext := range exts {
		path := filepath.Join(
			"img", fmt.Sprintf("%s-%s%s.%s", page, graphic, suffix, ext),
		)
		if _, err := os.Stat(filepath.Join("dist", path)); err != nil {
			if errors.Is(err, os.ErrNotExist) {
				continue
			}
			return "", fmt.Errorf("can't discover graphic: %w", err)
		}

		return SRC(path), nil
	}

	// Check again with uppercase extensions
	for ext := range exts {
		path := filepath.Join(
			"img",
			fmt.Sprintf("%s-%s%s.%s", page, graphic, suffix, strings.ToUpper(ext)),
		)
		if _, err := os.Stat(filepath.Join("dist", path)); err != nil {
			if errors.Is(err, os.ErrNotExist) {
				continue
			}
			return "", fmt.Errorf("can't discover graphic: %w", err)
		}

		return SRC(path), nil
	}

	return "", os.ErrNotExist
}
