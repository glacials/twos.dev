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
	"golang.org/x/image/draw"
)

type galleryDocument struct {
	PageLink string
	Alt      string
	Camera   string
	TakenAt  time.Time
	WebPath  string

	Prev *galleryDocument
	Next *galleryDocument

	cfg       Config
	localPath string
}

// NewGalleryDocument returns a Document that represents an HTML page which
// wraps a single, large-display image.
func NewGalleryDocument(src string, cfg Config) (*galleryDocument, error) {
	relpath, err := filepath.Rel("src", src)
	if err != nil {
		return nil, fmt.Errorf("can't get relpath for photo `%s`: %w", src, err)
	}

	return &galleryDocument{
		localPath: src,
		PageLink:  fmt.Sprintf("/%s.html", relpath),
		WebPath:   relpath,
		cfg:       cfg,
	}, nil
}

// Build builds the gallery document.
func (d *galleryDocument) Build() ([]byte, error) {
	imgdest := filepath.Join("dist", d.WebPath) // TODO: Do from substructure
	thmdest := strings.Replace(
		imgdest,
		filepath.FromSlash("/img/"),
		filepath.FromSlash("/img/thumb/"),
		1,
	)

	if err := os.MkdirAll(filepath.Dir(imgdest), 0755); err != nil {
		return nil, fmt.Errorf("can't mkdir `%s`: %w", filepath.Dir(imgdest), err)
	}

	srcf, err := os.Open(d.localPath)
	if err != nil {
		return nil, fmt.Errorf("can't read `%s`: %w", d.localPath, err)
	}
	defer srcf.Close()

	dstf, err := os.Create(imgdest)
	if err != nil {
		return nil, fmt.Errorf("can't write photo `%s`: %w", imgdest, err)
	}
	defer dstf.Close()

	if _, err := io.Copy(dstf, srcf); err != nil {
		return nil, fmt.Errorf(
			"can't copy `%s` to `%s`: %w",
			d.localPath,
			imgdest,
			err,
		)
	}

	if err := genThumbnail(d.localPath, thmdest, 300); err != nil {
		return nil, fmt.Errorf("can't generate thumbnails: %w", err)
	}

	d.Camera, d.TakenAt, err = photodata(d.localPath)
	if err != nil {
		return nil, fmt.Errorf(
			"can't get camera for `%s`: %w",
			d.localPath,
			err,
		)
	}

	return tmplByName(galname)
}

// Dependencies returns the filepaths the gallery document depends on.
func (d *galleryDocument) Dependencies() map[string]struct{} {
	return map[string]struct {
	}{
		galname: {},
	}
}

// Dest returns the destination path for the final gallery document.
func (d *galleryDocument) Dest() (string, error) {
	relpath, err := filepath.Rel("src", d.localPath)
	if err != nil {
		return "", fmt.Errorf(
			"can't get relpath for photo `%s`: %w",
			d.localPath,
			err,
		)
	}

	return fmt.Sprintf("%s.html", relpath), nil
}

func (d *galleryDocument) Execute(w io.Writer, t *template.Template) error {
	return t.Execute(w, d)
}

func (d *galleryDocument) IsPost() bool  { return false }
func (d *galleryDocument) IsDraft() bool { return false }

// LoadTemplates loads the needed templates for the gallery document.
func (d *galleryDocument) LoadTemplates(t *template.Template) error {
	return nil
}

// Layout returns the image container template name.
func (d *galleryDocument) Layout() string { return "imgcontainer" }

// Title returns a generic gallery title.
func (d *galleryDocument) Title() string { return fmt.Sprintf("%s Photo Viewer", d.cfg.Name) }

// CreatedAt returns the date and time the photo was taken, or a zero time if
// unknown.
func (d *galleryDocument) CreatedAt() time.Time { return d.TakenAt }

// UpdatedAt returns the date and time the photo was last updated, or a zero
// time if unknown.
func (d *galleryDocument) UpdatedAt() time.Time { return time.Time{} }

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
			"can't make thumbnail directory `%s`: %w",
			filepath.Dir(dst),
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
