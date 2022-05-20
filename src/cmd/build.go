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
	"os"
	"path/filepath"

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

		if err := copy.Copy(filepath.Join(sourceDir, "img"), filepath.Join(destinationDir, "img")); err != nil {
			return fmt.Errorf("can't copy original images to build dir: %w", err)
		}

		if err := genThumbnails(
			filepath.Join(sourceDir, "img"),
			filepath.Join(destinationDir, "img"),
		); err != nil {
			return fmt.Errorf("can't generate thumbnails: %w", err)
		}

		// TODO: Generate html pages (one per image) for bigger image viewing
		// TODO: Move some pre-Markdoc things in here

		return nil
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)
}

func genThumbnails(sourceDir, destinationDir string) error {
	if err := filepath.WalkDir(sourceDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("can't walk path `%s`: %w", path, err)
		}

		if d.IsDir() {
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
		w, h := 150, ((150*p.X/p.Y)+1)&-1
		destinationImage := image.NewRGBA(image.Rect(0, 0, w, h))

		draw.CatmullRom.Scale(
			destinationImage,
			image.Rectangle{image.Point{0, 0}, image.Point{150, 150}},
			sourceImage,
			image.Rectangle{image.Point{0, 0}, sourceImage.Bounds().Size()},
			draw.Over,
			nil,
		)

		relativePath, err := filepath.Rel(sourceDir, path)
		if err != nil {
			return fmt.Errorf("can't get relative path of `%s`: %w", path, err)
		}

		destinationPath := filepath.Join(filepath.Join(destinationDir, "thumb", relativePath))
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
