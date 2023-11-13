package winter

import (
	"bytes"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/yargevad/filepathx"
)

// SaveNewURIs indexes every HTML file in dist and saves their existence to disk.
// Later, validateURIsDidNotChange can read that file and ensure no file is missing.
//
// The database file's location can be customized with winter.yml.
// It should be commited to the repository.
func (s *Substructure) SaveNewURIs(dist string) error {
	dist = filepath.Clean(dist)
	uris := map[url.URL]struct{}{}
	files, err := filepathx.Glob(filepath.Join(dist, "**", "*"))
	if err != nil {
		return fmt.Errorf("cannot glob dist dir %q: %w", dist, err)
	}
	for _, path := range files {
		path = strings.TrimPrefix(path, dist)
		if len(path) == 0 {
			continue
		}
		if stat, err := os.Stat(path); err != nil {
			return fmt.Errorf("cannot stat path to save %q: %w", path, err)
		} else if stat.IsDir() {
			continue
		}
		uris[url.URL{Path: path}] = struct{}{}
	}

	union := map[url.URL]struct{}{}

	currentPaths, err := os.ReadFile(s.cfg.Known.URIs)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			if err := os.WriteFile(s.cfg.Known.URIs, []byte(""), 0o644); err != nil {
				return fmt.Errorf("cannot make known URIs file at %q: %w", s.cfg.Known.URIs, err)
			}
			return s.SaveNewURIs(dist)
		}
		return fmt.Errorf("cannot get existing URLs: %w", err)
	}
	for _, uri := range bytes.Split(currentPaths, []byte("\n")) {
		if len(uri) == 0 {
			continue
		}
		union[url.URL{Path: string(uri)}] = struct{}{}
	}
	for uri := range uris {
		union[uri] = struct{}{}
	}
	list := make([]string, 0, len(union))
	for uri := range union {
		list = append(list, uri.Path)
	}
	slices.Sort(list)
	f, err := os.Create(s.cfg.Known.URIs)
	if err != nil {
		return fmt.Errorf("cannot open known URIs file %q for writing: %w", s.cfg.Known.URIs, err)
	}
	for _, uri := range list {
		if _, err := f.WriteString(uri); err != nil {
			return fmt.Errorf("cannot write URI %q to known URIs file %q: %w", uri, s.cfg.Known.URIs, err)
		}
		if _, err := f.WriteString("\n"); err != nil {
			return fmt.Errorf("cannot write newline to known URIs file %q: %w", s.cfg.Known.URIs, err)
		}
	}
	return nil
}

// validateURIsDidNotChange returns an error if this build neglected to produce
// an HTML file that was previously present on the site.
//
// To update the list validateURIsDidNotChange uses, run:
//
//	winter freeze
//
// For more information about the "cool URIs don't change" rule, see:
// https://www.w3.org/Provider/Style/URI
func (s *Substructure) validateURIsDidNotChange(dist string) error {
	paths, err := os.ReadFile(s.cfg.Known.URIs)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			if err := os.WriteFile(s.cfg.Known.URIs, []byte(""), 0o644); err != nil {
				return fmt.Errorf("cannot create new known URIs file: %w", err)
			}
			return s.validateURIsDidNotChange(dist)
		}
		return fmt.Errorf("cannot read known URLs file: %w", err)
	}
	changedURIs := []string{}
	for _, pathBytes := range bytes.Split(paths, []byte{'\n'}) {
		if len(pathBytes) == 0 {
			continue
		}
		_, err := os.Stat(filepath.Join(dist, string(pathBytes)))
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				changedURIs = append(changedURIs, string(pathBytes))
			} else {
				return fmt.Errorf("cannot stat %q: %w", pathBytes, err)
			}
		}
	}
	if len(changedURIs) > 0 {
		return fmt.Errorf("cool URIs do not change, but these ones were removed by this build:\n\n- %s", strings.Join(changedURIs, "\n- "))
	}
	return nil
}
