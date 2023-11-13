package winter

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/yargevad/filepathx"
)

// Substructure is a graph of documents on the website.
type Substructure struct {
	cfg Config
	// devURL is [twos.dev/winter.Substructure.cfg.Development.URL] unmarshaled into a [net/url.URL].
	devURL *url.URL
	docs   documents
	// photos is a map of gallery name to slice of photos in that gallery.
	photos map[string][]*img
}

var (
	galName = regexp.MustCompile(
		`img\/.+\/(.+)/.*`,
	)

	// ignorePaths are file path regex patterns that should not be
	// treated as documents to be generated, even though they are in source
	// directories. The regex is a full match so must start with ^ and end with $.
	//
	// TODO: Migrate some of these to a directory blocklist instead.
	ignorePaths = map[*regexp.Regexp]struct{}{
		regexp.MustCompile(`(^\.)|(^.*\/\.)\.[^\/]*$`):         {},
		regexp.MustCompile(`^(.*\/)?#.*$`):                     {},
		regexp.MustCompile(`^(.*\/)?README.md$`):               {},
		regexp.MustCompile(`^(.*\/)?.DS_Store$`):               {},
		regexp.MustCompile(`^(.*\/)?imgcontainer.html.tmpl$`):  {},
		regexp.MustCompile(`^(.*\/)?text_document.html.tmpl$`): {},
		regexp.MustCompile(`^(.*\/)?_gallery.html.tmpl$`):      {},
		regexp.MustCompile(`^(.*\/)?_icon.html.tmpl$`):         {},
		regexp.MustCompile(`^(.*\/)?_imgs.html.tmpl$`):         {},
		regexp.MustCompile(`^(.*\/)?_nav.html.tmpl$`):          {},
		regexp.MustCompile(`^(.*\/)?_subtoc.html.tmpl$`):       {},
		regexp.MustCompile(`^(.*\/)?_toc.html.tmpl$`):          {},
		regexp.MustCompile(`^(.*\/)?_videos.html.tmpl$`):       {},
	}
	pad = newPadder()
)

// NewSubstructure returns a substructure with the given configuration.
// Upon initialization, a substructure is the result of a discovery phase of content on the filesystem.
// Further calls are needed to build the full graph of content and render it to HTML.
func NewSubstructure(cfg Config) (*Substructure, error) {
	devURL, err := url.Parse(cfg.Development.URL)
	if err != nil {
		return nil, fmt.Errorf("cannot parse config development.url %q: %w", cfg.Development.URL, err)
	}
	s := Substructure{
		cfg:    cfg,
		devURL: devURL,
	}
	return &s, s.discover()
}

// add adds the given document to the substructure,
// removing any old versions in the process.
func (s *Substructure) add(d Document) {
	// dedupe
	for i, doc := range s.docs {
		if doc.Metadata().SourcePath == d.Metadata().SourcePath {
			s.docs[i] = d
			return
		}
	}
	s.docs = append(s.docs, d)
}

func (s *Substructure) addIMG(im *img) error {
	if s.photos == nil {
		s.photos = map[string][]*img{}
	}
	match := galName.FindStringSubmatch(im.WebPath)
	if len(match) < 1 {
		return fmt.Errorf("cannot find a gallery name in photo path %q", im.WebPath)
	}
	name := match[1]
	s.photos[name] = append(s.photos[name], im)
	return nil
}

// discover clears the substructure of any known documents and
// discovers all documents from scratch on the filesystem.
func (s *Substructure) discover() error {
	paths := append(s.cfg.Src, "src")
	for _, path := range paths {
		if err := s.discoverAtPath(path); err != nil {
			return err
		}
	}
	if err := s.discoverStatic("public"); err != nil {
		return err
	}
	return nil
}

// discoverAtPath discovers all documents in or at the given path
// glob.
func (s *Substructure) discoverAtPath(path string) error {
	if err := s.discoverHTML(path); err != nil {
		return err
	}
	if err := s.discoverMarkdown(path); err != nil {
		return err
	}
	if err := s.discoverOrg(path); err != nil {
		return err
	}
	if err := s.discoverPhotos(path); err != nil {
		return err
	}
	if err := s.discoverTemplates(path); err != nil {
		return err
	}

	sort.Sort(s.docs)
	return nil
}

// galleryGlobs are the relative path components with which to discover gallery images.
// Each is appended to the path supplied to discoverGalleriesAtPath and used to perform a glob.
//
// The glob supports double asterisks, which mean "any character, including a path separator".
// Otherwise, syntax is identical to that of [filepath.Glob].
var galleryGlobs []string = []string{
	"img/**/*.[jJ][pP][gG]",
	"img/**/*.[jJ][pP][eE][gG]",
}

