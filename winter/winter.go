// Package winter supplies a static website generator.
// It has these goals:
//
//  1. New content must be easy to edit. Old content must be hard to break.
//  2. Be replaceable. Special syntax should look fine as plaintext if moved elsewhere.
//
// Winter is part of the [TadiWeb].
//
// # Easy to edit, hard to break
//
// To accomplish its first goal, Winter employs two mechanisms.
//
// First, content managed by Winter is either warm or cold.
// Warm content is welcome to be synchronized into a Winter-managed directory by external tools,
// such as your mobile text editor or note-taking app of choice.
// These external tools and their synchronization steps are not provided by Winter,
// and they are not required to use it,
// but Winter's design goals assume they exist.
// Warm content is content you are actively working on.
// It is not ready to be listed anywhere,
// but it is published to a page on your website so you can view it in context and share it with friends for review.
//
// Cold content is finished content. [Cool URIs don't change.]
// When your warm content is ready to be published,
// you should freeze it:
//
//	winter freeze src/cold/my_post.md
//
// This converts it into cold content.
// The most important difference between the two is that your synchronization tools never, ever touch cold content.
// Sychronization tools, jobs, and scripts are run frequently and probably aren't your most robust pipelines.
// To limit the blast radius of any bugs they experience,
// content that you don't expect to edit often must be moved away from the directories they control.
// This is cold content.
//
// Cold content is also guaranteed by Winter's quality checks.
// When a piece of content becomes cold, Winter commits its URL to memory and ensures that URL stays accessible in future runs.
// If ever a URL that once was accessible becomes inaccessible, Winter alerts you and prevents that build from succeeding.
// This protects cold content from you, your tools, and from Winter's current and future versions.
//
// # No lock-in
//
// To accomplish its second goal, Winter provides cherries on top of Markdown.
// However, using the same Markdown with any other static website generator results in no degradation.
// Just like Markdown itself, Winter Markdown looks great when not parsed.
//
// # Galleries
//
// A block containing only images is treated as a gallery,
// with the images placed in a responsive grid.
// Images can be clicked to zoom in or out.
//
//	![An image of a cat.](/img/cat.jpg)
//	![An image of a dog.](/img/dog.jpg)
//
// # Image captions
//
// A block of all-italic text immediately below an image or gallery is treated as a caption,
// and given a special visual treatment and accessibility structure.
//
//	![An image of a cat.](/img/cat.jpg)
//	![An image of a dog.](/img/dog.jpg)
//
//	_My cat and dog like to play with each other._
//
// # Photos as first-class citizens
//
// For photographers, Winter supports photographs as first-class citizens.
// EXIF data is automatically extracted and can be displayed using template variables.
// Photos organized into galleries display next to each other neatly,
// with a built-in lightbox.
//
// # Photo EXIF safety
//
// Any photos Winter processes that have GPS data embedded in them are loudly rejected,
// failing the build.
//
// # Reduced load times
//
// Images and references to them are automatically converted into WebP format and several thumbnails are generated for each.
// Each image renders with [<img srcset>] to ensure only the smallest possible image that saturates the display density is loaded.
//
// # LaTeX
//
// Surround text in $dollar signs$ to render LaTeX,
// implemented via [KaTeX].
//
//	$\LaTeX$ users rejoice!
//
// # Tables of contents
//
// A table of contents can be requested by setting the toc variable to true in frontmatter.
// Tables of contents are always rendered immediately above the first level-2 heading.
//
//	---
//	toc: true
//	---
//
//	# My Article
//
//	(table of contents will be rendered here)
//
//	## My Cat
//
//	...
//
//	## My Dog
//
//	...
//
// # Syntax highlighting
//
// Fenced code blocks that specify a language are syntax-highlighted.
//
//	```go
//	package main
//	func main() {
//	  fmt.Println("I'm syntax highlighted!")
//	}
//	```
//
// # External links in new tabs
//
// Any links that navigate to external websites will automatically have a target=_blank set during generation.
//
// # Dark mode images
//
// If an image's extensionless filename ends in "-dark" or "-light",
// and another image exists at the same path but with the opposite suffix,
// the correct one will be rendered to the user based on their light-/dark-mode preference.
//
// [<img srcset>]: https://developer.mozilla.org/en-US/docs/Learn/HTML/Multimedia_and_embedding/Responsive_images
// [Cool URIs don't change.]: https://www.w3.org/Provider/Style/URI
// [KaTeX]: https://katex.org
// [TadiWeb]: https://www.tadiweb.com
package winter

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/yargevad/filepathx"
)

