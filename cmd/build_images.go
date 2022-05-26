package cmd

import (
	"fmt"
	"html/template"
	"image"
	"image/jpeg"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/otiai10/copy"
	"golang.org/x/image/draw"
)

func buildImages(sourceDir, destinationDir string) error {
	imageSourceDir := filepath.Join(sourceDir, "img")
	imageDestinationDir := filepath.Join(destinationDir, "img")

	// Copy original images to build dir
	if err := copy.Copy(imageSourceDir, imageDestinationDir); err != nil {
		return fmt.Errorf("can't copy original images to build dir: %w", err)
	}

	// Generate thumbnails and place into build dir
	if err := genThumbnails(imageSourceDir, imageDestinationDir, 300); err != nil {
		return fmt.Errorf("can't generate thumbnails: %w", err)
	}

	// Create gallery pages for each image
	if err := genImageContainers(sourceDir, destinationDir); err != nil {
		return fmt.Errorf("can't generate image container pages: %w", err)
	}

	return nil
}

// genThumbnails makes thumbnails from every image (recursively) in sourceDir, copying
// them to equivalent paths in a "thumb" subdirectory of destinationDir.
//
// The thumbnails have the given width. Height is automatically calculated to maintain
// image ratio.
func genThumbnails(imageSourceDir, imageDestinationDir string, newWidth int) error {
	if err := filepath.WalkDir(imageSourceDir, func(path string, d fs.DirEntry, err error) error {
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

		relativePath, err := filepath.Rel(imageSourceDir, path)
		if err != nil {
			return fmt.Errorf("can't get relative path of `%s`: %w", path, err)
		}

		destinationPath := filepath.Join(imageDestinationDir, "thumb", relativePath)
		if err := os.MkdirAll(filepath.Dir(destinationPath), 0755); err != nil {
			return fmt.Errorf("can't make thumbnail directory in path `%s`: %w", imageDestinationDir, err)
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

func genImageContainers(imageSourceDir, imageDestinationDir string) error {
	templateHTML, err := ioutil.ReadFile(filepath.Join(imageSourceDir, imageContainerTemplatePath))
	if err != nil {
		return fmt.Errorf("can't read image container template: %w", err)
	}

	t, err := template.New("imgcontainer").Parse(string(templateHTML))
	if err != nil {
		return fmt.Errorf("can't create imgcontainer template: %w", err)
	}

	var imgs []string

	if err := filepath.WalkDir(imageSourceDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("can't walk path `%s` for imgcontainers: %w", path, err)
		}

		if d.IsDir() {
			return nil
		}

		// TODO: Support other filetypes
		matched, err := filepath.Match("*.jpg", strings.ToLower(d.Name()))
		if err != nil {
			return fmt.Errorf("can't match against filename `%s` for imgcontainers: %w", d.Name(), err)
		}
		if !matched {
			return nil
		}

		imgs = append(imgs, path)

		return nil
	}); err != nil {
		return fmt.Errorf("can't perform imgcontainer loop: %w", err)
	}

	for i, path := range imgs {
		relativeImagePath, err := filepath.Rel(imageSourceDir, path)
		if err != nil {
			return fmt.Errorf("can't get relative path of `%s`: %w", path, err)
		}

		f, err := os.Create(fmt.Sprintf("%s.html", filepath.Join(imageDestinationDir, relativeImagePath)))
		if err != nil {
			return fmt.Errorf("can't create imgcontainer file for `%s`: %w", path, err)
		}

		// TODO: Make prev/next only go through images in the same directory
		var prevLink, nextLink string
		if i > 0 {
			prevLink = fmt.Sprintf("%s.html", filepath.Base(imgs[i-1]))
		}
		if i < len(imgs)-1 {
			nextLink = fmt.Sprintf("%s.html", filepath.Base(imgs[i+1]))
		}

		v := imageContainerVars{
			PrevLink: prevLink,
			CurImage: filepath.Base(path),
			NextLink: nextLink,
		}

		if err := t.Execute(f, v); err != nil {
			return fmt.Errorf("can't execute imgcontainer template: %w", err)
		}
	}

	return nil
}
