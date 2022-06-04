package cmd

import (
	"fmt"
	"html/template"
	"image"
	"image/jpeg"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"golang.org/x/exp/slices"
	"golang.org/x/image/draw"
)

func imageBuilder(src, dst string) error {
	relsrc, err := filepath.Rel("src", src)
	if err != nil {
		return fmt.Errorf("can't get relpath for image `%s`: %w", src, err)
	}

	imgdst := filepath.Join(dst, relsrc)
	thmdst := filepath.Join(dst, "thumb", relsrc)

	if err := os.MkdirAll(filepath.Dir(imgdst), 0755); err != nil {
		return fmt.Errorf("can't mkdir `%s`: %w", filepath.Dir(imgdst), err)
	}

	srcf, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("can't read `%s`: %w", src, err)
	}

	dstf, err := os.Create(imgdst)
	if err != nil {
		return fmt.Errorf("can't write image `%s`: %w", imgdst, err)
	}

	if _, err := io.Copy(dstf, srcf); err != nil {
		return fmt.Errorf("can't copy `%s` to `%s`: %w", src, imgdst, err)
	}

	// Generate thumbnails and place into build dir
	if err := genThumbnail(src, filepath.Dir(thmdst), 300); err != nil {
		return fmt.Errorf("can't generate thumbnails: %w", err)
	}

	if err := genGalleryPage(src, fmt.Sprintf("%s.html", imgdst)); err != nil {
		return fmt.Errorf("can't generate image container pages: %w", err)
	}

	return nil
}

// genThumbnail makes thumbnails from every image (recursively) in src, copying
// them to equivalent paths in a "thumb" subdirectory of dst.
//
// The thumbnails have the given width. Height is automatically calculated to
// maintain ratio.
func genThumbnail(src, dst string, width int) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("can't open image at path `%s`: %w", src, err)
	}

	sourceImage, err := jpeg.Decode(sourceFile)
	if err != nil {
		return fmt.Errorf("can't decode image at path `%s` (maybe not an image?): %w", src, err)
	}
	p := sourceImage.Bounds().Size()
	w, h := width, ((width*p.X/p.Y)+1)&-1
	destinationImage := image.NewRGBA(image.Rect(0, 0, w, h))

	draw.CatmullRom.Scale(
		destinationImage,
		image.Rectangle{image.Point{0, 0}, image.Point{width, width}},
		sourceImage,
		image.Rectangle{image.Point{0, 0}, sourceImage.Bounds().Size()},
		draw.Over,
		nil,
	)

	relativePath, err := filepath.Rel(commonPathAncestor(src, dst), src)
	if err != nil {
		return fmt.Errorf("can't get relative path of `%s`: %w", src, err)
	}

	destinationPath := filepath.Join(dst, "thumb", relativePath)
	if err := os.MkdirAll(filepath.Dir(destinationPath), 0755); err != nil {
		return fmt.Errorf("can't make thumbnail directory in path `%s`: %w", dst, err)
	}

	destinationFile, err := os.Create(destinationPath)
	if err != nil {
		return fmt.Errorf("can't create thumbnail file for image at path `%s`: %w", src, err)
	}
	defer destinationFile.Close()

	if err := jpeg.Encode(destinationFile, destinationImage, nil); err != nil {
		return fmt.Errorf("can't encode to destination file at path `%s`: %w", destinationPath, err)
	}

	return nil
}

func genGalleryPage(src, dst string) error {
	templateHTML, err := ioutil.ReadFile(imageContainerTemplatePath)
	if err != nil {
		return fmt.Errorf("can't read image container template: %w", err)
	}

	t, err := template.New("imgcontainer").Parse(string(templateHTML))
	if err != nil {
		return fmt.Errorf("can't create imgcontainer template: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return fmt.Errorf("can't create imgcontainer directory `%s`: %w", filepath.Dir(dst), err)
	}

	f, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("can't create imgcontainer file for `%s`: %w", src, err)
	}

	files, err := filepath.Glob(filepath.Join(filepath.Dir(src), "*.[jJ][pP][gG]"))
	if err != nil {
		return fmt.Errorf("can't look into image directory `%s` for ordering: %w", filepath.Dir(src), err)
	}

	i := slices.IndexFunc(
		files,
		func(file string) bool { return filepath.Base(file) == filepath.Base(src) },
	)

	var prevLink, nextLink string
	if i > 0 {
		prevLink = fmt.Sprintf("%s.html", filepath.Base(files[i-1]))
	}
	if i < len(files)-1 {
		nextLink = fmt.Sprintf("%s.html", filepath.Base(files[i+1]))
	}

	v := imageContainerVars{
		PrevLink:  prevLink,
		CurImage:  filepath.Base(src),
		NextLink:  nextLink,
		SourceURL: fmt.Sprintf("https://github.com/glacials/twos.dev/blob/main/%s", src),
	}

	if err := t.Execute(f, v); err != nil {
		return fmt.Errorf("can't execute imgcontainer template: %w", err)
	}

	return nil
}
