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
	// cfg holds user preferences, specified by winter.yml.
	cfg Config
	// devURL is [twos.dev/winter.Substructure.cfg.Development.URL] unmarshaled into a [net/url.URL].
	devURL *url.URL
	// docs holds the documents known to the substructure.
	docs documents
	// photos is a map of gallery name to slice of photos in that gallery.
	photos map[string][]*img
}

var (
	galName = regexp.MustCompile(
		`img\/.+\/(.+)/.*`,
	)

	// ignorePaths are file path regex patterns that should not be treated as documents to be generated,
	// even though they are in source directories.
	// The regex is a full match so must start with ^ and end with $.
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

// addIMG adds the given image to the substructure,
// removing any old versions in the process.
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

// discover clears the substructure of any known documents and discovers all documents from scratch on the filesystem.
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

// discoverAtPath discovers all documents in or at the given path glob.
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

// imgGlobs are the relative path components with which to discover images.
// Each is appended to the path supplied to discoverPhotos and used to perform a glob.
//
// The glob supports double asterisks, which mean "any character, including a path separator".
// Otherwise, syntax is identical to that of [filepath.Glob].
var imgGlobs []string = []string{
	"img/**/*.[jJ][pP][gG]",
	"img/**/*.[jJ][pP][eE][gG]",
}

// discoverHTML adds all *.html documents in or at the given path glob to the substructure.
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
		if p, _, ok := strings.Cut(d.Metadata().WebPath, "_"); ok {
			parent, ok := s.DocBySourcePath(p)
			if ok {
				d.Metadata().Parent = parent.Metadata().SourcePath
			}
		}
	}

	return nil
}

// discoverMarkdown adds all *.md documents in or at the given path glob to the substructure.
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

// discoverOrg adds all *.org documents in or at the given path glob to the substructure.
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

// discoverPhotos adds all documents matching galleryGlobs in or at the given path glob to the substructure.
func (s *Substructure) discoverPhotos(src string) error {
	stat, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("cannot discovery galleries at %q: %w", src, err)
	}
	if !stat.IsDir() {
		return fmt.Errorf("discoverGalleriesAtPath expects a directory, but got a file: %q", src)
	}

	var files []string
	for _, g := range imgGlobs {
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

// discoverStatic adds all documents in or at the given path glob to the substructure.
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
		webPath, err := filepath.Rel(path, src)
		if err != nil {
			return fmt.Errorf("cannot determine desired web path of %q: %w", src, err)
		}
		s.add(NewStaticDocument(src, webPath))
	}

	return nil
}

// discoverTemplates adds all *.tmpl documents in or at the given path glob to the substructure.
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
	} else if strings.HasSuffix(path, ".tmpl") {
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
	fmt.Printf("%s ↓\n", src)
	for _, doc := range s.docs {
		if doc.Metadata().SourcePath != src && !doc.DependsOn(src) {
			continue
		}
		r, err := os.Open(doc.Metadata().SourcePath)
		if err != nil {
			return fmt.Errorf("cannot read %q for rebuilding %q: %w", doc.Metadata().SourcePath, doc.Metadata().Title, err)
		}
		if err := doc.Load(r); err != nil {
			return fmt.Errorf("cannot load %q for rebuilding %q: %w", doc.Metadata().SourcePath, doc.Metadata().Title, err)
		}
		dest := filepath.Join(dist, doc.Metadata().WebPath)
		fmt.Printf("  → %s", pad(dest))
		if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
			return fmt.Errorf("cannot make directory structure for %q: %w", dest, err)
		}
		w, err := os.Create(dest)
		if err != nil {
			return fmt.Errorf("cannot build %q (web path %q) into %q: %w", doc.Metadata().SourcePath, doc.Metadata().WebPath, dest, err)
		}
		defer w.Close()
		if err := doc.Render(w); err != nil {
			return fmt.Errorf("cannot render %q for rebuilding: %w", doc.Metadata().SourcePath, err)
		}
		fmt.Println(" ✓")
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
	builtDocs := map[string]Document{}
	for _, doc := range s.docs {
		if prev, ok := builtDocs[doc.Metadata().WebPath]; ok {
			return fmt.Errorf(
				"both %s (%T) and %s (%T) wanted to build to %s/%s; remove one",
				doc.Metadata().SourcePath,
				doc,
				prev.Metadata().SourcePath,
				prev,
				s.cfg.Hostname,
				doc.Metadata().WebPath,
			)
		}
		if err := s.Rebuild(doc.Metadata().SourcePath, dist); err != nil {
			return fmt.Errorf(
				"cannot execute %q during ExecuteAll: %w",
				doc.Metadata().SourcePath,
				err,
			)
		}
		builtDocs[doc.Metadata().WebPath] = doc
	}

	builtIMGs := map[string]*img{}
	for _, imgs := range s.photos {
		for _, img := range imgs {
			if prev, ok := builtIMGs[img.WebPath]; ok {
				return fmt.Errorf(
					"both %s (%T) and %q (%T) wanted to build to %q/%q; remove one",
					img.SourcePath,
					img,
					prev.SourcePath,
					prev,
					s.cfg.Hostname,
					img.WebPath,
				)
			}
			dest := filepath.Join(dist, img.WebPath)
			f, err := os.Create(dest)
			if err != nil {
				return fmt.Errorf("cannot write image %q to %q during ExecuteAll: %w", img.SourcePath, dest, err)
			}
			if err := img.Render(f); err != nil {
				return err
			}
			if err := s.Rebuild(img.SourcePath, dist); err != nil {
				return fmt.Errorf("cannot rebuild image %q during ExecuteAll: %w", img.SourcePath, err)
			}
			builtIMGs[img.WebPath] = img
		}
	}

	if err := s.writefeed(); err != nil {
		return fmt.Errorf("cannot generate feed: %w", err)
	}

	return s.validateURIsDidNotChange(dist)
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
