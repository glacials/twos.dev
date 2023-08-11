package winter

import (
	"fmt"
	"html/template"
	"io/ioutil"
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
	cfg  Config
	docs documents
}

// substructureDocument is a wrapper around Document so Substructure can store
// some extra fields on it.
type substructureDocument struct {
	Document

	Parent *substructureDocument
	// Source is the path to the source file for the document.
	Source string
}

var (
	// ignorePaths are file path regex patterns that should not be
	// treated as documents to be generated, even though they are in source
	// directories. The regex is a full match so must start with ^ and end with $.
	//
	// TODO: Migrate some of these to a directory blocklist instead.
	ignorePaths = map[*regexp.Regexp]struct{}{
		regexp.MustCompile("(^\\.)|(^.*\\/\\.)\\.[^\\/]*$"):     {},
		regexp.MustCompile("^(.*\\/)?#.*$"):                     {},
		regexp.MustCompile("^(.*\\/)?README.md$"):               {},
		regexp.MustCompile("^(.*\\/)?.DS_Store$"):               {},
		regexp.MustCompile("^(.*\\/)?imgcontainer.html.tmpl$"):  {},
		regexp.MustCompile("^(.*\\/)?text_document.html.tmpl$"): {},
		regexp.MustCompile("^(.*\\/)?_icon.html.tmpl$"):         {},
		regexp.MustCompile("^(.*\\/)?_imgs.html.tmpl$"):         {},
		regexp.MustCompile("^(.*\\/)?_nav.html.tmpl$"):          {},
		regexp.MustCompile("^(.*\\/)?_subtoc.html.tmpl$"):       {},
		regexp.MustCompile("^(.*\\/)?_toc.html.tmpl$"):          {},
		regexp.MustCompile("^(.*\\/)?_videos.html.tmpl$"):       {},
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
	s := Substructure{cfg: cfg}
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

// discoverGalleriesAtPath discovers all galleries in or at the given path glob.
// The files are wrapped in galleryDocument types and then added to the
// substructure.
func (s *Substructure) discoverGalleriesAtPath(path string) error {
	var jpgFiles []string
	if stat, err := os.Stat(path); err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	} else if stat.IsDir() {
		jpgFiles, err = filepathx.Glob(filepath.Join(path, "**", "*.[jJ][pP][gG]"))
		if err != nil {
			return err
		}
	} else if strings.ToLower(filepath.Ext(path)) == ".jpg" {
		jpgFiles = append(jpgFiles, path)
	}

	var jpegFiles []string
	if stat, err := os.Stat(path); err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	} else if stat.IsDir() {
		jpegFiles, err = filepathx.Glob(filepath.Join(path, "**", "*.[jJ][pP][eE][gG]"))
		if err != nil {
			return err
		}
	} else if strings.ToLower(filepath.Ext(path)) == ".jpeg" {
		jpegFiles = append(jpegFiles, path)
	}

	var prev, next *galleryDocument
	for _, src := range append(jpgFiles, jpegFiles...) {
		if shouldIgnore(src) {
			continue
		}
		doc, err := NewGalleryDocument(src, s.cfg)
		if err != nil {
			return fmt.Errorf("cannot create gallery document from %s: %w", src, err)
		}
		doc.Prev, prev = prev, doc
		s.add(&substructureDocument{Document: doc, Source: src})
	}
	for d := prev; d != nil; d = d.Prev {
		d.Next, next = next, d
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
			name := tmplPathToName(src)
			for d := range doc.Dependencies() {
				if d == name || d == src {
					fmt.Printf("  ↗ %s", pad(dest))
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
				fmt.Printf("  → %s", pad(dest))
				if err := s.execute(doc, dist); err != nil {
					return fmt.Errorf("can't rebuild changed file %s: %w", src, err)
				}
				fmt.Println(" ✓")
				built = true
			}
		}
		if d := s.DocBySrc(src); d != nil && d.IsPost() {
			fmt.Printf("  ↘ %s", pad("index.html"))
			if err := s.execute(s.DocByShortname("index"), dist); err != nil {
				return fmt.Errorf("can't rebuild index: %w", err)
			}
			fmt.Println(" ✓")
			fmt.Printf("  ↘ %s", pad("archives.html"))
			if err := s.execute(s.DocByShortname("archives"), dist); err != nil {
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
		base, _, _ = strings.Cut(dest, ".")
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
		feed.Items = append(feed.Items, &feeds.Item{
			Id:          dest,
			Title:       post.Title(),
			Author:      feed.Author,
			Content:     string(body),
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
	if err := ioutil.WriteFile("dist/feed.atom", []byte(atom), 0o644); err != nil {
		return err
	}

	rss, err := feed.ToRss()
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile("dist/feed.rss", []byte(rss), 0o644); err != nil {
		return err
	}

	return nil
}

// ExecuteAll builds all documents known to the substructure, as well as any
// site-scoped non-documents such as RSS feeds.
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
				"can't execute %s while executing all: %w",
				d.Source,
				err,
			)
		}
		built[dest] = d
	}

	s.writefeed()

	return nil
}

func (s *Substructure) internalize(d Document) (*substructureDocument, error) {
	for _, subdoc := range s.docs {
		fmt.Println("subdoc.Document is", subdoc.Document, "d is", d)
		if d == subdoc.Document {
			return subdoc, nil
		}
	}

	return nil, ErrNotTracked{path: d.Title()}
}

// execute builds the given document into the given directory. To build a
// document and its dependencies and dependants, use Rebuild instead.
func (s *Substructure) execute(d *substructureDocument, dist string) error {
	iconFunc, err := icon(d.Shortname())
	if err != nil {
		return err
	}
	imgsFunc, err := imgs(d.Shortname())
	if err != nil {
		return err
	}
	videoFunc, err := videos(d.Shortname())
	if err != nil {
		return err
	}

	dest, err := d.Dest()
	if err != nil {
		return err
	}
	dest = filepath.Join(dist, dest)

	if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
		return fmt.Errorf(
			"can't create %s directory `%s`: %w",
			tmplPathToName(galTmplPath),
			filepath.Dir(dest),
			err,
		)
	}

	bodyBytes, err := d.Build()
	if err != nil {
		return fmt.Errorf("can't build %s: %w", d.Shortname(), err)
	}
	if d.Layout() == "" {
		return os.WriteFile(dest, bodyBytes, 0o644)
	}

	layoutBytes, err := tmplByName(d.Layout())
	if err != nil {
		return fmt.Errorf("can't find `%s`: %w", d.Layout(), err)
	}
	t := template.New(d.Layout())
	now := time.Now()
	_ = t.Funcs(template.FuncMap{
		"add": add,
		"sub": sub,

		"now":    func() time.Time { return now },
		"parent": func() *substructureDocument { return d.Parent },
		"src":    func() string { return d.Source },

		"archives":   s.archives,
		"categories": s.categories,
		"icon":       iconFunc,
		"img":        imgsFunc,
		"imgs":       imgsFunc,
		"posts":      s.posts,
		"video":      videoFunc,
		"videos":     videoFunc,
	})
	_, err = t.Parse(string(layoutBytes))
	if err != nil {
		return fmt.Errorf("can't parse `%s`: %w", d.Layout(), err)
	}

	if _, err := t.New("body").Parse(string(bodyBytes)); err != nil {
		return fmt.Errorf("can't parse %s: %w", d.Shortname(), err)
	}

	if err := loadAllDeps(t); err != nil {
		return fmt.Errorf("can't load dependencies: %s", err)
	}

	// Hardcode that every document with a layout depends on our CSS
	deps := d.Dependencies() // TODO: Consider a d.AddDependency method
	deps[filepath.Join("public", "style.css")] = struct{}{}

	for _, dep := range t.Templates() {
		if dep.Name() != "body" && dep.Name() != d.Shortname() {
			deps[dep.Name()] = struct{}{}
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
