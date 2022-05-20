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
package cmd

import (
	"fmt"
	"image"
	"image/jpeg"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/otiai10/copy"
	"github.com/spf13/cobra"
	"golang.org/x/image/draw"
)

// buildCmd represents the build command
var buildCmd = &cobra.Command{
	Use:   "build SOURCE DESTINATION",
	Short: "Build twos.dev",
	Long:  `Build twos.dev from the given source directory into the given destination directory.`,
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		sourceDir, err := filepath.Abs(args[0])
		if err != nil {
			return fmt.Errorf("can't resolve path `%s`: %w", args[0], err)
		}
		destinationDir, err := filepath.Abs(args[1])
		if err != nil {
			return fmt.Errorf("can't resolve path `%s`: %w", args[1], err)
		}

		if err := buildImages(sourceDir, destinationDir); err != nil {
			return fmt.Errorf("can't build images: %w", err)
		}

		// TODO: Generate html pages (one per image) for bigger image viewing
		// TODO: Move some pre-Markdoc things in here

		return nil
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)
}

func buildImages(sourceDir, destinationDir string) error {
	imageSourceDir := filepath.Join(sourceDir, "img")
	imageDestinationDir := filepath.Join(destinationDir, "img")

	// Copy original images to build dir
	if err := copy.Copy(sourceDir, destinationDir); err != nil {
		return fmt.Errorf("can't copy original images to build dir: %w", err)
	}

	// Generate thumbnails and place into build dir
	if err := genThumbnails(imageSourceDir, imageDestinationDir, 300); err != nil {
		return fmt.Errorf("can't generate thumbnails: %w", err)
	}

	// Replace any <img> tag with the thumbnail version
	if err := filepath.WalkDir(destinationDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("can't walk path `%s`: %w", path, err)
		}

		if d.IsDir() {
			return nil
		}

		matched, err := filepath.Match("*.html", strings.ToLower(d.Name()))
		if err != nil {
			return fmt.Errorf("can't match against filename `%s`: %w", d.Name(), err)
		}
		if !matched {
			return nil
		}

		originalHTML, err := ioutil.ReadFile(path)
		if err != nil {
			return fmt.Errorf("can't read HTML file `%s`: %w", path, err)
		}

		re, err := regexp.Compile("(<img.+src=[\"']img\\/)(.*\\.[jJ][pP][gG][\"'].*>)")
		if err != nil {
			return fmt.Errorf("can't compile regexp: %w", err)
		}

		replacedHTML := re.ReplaceAll(originalHTML, []byte("${1}thumb/${2}"))

		if err := ioutil.WriteFile(path, replacedHTML, 0); err != nil {
			return fmt.Errorf("can't write thumbnail replacements to file: %w", err)
		}

		return nil
	}); err != nil {
		return fmt.Errorf("can't iterate over source to replace images with thumbnails: %s", err)
	}

	return nil
}

// genThumbnails makes thumbnails from every image (recursively) in sourceDir, copying
// them to equivalent paths in a "thumb" subdirectory of destinationDir.
//
// The thumbnails have the given width. Height is automatically calculated to maintain
// image ratio.
func genThumbnails(sourceDir, destinationDir string, newWidth int) error {
	if err := filepath.WalkDir(sourceDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("can't walk path `%s`: %w", path, err)
		}

		if d.IsDir() {
			return nil
		}

		// TODO: Support other filetypes
		matched, err := filepath.Match("*.jpg", strings.ToLower(d.Name()))
		if err != nil {
			return fmt.Errorf("can't match against filename `%s`: %w", d.Name(), err)
		}
		if !matched {
			return nil
		}

		path, err = filepath.Abs(path)
		if err != nil {
			return fmt.Errorf("can't resolve path `%s`: %w", path, err)
		}

		sourceFile, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("can't open image at path `%s`: %w", path, err)
		}

		sourceImage, err := jpeg.Decode(sourceFile)
		if err != nil {
			return fmt.Errorf("can't decode image at path `%s` (maybe not an image?): %w", path, err)
		}
		p := sourceImage.Bounds().Size()
		w, h := newWidth, ((newWidth*p.X/p.Y)+1)&-1
		destinationImage := image.NewRGBA(image.Rect(0, 0, w, h))

		draw.CatmullRom.Scale(
			destinationImage,
			image.Rectangle{image.Point{0, 0}, image.Point{newWidth, newWidth}},
			sourceImage,
			image.Rectangle{image.Point{0, 0}, sourceImage.Bounds().Size()},
			draw.Over,
			nil,
		)

		relativePath, err := filepath.Rel(sourceDir, path)
		if err != nil {
			return fmt.Errorf("can't get relative path of `%s`: %w", path, err)
		}

		destinationPath := filepath.Join(destinationDir, "thumb", relativePath)
		if err := os.MkdirAll(filepath.Dir(destinationPath), 0755); err != nil {
			return fmt.Errorf("can't make thumbnail directory in path `%s`: %w", destinationDir, err)
		}

		destinationFile, err := os.Create(destinationPath)
		if err != nil {
			return fmt.Errorf("can't create thumbnail file for image at path `%s`: %w", path, err)
		}
		defer destinationFile.Close()

		if err := jpeg.Encode(destinationFile, destinationImage, nil); err != nil {
			return fmt.Errorf("can't encode to destination file at path `%s`: %w", destinationPath, err)
		}

		return nil
	}); err != nil {
		return fmt.Errorf("can't perform thumbnail loop: %w", err)
	}

	return nil
}
