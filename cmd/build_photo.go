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
	"html/template"
	"image"
	"image/jpeg"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/mknote"
	"golang.org/x/exp/slices"
	"golang.org/x/image/draw"
)

func photoBuilder(src, dst string) error {
	relsrc, err := filepath.Rel("src", src)
	if err != nil {
		return fmt.Errorf("can't get relpath for photo `%s`: %w", src, err)
	}

	dst = filepath.Join(dst, relsrc)
	thmdst := strings.Replace(
		dst,
		filepath.FromSlash("/img/"),
		filepath.FromSlash("/img/thumb/"),
		1,
	)

	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return fmt.Errorf("can't mkdir `%s`: %w", filepath.Dir(dst), err)
	}

	srcf, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("can't read `%s`: %w", src, err)
	}
	defer srcf.Close()

	dstf, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("can't write photo `%s`: %w", dst, err)
	}
	defer dstf.Close()

	if _, err := io.Copy(dstf, srcf); err != nil {
		return fmt.Errorf("can't copy `%s` to `%s`: %w", src, dst, err)
	}

	// Generate thumbnails and place into build dir
	if err := genThumbnail(src, thmdst, 300); err != nil {
		return fmt.Errorf("can't generate thumbnails: %w", err)
	}

	if err := genGalleryPage(src, fmt.Sprintf("%s.html", dst)); err != nil {
		return fmt.Errorf("can't generate photo container pages: %w", err)
	}

	return nil
}

// genThumbnail makes thumbnails from every photo (recursively) in src, copying
// them to equivalent paths in a "thumb" subdirectory of dst.
//
// The thumbnails have the given width. Height is automatically calculated to
// maintain ratio.
func genThumbnail(src, dst string, width int) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("can't open photo at path `%s`: %w", src, err)
	}
	defer sourceFile.Close()

	srcPhoto, err := jpeg.Decode(sourceFile)
	if err != nil {
		return fmt.Errorf(
			"can't decode photo at path `%s` (maybe not an image?): %w",
			src,
			err,
		)
	}
	p := srcPhoto.Bounds().Size()
	w, h := width, ((width*p.X/p.Y)+1)&-1
	dstPhoto := image.NewRGBA(image.Rect(0, 0, w, h))

	draw.CatmullRom.Scale(
		dstPhoto,
		image.Rectangle{image.Point{0, 0}, image.Point{width, width}},
		srcPhoto,
		image.Rectangle{image.Point{0, 0}, srcPhoto.Bounds().Size()},
		draw.Over,
		nil,
	)

	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return fmt.Errorf(
			"can't make thumbnail directory in path `%s`: %w",
			dst,
			err,
		)
	}

	destinationFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf(
			"can't create thumbnail file for photo at path `%s`: %w",
			src,
			err,
		)
	}
	defer destinationFile.Close()

	if err := jpeg.Encode(destinationFile, dstPhoto, nil); err != nil {
		return fmt.Errorf(
			"can't encode to destination file at path `%s`: %w",
			dst,
			err,
		)
	}

	return nil
}

func genGalleryPage(src, dst string) error {
	templateHTML, err := ioutil.ReadFile(photoGalleryTemplatePath)
	if err != nil {
		return fmt.Errorf("can't read photo gallery template: %w", err)
	}

	t, err := template.New("imgcontainer").Parse(string(templateHTML))
	if err != nil {
		return fmt.Errorf("can't create imgcontainer template: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return fmt.Errorf(
			"can't create imgcontainer directory `%s`: %w",
			filepath.Dir(dst),
			err,
		)
	}

	f, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("can't create imgcontainer file for `%s`: %w", src, err)
	}
	defer f.Close()

	files, err := filepath.Glob(
		filepath.Join(filepath.Dir(src), "*.[jJ][pP][gG]"),
	)
	if err != nil {
		return fmt.Errorf(
			"can't look into photo directory `%s` for ordering: %w",
			filepath.Dir(src),
			err,
		)
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

	camera, err := camera(src)
	if err != nil {
		return fmt.Errorf("can't get camera for `%s`: %w", src, err)
	}

	v := galleryPageVars{
		Alt:    "testing",
		Camera: camera,
		URL:    filepath.Base(src),

		Prev: prevLink,
		Next: nextLink,

		pageVars: pageVars{
			SourceURL: fmt.Sprintf(
				"https://github.com/glacials/twos.dev/blob/main/%s",
				src,
			),
		},
	}

	if err := t.Execute(f, v); err != nil {
		return fmt.Errorf("can't execute imgcontainer template: %w", err)
	}

	return nil
}

// camera extracts the camera string (including lens, etc.) from the image at
// the given path.
func camera(src string) (string, error) {
	f, err := os.Open(src)
	if err != nil {
		return "", fmt.Errorf("can't open photo: %w", err)
	}
	defer f.Close()

	exif.RegisterParsers(mknote.All...)

	x, err := exif.Decode(f)
	if err != nil {
		return "", fmt.Errorf("can't read exif data: %w", err)
	}

	focalLength, err := exifFractionToDecimal(x, exif.FocalLength)
	if err != nil {
		return "", fmt.Errorf("can't get focal length: %w", err)
	}

	camModel, err := x.Get(exif.Model) // normally, don't ignore errors!
	if err != nil {
		return "", fmt.Errorf("can't get camera model: %w", err)
	}

	cam, err := camModel.StringVal()
	if err != nil {
		return "", fmt.Errorf("can't render camera as string: %w", err)
	}

	fnum, err := exifFractionToDecimal(x, exif.FNumber)
	if err != nil {
		return "", fmt.Errorf("can't get focal length: %w", err)
	}

	exposure, err := x.Get(exif.ExposureTime)
	if err != nil {
		return "", fmt.Errorf("can't get exposure: %w", err)
	}

	iso, err := x.Get(exif.ISOSpeedRatings)
	if err != nil {
		return "", fmt.Errorf("can't get ISO: %w", err)
	}

	return fmt.Sprintf(
		"%s • %.0fmm • ƒ%.1f • %ss • ISO %s",
		cam,
		focalLength,
		fnum,
		strings.Replace(exposure.String(), "\"", "", 2),
		iso,
	), nil
}

func exifFractionToDecimal(
	x *exif.Exif,
	field exif.FieldName,
) (float64, error) {
	fraction, err := x.Get(field)
	if err != nil {
		return 0, fmt.Errorf("can't get field %s: %w", field, err)
	}
	parts := strings.Split(
		strings.Replace(fraction.String(), "\"", "", 2),
		"/",
	)
	numer, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, fmt.Errorf(
			"can't convert %s (numerator of %s, %s) to int: %w",
			parts[0],
			field,
			fraction,
			err,
		)
	}
	denom, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, fmt.Errorf(
			"can't convert %s (denominator of %s, %s) to int: %w",
			parts[0],
			field,
			fraction,
			err,
		)
	}

	return float64(numer) / float64(denom), nil
}
