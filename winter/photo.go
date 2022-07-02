package winter

import (
	"errors"
	"fmt"
	"html/template"
	"image"
	"image/jpeg"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/mknote"
	"golang.org/x/exp/slices"
	"golang.org/x/image/draw"
)

type galleryVars struct {
	*Document

	ImageSRC string
	Alt      string
	Camera   string
	TakenAt  time.Time

	Prev *galleryImageVars
	Next *galleryImageVars
}

type galleryImageVars struct {
	PageLink string
	ImageSRC string
}

func BuildPhoto(src, dst string, cfg Config) error {
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

	if err := genGalleryPage(src, fmt.Sprintf("%s.html", dst), cfg); err != nil {
		return fmt.Errorf("can't generate photo container pages: %w", err)
	}

	return nil
}

func genGalleryPage(src, dst string, cfg Config) error {
	t := template.New("")

	if err := loadTemplates(t); err != nil {
		return fmt.Errorf("can't load templates: %w", err)
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

	var prev, next *galleryImageVars
	if i > 0 {
		img := filepath.Base(files[i-1])
		prev = &galleryImageVars{
			ImageSRC: img,
			PageLink: fmt.Sprintf("%s.html", img),
		}
	}
	if i < len(files)-1 {
		img := filepath.Base(files[i+1])
		next = &galleryImageVars{
			ImageSRC: img,
			PageLink: fmt.Sprintf("%s.html", img),
		}
	}

	camera, datetime, err := photodata(src)
	if err != nil {
		return fmt.Errorf("can't get camera for `%s`: %w", src, err)
	}

	shortname := filepath.Base(src)
	shortname, _, _ = strings.Cut(shortname, ".")

	v := galleryVars{
		Document: &Document{
			SourcePath:           src,
			Kind:                 gallery,
			FrontmatterShortname: shortname,
			FrontmatterTitle:     fmt.Sprintf("%s Photo Viewer", cfg.Name),
		},

		Alt:      "",
		Camera:   camera,
		ImageSRC: filepath.Base(src),
		TakenAt:  datetime,

		Prev: prev,
		Next: next,
	}

	if err := t.Lookup("imgcontainer").Execute(f, v); err != nil {
		fmt.Println(t)
		return fmt.Errorf("can't execute imgcontainer template: %w", err)
	}

	return nil
}

// photodata extracts the EXIF string (including lens, etc.) and timestamp from
// the image at the given path.
func photodata(src string) (string, time.Time, error) {
	f, err := os.Open(src)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("can't open photo: %w", err)
	}
	defer f.Close()

	exif.RegisterParsers(mknote.All...)

	x, err := exif.Decode(f)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("can't read exif data: %w", err)
	}

	focalLength, err := exifFractionToDecimal(x, exif.FocalLength)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("can't get focal length: %w", err)
	}

	camModel, err := x.Get(exif.Model) // normally, don't ignore errors!
	if err != nil {
		return "", time.Time{}, fmt.Errorf("can't get camera model: %w", err)
	}

	cam, err := camModel.StringVal()
	if err != nil {
		return "", time.Time{}, fmt.Errorf("can't render camera as string: %w", err)
	}

	fnum, err := exifFractionToDecimal(x, exif.FNumber)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("can't get focal length: %w", err)
	}

	exposure, err := x.Get(exif.ExposureTime)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("can't get exposure: %w", err)
	}

	iso, err := x.Get(exif.ISOSpeedRatings)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("can't get ISO: %w", err)
	}

	timestamp, err := x.DateTime()
	if err != nil {
		if errors.Is(err, exif.TagNotPresentError("")) {
			return "", time.Time{}, fmt.Errorf(
				"photo has no EXIF timestamp... what do?",
			)
		}
		return "", time.Time{}, fmt.Errorf("can't get photo datetime: %w", err)
	}

	_, err = x.Get(exif.GPSInfoIFDPointer)
	if err == nil {
		// location data is set! no no no!
		panic(fmt.Sprintf("photo %s has location data! please strip it.", src))
	}

	return fmt.Sprintf(
		"%s • %.0fmm • ƒ%.1f • %ss • ISO %s",
		cam,
		focalLength,
		fnum,
		strings.Replace(exposure.String(), "\"", "", 2),
		iso,
	), timestamp, nil
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
		image.Rectangle{
			image.Point{0, 0},
			image.Point{width, width},
		},
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
