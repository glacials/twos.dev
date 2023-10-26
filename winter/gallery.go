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
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/mknote"
	"golang.org/x/image/draw"
)

var extensionRegex = regexp.MustCompile(`(.*)\.([^\.]*)`)

type galleryDocument struct {
	EXIF

	PageLink   string
	Alt        string
	WebPath    string
	Thumbnails []*thumbnail

	Prev *galleryDocument
	Next *galleryDocument

	cfg       Config
	localPath string
}

type thumbnail struct {
	WebPath string
	Width   int
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
		WebPath:   fmt.Sprintf("/%s", relpath),
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

	if err := os.MkdirAll(filepath.Dir(imgdest), 0o755); err != nil {
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

	if err := d.genThumbnails(d.localPath, thmdest); err != nil {
		return nil, fmt.Errorf("can't generate thumbnails: %w", err)
	}

	d.EXIF, err = photodata(d.localPath)
	if err != nil {
		return nil, fmt.Errorf(
			"can't get camera for `%s`: %w",
			d.localPath,
			err,
		)
	}

	b, err := tmplByName(tmplPathToName(galTmplPath))
	if err != nil {
		return nil, err
	}

	for old, new := range replacements {
		re, err := regexp.Compile(old)
		if err != nil {
			return nil, err
		}
		b = re.ReplaceAll(b, []byte(new))
	}

	return b, nil
}

func (d *galleryDocument) Category() string { return "" }

// Dependencies returns the filepaths the gallery document depends on.
func (d *galleryDocument) Dependencies() map[string]struct{} {
	return map[string]struct{}{tmplPathToName(galTmplPath): {}}
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

// Preview returns the empty string.
func (d *galleryDocument) Preview() string { return "" }

// Title returns a generic gallery title.
func (d *galleryDocument) Title() string { return fmt.Sprintf("%s Photo Viewer", d.cfg.Name) }

// CreatedAt returns the date and time the photo was taken, or a zero time if
// unknown.
func (d *galleryDocument) CreatedAt() time.Time { return d.EXIF.TakenAt }

// UpdatedAt returns the date and time the photo was last updated, or a zero
// time if unknown.
func (d *galleryDocument) UpdatedAt() time.Time { return time.Time{} }

type EXIF struct {
	Camera       string
	FocalLength  float64
	Aperture     float64
	ShutterSpeed string
	ISO          string
	TakenAt      time.Time
}

// photodata extracts the EXIF string (including lens, etc.) and timestamp from
// the image at the given path.
func photodata(src string) (EXIF, error) {
	f, err := os.Open(src)
	if err != nil {
		return EXIF{}, fmt.Errorf("can't open photo: %w", err)
	}
	defer f.Close()

	exif.RegisterParsers(mknote.All...)

	x, err := exif.Decode(f)
	if err != nil {
		return EXIF{}, fmt.Errorf("can't read exif data: %w", err)
	}

	focalLength, err := exifFractionToDecimal(x, exif.FocalLength)
	if err != nil {
		return EXIF{}, fmt.Errorf("can't get focal length: %w", err)
	}

	camModel, err := x.Get(exif.Model) // normally, don't ignore errors!
	if err != nil {
		return EXIF{}, fmt.Errorf("can't get camera model: %w", err)
	}

	cam, err := camModel.StringVal()
	if err != nil {
		return EXIF{}, fmt.Errorf("can't render camera as string: %w", err)
	}

	fnum, err := exifFractionToDecimal(x, exif.FNumber)
	if err != nil {
		return EXIF{}, fmt.Errorf("can't get focal length: %w", err)
	}

	exposure, err := x.Get(exif.ExposureTime)
	if err != nil {
		return EXIF{}, fmt.Errorf("can't get exposure: %w", err)
	}

	iso, err := x.Get(exif.ISOSpeedRatings)
	if err != nil {
		return EXIF{}, fmt.Errorf("can't get ISO: %w", err)
	}

	timestamp, err := x.DateTime()
	if err != nil {
		if errors.Is(err, exif.TagNotPresentError("")) {
			return EXIF{}, fmt.Errorf(
				"photo has no EXIF timestamp... what do?",
			)
		}
		return EXIF{}, fmt.Errorf("can't get photo datetime: %w", err)
	}

	_, err = x.Get(exif.GPSInfoIFDPointer)
	if err == nil {
		// location data is set! no no no!
		panic(fmt.Sprintf("photo %s has location data! please strip it.", src))
	}

	return EXIF{
		cam,
		focalLength,
		fnum,
		strings.Replace(exposure.String(), "\"", "", 2),
		iso.String(),
		timestamp,
	}, nil
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

// genThumbnails makes thumbnails of several sizes from the photo located at src,
// placing them in dst with dimensions added to the filename, just before the extension.
//
// It generates thumbnails with widths of powers of 2,
// from 1 until the largest width possible that is still smaller than the source image.
// Heights are automatically calculated to mantain aspect ratio.
//
// For example, a 500x500 image would have thumbnails of
// 1x1, 2x2, 4x4, 8x8, 16x16, 32x32, 64x64, 128x128, and 256x256
// generated.
func (d *galleryDocument) genThumbnails(src, dst string) error {
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

	for width := 1; width < p.X; width *= 2 {
		height := (width * p.X / p.Y) & -1
		dstPhoto := image.NewRGBA(image.Rect(0, 0, width, height))

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

		if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
			return fmt.Errorf(
				"can't make thumbnail directory `%s`: %w",
				filepath.Dir(dst),
				err,
			)
		}

		srcFilename := filepath.Base(src)
		matches := extensionRegex.FindStringSubmatch(srcFilename)
		if len(matches) < 3 {
			return fmt.Errorf("cannot add dimensions to filename %q with no extension", srcFilename)
		}
		dstFilename := fmt.Sprintf("%s.%dx%d.%s", matches[1], width, height, matches[2])
		dstPath := filepath.Join(filepath.Dir(dst), dstFilename)
		destinationFile, err := os.Create(dstPath)
		if err != nil {
			return fmt.Errorf(
				"can't create thumbnail file for photo at path `%s`: %w",
				src,
				err,
			)
		}
		defer destinationFile.Close()

		/*
			palette := []color.Color{
				color.RGBA{0x00, 0x00, 0x00, 0xff},
				color.RGBA{0x00, 0x00, 0xff, 0xff},
				//		color.RGBA{0x00, 0xff, 0x00, 0xff},
				//		color.RGBA{0x00, 0xff, 0xff, 0xff},
				color.RGBA{0xff, 0x00, 0x00, 0xff},
				color.RGBA{0xff, 0x00, 0xff, 0xff},
				//		color.RGBA{0xff, 0xff, 0x00, 0xff},
				//		color.RGBA{0x00, 0xff, 0xff, 0xff},
				color.RGBA{0xff, 0xff, 0xff, 0xff},
			}
			d := dither.NewDitherer(palette)
			d.Matrix = dither.FloydSteinberg

			// Dither tries to modify the original first, and may return nil if it does,
			// so guard against that.
			new := d.Dither(dstPhoto)
			if new != nil {
				dstPhoto = nil
			}
		*/

		if err := jpeg.Encode(destinationFile, dstPhoto, nil); err != nil {
			return fmt.Errorf(
				"can't encode to destination file at path `%s`: %w",
				dst,
				err,
			)
		}
		webPath, err := filepath.Rel("dist", dstPath)
		if err != nil {
			return fmt.Errorf("cannot get relative path for thumbnail %q: %w", dst, err)
		}
		webPath = fmt.Sprintf("/%s", webPath)
		d.Thumbnails = append(d.Thumbnails, &thumbnail{
			WebPath: webPath,
			Width:   width,
		})
	}

	return nil
}