const (
	AppName = "winter"
)

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
		//		regexp.MustCompile(`^(.*\/)?_gallery.html.tmpl$`):      {},
		//		regexp.MustCompile(`^(.*\/)?_icon.html.tmpl$`):         {},
		//		regexp.MustCompile(`^(.*\/)?_imgs.html.tmpl$`):         {},
		//		regexp.MustCompile(`^(.*\/)?_nav.html.tmpl$`):          {},
		//		regexp.MustCompile(`^(.*\/)?_subtoc.html.tmpl$`):       {},
		//		regexp.MustCompile(`^(.*\/)?_toc.html.tmpl$`):          {},
		//		regexp.MustCompile(`^(.*\/)?_videos.html.tmpl$`):       {},
	}
	// ignoreDirs are file path directories from which no documents should be generated.
	ignoreDirs = map[string]struct{}{
		"src/templates/": {},
	}
	pad = newPadder()
)

type ErrNotTracked struct{ path string }

func (err ErrNotTracked) Error() string {
	return fmt.Sprintf(
		"%s is not tracked; restart to track",
		err.path,
	)
}

// Substructure is a graph of documents on the website.
type Substructure struct {
	// cfg holds user preferences, specified by winter.yml.
	cfg Config
	// devURL is [twos.dev/winter.Substructure.cfg.Development.URL] unmarshaled into a [net/url.URL].
	devURL *url.URL
	// docs holds the documents known to the substructure.
	docs documents
	// galleries is a map of gallery name to slice of galleries in that gallery.
	galleries map[string][]*img
}

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

// Build builds doc.
//
// To also build downstream dependencies, use [Rebuild] instead.
func (s *Substructure) Build(doc Document) error {
	r, err := os.Open(doc.Metadata().SourcePath)
	if err != nil {
		return fmt.Errorf("cannot read %q for building %q: %w", doc.Metadata().SourcePath, doc.Metadata().Title, err)
	}
	defer r.Close()
	if err := doc.Load(r); err != nil {
		return fmt.Errorf("cannot load %q for building %q: %w", doc.Metadata().SourcePath, doc.Metadata().Title, err)
	}
	dest := filepath.Join(s.cfg.Dist, doc.Metadata().WebPath)
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
		return fmt.Errorf("cannot render %q for building: %w", doc.Metadata().SourcePath, err)
	}
	fmt.Println(" ✓")
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
	builtIMGs := map[string]*img{}
	for _, gallery := range s.galleries {
		for _, img := range gallery {
			if prev, ok := builtIMGs[img.WebPath]; ok {
				return fmt.Errorf(
					"both %s (%T) and %q (%T) wanted to build to %q/%q; remove one",
					img.SourcePath,
					img,
					prev.SourcePath,
					prev,
					s.cfg.Production.URL,
					img.WebPath,
				)
			}
			srcf, err := os.Open(img.SourcePath)
			if err != nil {
				return fmt.Errorf("cannot open image: %w", err)
			}
			defer srcf.Close()
			dest := filepath.Join(dist, img.WebPath)
			if err := img.Load(srcf); err != nil {
				return fmt.Errorf("cannot load image: %w", err)
			}
			fresh, err := img.generatedPhotosAreFresh(img.SourcePath)
			if err != nil {
				return fmt.Errorf("cannot check freshness of %q: %w", img.SourcePath, err)
			}
			if fresh {
				continue
			}
			if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
				return fmt.Errorf("cannot make gallery dir %q: %w", filepath.Dir(dest), err)
			}
			destf, err := os.Create(dest)
			if err != nil {
				return fmt.Errorf("cannot write image %q to %q during ExecuteAll: %w", img.SourcePath, dest, err)
			}
			defer destf.Close()
			if err := img.Render(destf); err != nil {
				return err
			}
			if err := s.Rebuild(img.SourcePath); err != nil {
				return fmt.Errorf("cannot rebuild image %q during ExecuteAll: %w", img.SourcePath, err)
			}
			builtIMGs[img.WebPath] = img
		}
	}
	builtDocs := map[string]Document{}
	for _, doc := range s.docs {
		if prev, ok := builtDocs[doc.Metadata().WebPath]; ok {
			return fmt.Errorf(
				"both %s (%T) and %s (%T) wanted to build to %s/%s; remove one",
				doc.Metadata().SourcePath,
				doc,
				prev.Metadata().SourcePath,
				prev,
				s.cfg.Production.URL,
				doc.Metadata().WebPath,
			)
		}
		if err := s.Rebuild(doc.Metadata().SourcePath); err != nil {
			return fmt.Errorf(
				"cannot execute %q during ExecuteAll: %w",
				doc.Metadata().SourcePath,
				err,
			)
		}
		builtDocs[doc.Metadata().WebPath] = doc
	}

	if err := s.writefeed(); err != nil {
		return fmt.Errorf("cannot generate feed: %w", err)
	}

	return s.validateURIsDidNotChange(dist)
}

