package winter

import (
	"bytes"
	"fmt"
	"html/template"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/gorilla/feeds"
	"github.com/yargevad/filepathx"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// Substructure is a graph of documents on the website, generated with read-only
// operations. It will later be fed into a renderer.
type Substructure struct {
	cfg Config
	// devURL is [twos.dev/winter.Substructure.cfg.Development.URL] unmarshaled into a [net/url.URL].
	devURL *url.URL
	docs   documents
	// photos is a map of gallery name to slice of photos in that gallery.
	photos map[string][]*galleryDocument
}

// substructureDocument is a Document wrapper allowing the Substructure to manage its own data.
type substructureDocument struct {
	Document

	// Parent, if non-nil, points to the document that this one is a child of.
	// A child document can be used to provide a navigational hierarchy to documents.
	Parent *substructureDocument
	// Source is the path to the source file for the document,
	// relative to the project root.
	Source string
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

func (d *substructureDocument) Shortname() string {
	dest, err := d.Dest()
	if err != nil {
		return ""
	}
	filename := filepath.Base(dest)
	shortname, _, _ := strings.Cut(filename, ".")
	return shortname
}

// NewSubstructure returns a substructure with the given configuration. Upon
// initialization, a substructure is the result of a discovery phase of content
// on the filesystem. Further calls are needed to build the full graph of
// content and render it to HTML.
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

// add adds the given document to the substructure, removing any old
// versions in the process.
func (s *Substructure) add(d *substructureDocument) {
	// dedupe
	for i, doc := range s.docs {
		if doc.Source == d.Source {
			s.docs[i] = d
			return
		}
	}

	if gal, ok := d.Document.(*galleryDocument); ok {
		if s.photos == nil {
			s.photos = map[string][]*galleryDocument{}
		}
		match := galName.FindStringSubmatch(gal.WebPath)
		if len(match) < 1 {
			panic(fmt.Sprintf("cannot find a gallery name in photo path %s", gal.WebPath))
		}
		name := match[1]
		s.photos[name] = append(s.photos[name], gal)
	}

	s.docs = append(s.docs, d)
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
	if err := s.discoverStaticAtPath("public"); err != nil {
		return err
	}
	return nil
}

// discoverAtPath discovers all documents in or at the given path
// glob.
func (s *Substructure) discoverAtPath(path string) error {
	if err := s.discoverHTMLAtPath(path); err != nil {
		return err
	}

	if err := s.discoverMarkdownAtPath(path); err != nil {
		return err
	}

	if err := s.discoverOrgAtPath(path); err != nil {
		return err
	}

	if err := s.discoverGalleriesAtPath(path); err != nil {
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

// discoverGalleriesAtPath discovers all galleries in or at the given path glob.
// The files are wrapped in galleryDocument types and then added to the
// substructure.
func (s *Substructure) discoverGalleriesAtPath(src string) error {
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
	var prev *galleryDocument
	for _, src := range files {
		if shouldIgnore(src) {
			return nil
		}
		doc, err := NewGalleryDocument(src, s.cfg)
		if err != nil {
			return fmt.Errorf("cannot create gallery document from %s: %w", src, err)
		}
		// Set the last doc's next to this doc.
		if prev != nil {
			prev.Next = doc
		}
		// Set this doc's previous to the last doc.
		// Update prev so the next doc can do the same.
		doc.Prev, prev = prev, doc
		s.add(&substructureDocument{Document: doc, Source: doc.Source})
	}

	return nil
}

// discoverHTMLAtPath discovers all HTML documents in or at the given path glob.
// The files are wrapped in textDocument types and then added to the
// substructure.
func (s *Substructure) discoverHTMLAtPath(path string) error {
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

	for _, src := range append(htmlFiles, tmplFiles...) {
		if shouldIgnore(src) {
			continue
		}
		doc, err := NewHTMLDocument(src)
		if err != nil {
			return fmt.Errorf("cannot create HTML document from %s: %w", src, err)
		}
		s.add(&substructureDocument{Document: doc, Source: src})
	}
	for _, d := range s.docs {
		if p, _, ok := strings.Cut(d.Shortname(), "_"); ok {
			d.Parent = s.DocByShortname(p)
		}
	}

	return nil
}

// discoverMarkdownAtPath discovers all Markdown documents in or at the given
// path glob. The files are wrapped in textDocument types and then added to the
// substructure.
func (s *Substructure) discoverMarkdownAtPath(path string) error {
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
		doc, err := NewMarkdownDocument(src)
		if err != nil {
			return fmt.Errorf("cannot create Markdown document from %s: %w", src, err)
		}
		s.add(&substructureDocument{Document: doc, Source: src})
	}

	return nil
}

// discoverOrgAtPath discovers all Org documents in or at the given path glob.
// The files are wrapped in textDocument types and then added to the
// substructure.
func (s *Substructure) discoverOrgAtPath(path string) error {
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
		doc, err := NewOrgDocument(src)
		if err != nil {
			return err
		}
		s.add(&substructureDocument{Document: doc, Source: src})
	}

	return nil
}

// discoverStaticAtPath discovers all static files in or at the given path glob.
// The files are wrapped with staticDocument types then added to the
// substructure.
func (s *Substructure) discoverStaticAtPath(path string) error {
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
		doc, err := NewStaticDocument(src)
		if err != nil {
			return err
		}
		s.add(&substructureDocument{Document: doc, Source: src})
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
// If the path is a template, any documents that use it will be rebuilt
// afterwards. If the path is a document, any templates it uses will be rebuilt
// first. If the path isn't known to the substructure, Rebuild returns
// ErrNotTracked.
//
// As a special case, the index page will be rebuilt after any post is built, so
// that it can display its up to date title.
func (s *Substructure) Rebuild(src, dist string) error {
	var built bool

	// Try first as-is; if that fails we'll discover() and try again.
	for i := 0; i < 2; i++ {
		for _, doc := range s.docs {
			dest, err := doc.Dest()
			if err != nil {
				return err
			}
			for d := range doc.Dependencies() {
				if d == src {
					fmt.Printf("  ↗ %s", pad(s.devURL.JoinPath(dest).String()))
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
			if doc.Source == src {
				fmt.Printf("  → %s", pad(s.devURL.JoinPath(dest).String()))
				if err := s.execute(doc, dist); err != nil {
					return fmt.Errorf("can't rebuild changed file %s: %w", src, err)
				}
				fmt.Println(" ✓")
				built = true
			}
		}
		if d := s.DocBySrc(src); d != nil && d.IsPost() {
			fmt.Printf("  ↘ %s", pad(s.devURL.JoinPath("archives.html").String()))
			if err := s.execute(s.DocByShortname("archives"), dist); err != nil {
				return fmt.Errorf("can't rebuild index: %w", err)
			}
			fmt.Println(" ✓")
			fmt.Printf("  ↘ %s", pad(s.devURL.JoinPath("writing.html").String()))
			if err := s.execute(s.DocByShortname("writing"), dist); err != nil {
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

// DocByShortname returns the document with the given shortname, or nil if none
// exists. The returned struct implements Document.
func (s *Substructure) DocByShortname(shortname string) *substructureDocument {
	for _, d := range s.docs {
		dest, err := d.Dest()
		if err != nil {
			continue
		}
		base := filepath.Base(dest)
		base, _, _ = strings.Cut(base, ".")
		if base == shortname {
			return d
		}
	}
	return nil
}

// DocBySrc returns the document with the given source file path, or nil if none
// exists. The returned struct implements Document.
func (s *Substructure) DocBySrc(path string) *substructureDocument {
	for _, d := range s.docs {
		if d.Source == path {
			return d
		}
	}
	return nil
}

// archives returns posts grouped by year.
func (s *Substructure) archives() (archivesVars, error) {
	m := map[int]documents{}
	for _, d := range s.posts() {
		year := d.CreatedAt().Year()
		m[year] = append(m[year], d)
	}

	var archives archivesVars
	for year, docs := range m {
		sort.Sort(docs)
		archives = append(archives, archiveVars{
			Year:      year,
			Documents: docs,
		})
	}
	sort.Sort(archives)

	return archives, nil
}

// categories returns posts grouped by category. The key of the returned map is
// the non-pluralized category name in title casing.
//
// The empty string is a valid category and will appear in the map.
func (s *Substructure) categories() map[string]documents {
	cats := map[string]documents{}
	for _, d := range s.posts() {
		cat := cases.Title(language.English).String(d.Category())
		cats[cat] = append(cats[cat], d)
	}
	return cats
}

// gallery returns the gallery pages to be shown, grouped by gallery name.
func (s *Substructure) gallery() map[string][]*galleryDocument {
	return s.photos
}

// posts returns text documents of type post.
func (s *Substructure) posts() (docs documents) {
	for _, d := range s.docs {
		if d.IsPost() {
			docs = append(docs, d)
		}
	}
	return
}

func (s *Substructure) writefeed() error {
	now := time.Now()
	feed := feeds.Feed{
		Title:       s.cfg.Name,
		Description: s.cfg.Description,
		Author: &feeds.Author{
			Name:  s.cfg.Author.Name,
			Email: s.cfg.Author.Email,
		},
		Link: &feeds.Link{Href: (&url.URL{Scheme: "https", Host: s.cfg.Hostname}).String()},
		Copyright: fmt.Sprintf(
			"Copyright %d–%d %s",
			s.cfg.Since,
			now.Year(),
			s.cfg.Author.Name,
		),
		Items: []*feeds.Item{},

		Created: now,
		Updated: now,
	}

	for _, post := range s.posts() {
		body, err := post.Build()
		if err != nil {
			return err
		}
		dest, err := post.Dest()
		if err != nil {
			return err
		}
		t := template.New(post.Source)
		funcs, err := s.funcmap(post)
		if err != nil {
			return fmt.Errorf("can't generate funcmap: %w", err)
		}
		_ = t.Funcs(funcs)
		_, err = t.Parse(string(body))
		if err != nil {
			return fmt.Errorf("cannot parse page for feed: %w", err)
		}
		if err := loadDeps(t); err != nil {
			return fmt.Errorf("can't load dependency templates for %q: %w", t.Name(), err)
		}
		var buf bytes.Buffer
		if err = post.Execute(&buf, t); err != nil {
			return fmt.Errorf("cannot execute document %q: %w", post.Source, err)
		}

		feed.Items = append(feed.Items, &feeds.Item{
			Id:          dest,
			Title:       post.Title(),
			Author:      feed.Author,
			Content:     string(buf.String()),
			Description: string(body),
			Link: &feeds.Link{
				Href: fmt.Sprintf("%s/%s.html", feed.Link.Href, dest),
			},
			Created: post.CreatedAt(),
			Updated: post.UpdatedAt(),
		})
	}

	atom, err := feed.ToAtom()
	if err != nil {
		return err
	}
	if err := os.WriteFile("dist/feed.atom", []byte(atom), 0o644); err != nil {
		return err
	}

	rss, err := feed.ToRss()
	if err != nil {
		return err
	}
	if err := os.WriteFile("dist/feed.rss", []byte(rss), 0o644); err != nil {
		return err
	}

	return nil
}

// ExecuteAll builds all documents known to the substructure,
// as well as any site-scoped non-documents such as RSS feeds.
func (s *Substructure) ExecuteAll(dist string, cfg Config) error {
	built := map[string]*substructureDocument{}
	for _, d := range s.docs {
		dest, err := d.Dest()
		if err != nil {
			return err
		}
		if prev, ok := built[dest]; ok {
			return fmt.Errorf(
				"both %s and %s wanted to build to %s/%s; remove one",
				d.Source,
				prev.Source,
				cfg.Hostname,
				dest,
			)
		}
		if err := s.execute(d, dist); err != nil {
			return fmt.Errorf(
				"cannot execute %s while executing all: %w",
				d.Source,
				err,
			)
		}
		built[dest] = d
	}

	return s.writefeed()
}

// execute builds the given document into the given directory. To build a
// document and its dependencies and dependants, use Rebuild instead.
func (s *Substructure) execute(d *substructureDocument, dist string) error {
	dest, err := d.Dest()
	if err != nil {
		return err
	}
	dest = filepath.Join(dist, dest)

	if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
		return fmt.Errorf(
			"can't create %s directory `%s`: %w",
			galTmplPath,
			filepath.Dir(dest),
			err,
		)
	}

	bodyBytes, err := d.Build()
	if err != nil {
		return fmt.Errorf("cannot build document %q: %w", d.Shortname(), err)
	}
	if d.Layout() == "" {
		return os.WriteFile(dest, bodyBytes, 0o644)
	}

	layoutBytes, err := os.ReadFile(d.Layout())
	if err != nil {
		return fmt.Errorf("cannot read %q to execute %q: %w", d.Layout(), d.Source, err)
	}
	funcs, err := s.funcmap(d)
	if err != nil {
		return fmt.Errorf("can't generate funcmap: %w", err)
	}
	t, err := template.New(d.Layout()).Funcs(funcs).Parse(string(layoutBytes))
	if err != nil {
		return fmt.Errorf("can't parse `%s`: %w", d.Layout(), err)
	}
	if _, err := t.New("body").Parse(string(bodyBytes)); err != nil {
		return fmt.Errorf("can't parse %s: %w", d.Shortname(), err)
	}
	if err := loadDeps(t); err != nil {
		return fmt.Errorf("can't load dependencies: %s", err)
	}

	// Hardcode that every document with a layout depends on our CSS.
	deps := d.Dependencies() // TODO: Consider a d.AddDependency method
	deps[filepath.Join("public", "style.css")] = struct{}{}

	for _, depTmpl := range t.Templates() {
		if depTmpl.Name() != "body" && depTmpl.Name() != d.Source {
			deps[depTmpl.Name()] = struct{}{}
		}
	}

	f, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("can't create %s: %w", dest, err)
	}
	defer f.Close()

	if err = d.Execute(f, t); err != nil {
		return fmt.Errorf("can't execute document `%s`: %w", d.Shortname(), err)
	}

	return nil
}

func (s *Substructure) funcmap(d *substructureDocument) (template.FuncMap, error) {
	iconFunc, err := d.icon()
	if err != nil {
		return nil, err
	}
	now := time.Now()

	return template.FuncMap{
		"add": add,
		"div": div,
		"mul": mul,
		"sub": sub,

		"now":    func() time.Time { return now },
		"parent": func() *substructureDocument { return d.Parent },
		"src":    func() string { return d.Source },

		"archives":   s.archives,
		"categories": s.categories,
		"gallery":    s.gallery,
		"icon":       iconFunc,
		"posts":      s.posts,
	}, nil
}

// shouldIgnore returns true if the given file path should not be built into the
// substructure, or false otherwise.
func shouldIgnore(src string) bool {
	for r := range ignorePaths {
		if r.MatchString(src) {
			return true
		}
	}
	return false
}
