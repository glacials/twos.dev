package img

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// lightdark returns the paths to the light and dark versions of the image or
// video specified. The returned paths, if nonempty, are guaranteed to be on
// disk and can have any image or video extension. If err is non-nil, at least
// one path is guaranteed nonempty.
//
// TODO: Images MUST be processed before this is called, because it looks in
// dist/; fix this perhaps to process images just-in-time if a page that
// references one is built.
func LightDark(page, image string) (string, string, error) {
	if strings.Contains(image, ".") {
		return "", "", fmt.Errorf("image name %s must not have any dots", image)
	}
	light, err := discoverExtension(page, image, "light")
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			light = ""
		} else {
			return "", "", fmt.Errorf("can't discover image: %w", err)
		}
	}

	dark, err := discoverExtension(page, image, "dark")
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			dark = ""
		} else {
			return "", "", fmt.Errorf("can't discover image: %w", err)
		}
	}

	if light == "" && dark == "" {
		return "", "", fmt.Errorf(
			"found neither light nor dark version of %s-%s",
			page,
			image,
		)
	}

	return light, dark, nil
}

var imgExts = map[string]struct{}{"png": {}, "jpg": {}, "jpeg": {}}

func discoverExtension(page, image, suffix string) (string, error) {
	for ext := range imgExts {
		path := filepath.Join(
			"img", fmt.Sprintf("%s-%s-%s.%s", page, image, suffix, ext),
		)
		if _, err := os.Stat(filepath.Join("dist", path)); err != nil {
			if errors.Is(err, os.ErrNotExist) {
				continue
			}
			return "", fmt.Errorf("can't discover image: %w", err)
		}

		return path, nil
	}

	// Check again with uppercase extensions
	for ext := range imgExts {
		path := filepath.Join(
			"img",
			fmt.Sprintf("%s-%s-%s.%s", page, image, suffix, strings.ToUpper(ext)),
		)
		if _, err := os.Stat(filepath.Join("dist", path)); err != nil {
			if errors.Is(err, os.ErrNotExist) {
				continue
			}
			return "", fmt.Errorf("can't discover image: %w", err)
		}

		return path, nil
	}

	return "", os.ErrNotExist
}