// Rebuild rebuilds the document or template at the given path.
// Then, it rebuilds any downstream dependencies.
//
// If src isn't known to the substructure, Rebuild no-ops and returns no error.
func (s *Substructure) Rebuild(src string) error {
	fmt.Printf("%s ↓\n", src)
	if doc, ok := s.DocBySourcePath(src); ok {
		if err := s.Build(doc); err != nil {
			return fmt.Errorf("cannot retrieve doc at %q: %w", src, err)
		}
	}
	for _, doc := range s.docs {
		if doc.DependsOn(src) {
			if err := s.Build(doc); err != nil {
				return fmt.Errorf("cannot build %q (dependent of %q): %w", doc.Metadata().SourcePath, src, err)
			}
		}
	}
	return nil
}

// add adds the given document to the substructure,
// removing any old versions in the process.
func (s *Substructure) add(doc Document) {
	// dedupe
	for i, d := range s.docs {
		if d.Metadata().SourcePath == doc.Metadata().SourcePath {
			s.docs[i] = doc
			// Sort again in case d's creation date changed.
			sort.Sort(s.docs)
			return
		}
	}
	s.docs = append(s.docs, doc)
	sort.Sort(s.docs)
}

// addIMG adds the given image to the substructure,
// removing any old versions in the process.
func (s *Substructure) addIMG(im *img) error {
	if s.galleries == nil {
		s.galleries = map[string][]*img{}
	}
	match := galName.FindStringSubmatch(im.WebPath)
	if len(match) < 1 {
		return fmt.Errorf("cannot find a gallery name in photo path %q", im.WebPath)
	}
	name := match[1]
	s.galleries[name] = append(s.galleries[name], im)
	return nil
}

// discover clears the substructure of any known documents and discovers all documents from scratch on the filesystem.
func (s *Substructure) discover() error {
	log.Println("Starting discovery.")
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
	if err := s.discoverGalleries(path); err != nil {
		return err
	}
	if err := s.discoverHTML(path); err != nil {
		return err
	}
	if err := s.discoverMarkdown(path); err != nil {
		return err
	}
	if err := s.discoverOrg(path); err != nil {
		return err
	}
	// TODO: Calling discoverTemplates once misses *.html.tmpl files.
	// Calling it twice doesn't. Fix?
	if err := s.discoverTemplates(path); err != nil {
		return err
	}
	if err := s.discoverTemplates(path); err != nil {
		return err
	}
	sort.Sort(s.docs)
	return nil
}