// discoverHTML discovers all HTML documents in or at the given path glob.
// The files are wrapped in textDocument types and then added to the
// substructure.
func (s *Substructure) discoverHTML(path string) error {
	var htmlFiles []string
	if stat, err := os.Stat(path); err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	} else if stat.IsDir() {
		htmlFiles, err = filepathx.Glob(filepath.Join(path, "**", "*.html"))
		if err != nil {
			return err
		}
	} else if strings.HasSuffix(path, ".html") {
		htmlFiles = append(htmlFiles, path)
	}

	for _, src := range htmlFiles {
		if shouldIgnore(src) {
			continue
		}
		s.add(NewHTMLDocument(src))
	}
	for _, d := range s.docs {
		if p, _, ok := strings.Cut(d.Metadata().Filename, "_"); ok {
			parent, ok := s.DocBySourcePath(p)
			if ok {
				d.Metadata().Parent = parent.Metadata().SourcePath
			}
		}
	}

	return nil
}

// discoverMarkdown discovers all Markdown documents in or at the given
// path glob. The files are wrapped in textDocument types and then added to the
// substructure.
func (s *Substructure) discoverMarkdown(path string) error {
	var markdownFiles []string
	if stat, err := os.Stat(path); err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	} else if stat.IsDir() {
		markdownFiles, err = filepathx.Glob(filepath.Join(path, "**", "*.md"))
		if err != nil {
			return err
		}
	} else if strings.ToLower(filepath.Ext(path)) == ".md" {
		markdownFiles = append(markdownFiles, path)
	}

	for _, src := range markdownFiles {
		if shouldIgnore(src) {
			continue
		}
		s.add(NewMarkdownDocument(src))
	}

	return nil
}

// discoverOrg discovers all Org documents in or at the given path glob.
// The files are wrapped in textDocument types and then added to the
// substructure.
func (s *Substructure) discoverOrg(path string) error {
	// TODO: Allow looking in user's org directory.
	// TODO: Allow rendering only a subsection of an org file.
	var orgFiles []string
	if stat, err := os.Stat(path); err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	} else if stat.IsDir() {
		orgFiles, err = filepathx.Glob(filepath.Join(path, "**", "*.org"))
		if err != nil {
			return err
		}
	} else {
		orgFiles = append(orgFiles, path)
	}

	for _, src := range orgFiles {
		if shouldIgnore(src) {
			continue
		}
		s.add(NewOrgDocument(src))
	}

	return nil
}

// discoverPhotos discovers all photos in or at the given path glob.
func (s *Substructure) discoverPhotos(src string) error {
	stat, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("cannot discovery galleries at %q: %w", src, err)
	}
	if !stat.IsDir() {
		return fmt.Errorf("discoverGalleriesAtPath expects a directory, but got a file: %q", src)
	}

	var files []string
	for _, g := range galleryGlobs {
		f, err := filepathx.Glob(filepath.Join(src, g))
		if err != nil {
			return err
		}
		files = append(files, f...)
	}
	sort.Sort(sort.Reverse((sort.StringSlice(files))))
	for _, src := range files {
		if shouldIgnore(src) {
			return nil
		}
		im, err := NewIMG(src, s.cfg)
		if err != nil {
			return fmt.Errorf("cannot create gallery document from %s: %w", src, err)
		}
		if err := s.addIMG(im); err != nil {
			return err
		}
	}

	return nil
}

// discoverStatic discovers all static files in or at the given path glob.
// The files are wrapped with staticDocument types then added to the
// substructure.
func (s *Substructure) discoverStatic(path string) error {
	var staticFiles []string
	if stat, err := os.Stat(path); err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	} else if stat.IsDir() {
		staticFiles, err = filepathx.Glob(filepath.Join(path, "**", "*"))
		if err != nil {
			return err
		}
	} else {
		staticFiles = append(staticFiles, path)
	}

	for _, src := range staticFiles {
		if shouldIgnore(src) {
			continue
		}
		if stat, err := os.Stat(src); err != nil {
			return err
		} else if stat.IsDir() {
			continue
		}
		s.add(NewStaticDocument(src))
	}

	return nil
}

func (s *Substructure) discoverTemplates(path string) error {
	var tmplFiles []string
	if stat, err := os.Stat(path); err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	} else if stat.IsDir() {
		tmplFiles, err = filepathx.Glob(filepath.Join(path, "**", "*.tmpl"))
		if err != nil {
			return err
		}
	} else if strings.HasSuffix(path, ".html") {
		tmplFiles = append(tmplFiles, path)
	}

	for _, src := range tmplFiles {
		if shouldIgnore(src) {
			continue
		}
		s.add(NewTemplateDocument(src, s.docs))
	}

	return nil
}

type ErrNotTracked struct{ path string }

func (err ErrNotTracked) Error() string {
	return fmt.Sprintf(
		"%s is not tracked; restart to track",
		err.path,
	)
}