// discoverGalleries adds all documents matching galleryGlobs in or at the given path glob to the substructure.
func (s *Substructure) discoverGalleries(src string) error {
	log.Println("Discovering galleries.")
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

// discoverHTML adds all *.html documents in or at the given path glob to the substructure.
func (s *Substructure) discoverHTML(path string) error {
	log.Println("Discovering HTML.")
	var htmlFiles []string
	if stat, err := os.Stat(path); err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	} else if stat.IsDir() {
		htmlf, err := filepathx.Glob(filepath.Join(path, "**", "*.html"))
		if err != nil {
			return err
		}
		htmlFiles = append(htmlFiles, htmlf...)
	} else if strings.HasSuffix(path, ".html") {
		htmlFiles = append(htmlFiles, path)
	}

	for _, src := range htmlFiles {
		if shouldIgnore(src) {
			continue
		}
		meta := NewMetadata(src, tmplPath)
		s.add(
			NewHTMLDocument(src, meta,
				NewTemplateDocument(src, meta, s.docs, s.galleries, nil),
			),
		)
	}
	for _, d := range s.docs {
		if p, _, ok := strings.Cut(d.Metadata().WebPath, "_"); ok {
			parent, ok := s.DocBySourcePath(p)
			if ok {
				d.Metadata().ParentFilename = parent.Metadata().WebPath
			}
		}
	}

	return nil
}

// discoverMarkdown adds all *.md documents in or at the given path glob to the substructure.
func (s *Substructure) discoverMarkdown(path string) error {
	log.Println("Discovering markdown.")
	var mdFiles []string
	if stat, err := os.Stat(path); err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	} else if stat.IsDir() {
		mdFiles, err = filepathx.Glob(filepath.Join(path, "**", "*.md"))
		if err != nil {
			return err
		}
	} else {
		mdFiles = append(mdFiles, path)
	}

	for _, src := range mdFiles {
		if shouldIgnore(src) {
			continue
		}
		meta := NewMetadata(src, tmplPath)
		s.add(
			NewMarkdownDocument(src, meta,
				NewHTMLDocument(src, meta,
					NewTemplateDocument(src, meta, s.docs, s.galleries, nil),
				),
			),
		)
	}

	return nil
}

// discoverOrg adds all *.org documents in or at the given path glob to the substructure.
func (s *Substructure) discoverOrg(path string) error {
	log.Println("Discovering Org.")
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
		meta := NewMetadata(src, tmplPath)
		s.add(
			NewOrgDocument(src, meta,
				NewHTMLDocument(src, meta,
					NewTemplateDocument(src, meta, s.docs, s.galleries, nil),
				),
			),
		)
	}

	return nil
}

// galleryGlobs are the relative path components with which to discover images.
// Each is appended to the path supplied to discoverPhotos and used to perform a glob.
//
// The glob supports double asterisks, which mean "any character, including a path separator".
// Otherwise, syntax is identical to that of [filepath.Glob].
var galleryGlobs []string = []string{
	"img/**/*.[jJ][pP][gG]",
	"img/**/*.[jJ][pP][eE][gG]",
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

// discoverTemplates adds all *.html.tmpl documents in or at the given path glob to the substructure.
func (s *Substructure) discoverTemplates(path string) error {
	log.Println("Discovering templates.")
	var tmplFiles []string
	if stat, err := os.Stat(path); err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	} else if stat.IsDir() {
		tmplf, err := filepathx.Glob(filepath.Join(path, "**", "*.html.tmpl"))
		if err != nil {
			return err
		}
		tmplFiles = append(tmplFiles, tmplf...)
	} else if strings.HasSuffix(path, ".html.tmpl") {
		tmplFiles = append(tmplFiles, path)
	}

	for _, src := range tmplFiles {
		if shouldIgnore(src) {
			continue
		}
		meta := NewMetadata(src, tmplPath)
		s.add(
			NewHTMLDocument(src, meta,
				NewTemplateDocument(src, meta, s.docs, s.galleries, nil),
			),
		)
	}
	for _, d := range s.docs {
		if p, _, ok := strings.Cut(d.Metadata().WebPath, "_"); ok {
			parent, ok := s.DocBySourcePath(p)
			if ok {
				d.Metadata().ParentFilename = parent.Metadata().WebPath
			}
		}
	}

	return nil
}

// shouldIgnore returns true if the given file path should not be built into the substructure,
// or false otherwise.
func shouldIgnore(src string) bool {
	for dir := range ignoreDirs {
		if strings.HasPrefix(filepath.Clean(src), filepath.Clean(dir)) {
			return true
		}
	}
	for r := range ignorePaths {
		if r.MatchString(src) {
			return true
		}
	}
	return false
}