// Rebuild rebuilds the document or template at the given path into the given
// dist directory.
//
// If src is a template, any documents that use it will be rebuilt afterwards.
//
// If src is a document, any templates it uses will be rebuilt first.
//
// If src is a post, the index, writing, and archives pages will be rebuilt after.
//
// If src isn't known to the substructure, Rebuild returns ErrNotTracked.
func (s *Substructure) Rebuild(src, dist string) error {
	var built bool

	// Try first as-is; if that fails we'll discover() and try again.
	for i := 0; i < 2; i++ {
		for _, doc := range s.docs {
			for d := range doc.Dependencies() {
				if d == src {
					fmt.Printf("  ↗ %s", pad(s.devURL.JoinPath(doc.Metadata().Filename).String()))
					if err := s.execute(doc, dist); err != nil {
						return fmt.Errorf(
							"can't rebuild %s upstream dependency %s: %w",
							src,
							d,
							err,
						)
					}
					fmt.Println(" ✓")
					built = true
				}
			}
			if doc.Metadata().SourcePath == src {
				fmt.Printf("  → %s", pad(s.devURL.JoinPath(doc.Metadata().Filename).String()))
				if err := s.execute(doc, dist); err != nil {
					return fmt.Errorf("can't rebuild changed file %s: %w", src, err)
				}
				fmt.Println(" ✓")
				built = true
			}
		}
		if d, ok := s.DocBySourcePath(src); ok && d.Metadata().Kind == post {
			fmt.Printf("  ↘ %s", pad(s.devURL.JoinPath("archives.html").String()))
			archives, ok := s.DocBySourcePath("src/cold/archives.html.tmpl")
			if !ok {
				return fmt.Errorf("src/cold/archives.html.tmpl not found")
			}
			if err := s.execute(archives, dist); err != nil {
				return fmt.Errorf("can't rebuild index: %w", err)
			}
			fmt.Println(" ✓")
			fmt.Printf("  ↘ %s", pad(s.devURL.JoinPath("writing.html").String()))
			writing, ok := s.DocBySourcePath("src/cold/writing.html")
			if !ok {
				return fmt.Errorf("src/cold/writing.html not found")
			}
			if err := s.execute(writing, dist); err != nil {
				return fmt.Errorf("can't rebuild index: %w", err)
			}
			fmt.Println(" ✓")
			fmt.Printf("  ↘ %s", pad(s.devURL.JoinPath("index.html").String()))
			index, ok := s.DocBySourcePath("src/cold/index.md")
			if !ok {
				return fmt.Errorf("src/cold/index.md not found")
			}
			if err := s.execute(index, dist); err != nil {
				return fmt.Errorf("can't rebuild index: %w", err)
			}
			fmt.Println(" ✓")
			built = true
		}

		if !built {
			if i == 0 {
				if err := s.discover(); err != nil {
					return err
				}
			} else {
				return ErrNotTracked{src}
			}
		}
	}

	return nil
}

// DocBySourcePath returns the document that originated from the file at path.
//
// If no such document exists, ok is false.
func (s *Substructure) DocBySourcePath(path string) (doc Document, ok bool) {
	for _, doc := range s.docs {
		if filepath.Clean(doc.Metadata().SourcePath) == filepath.Clean(path) {
			return doc, true
		}
	}
	return nil, false
}

// ExecuteAll builds all documents known to the substructure,
// as well as any site-scoped non-documents such as RSS feeds.
func (s *Substructure) ExecuteAll(dist string) error {
	built := map[string]Document{}
	for _, d := range s.docs {
		if prev, ok := built[d.Metadata().Filename]; ok {
			return fmt.Errorf(
				"both %s and %s wanted to build to %s/%s; remove one",
				d.Metadata().SourcePath,
				prev.Metadata().SourcePath,
				s.cfg.Hostname,
				d.Metadata().Filename,
			)
		}
		if err := s.execute(d, dist); err != nil {
			return fmt.Errorf(
				"cannot execute %s while executing all: %w",
				d.Metadata().SourcePath,
				err,
			)
		}
		built[d.Metadata().Filename] = d
	}

	if err := s.writefeed(); err != nil {
		return fmt.Errorf("cannot generate feed: %w", err)
	}

	return s.validateURIsDidNotChange(dist)
}

// execute builds the given document into the given directory.
// To also build its dependencies and dependants, use Rebuild instead.
func (s *Substructure) execute(d Document, dist string) error {
	dest := filepath.Join(dist, d.Metadata().Filename)
	if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
		return fmt.Errorf(
			"can't create %q directory: %w",
			filepath.Dir(dest),
			err,
		)
	}
	f, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("can't create %s: %w", dest, err)
	}
	defer f.Close()
	return d.Render(f)
}

// shouldIgnore returns true if the given file path should not be built into the substructure,
// or false otherwise.
func shouldIgnore(src string) bool {
	for r := range ignorePaths {
		if r.MatchString(src) {
			return true
		}
	}
	return false
}
