package winter

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/gorilla/feeds"
	"github.com/yargevad/filepathx"
)

var (
	ignoreFilenames = map[string]struct{}{
		"README.md": {},
		".DS_Store": {},
	}
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

// discover clears the substructure of any known documents and discovers all
// documents from scratch from the filesystem.
func (s *Substructure) discover() error {
	md, err := filepathx.Glob("src/**/*.md")
	if err != nil {
		return err
	}

	for _, src := range md {
		if _, ok := ignoreFilenames[filepath.Base(src)]; ok {
			continue
		}
		doc, err := NewMarkdownDocument(src)
		if err != nil {
			return err
		}
		s.docs = append(s.docs, &substructureDocument{Document: doc, Source: src})
	}

	warmhtml, err := filepathx.Glob("src/warm/*.html")
	if err != nil {
		return err
	}

	warmtmpl, err := filepathx.Glob("src/warm/*.tmpl")
	if err != nil {
		return err
	}

	coldhtml, err := filepathx.Glob("src/cold/*.html")
	if err != nil {
		return err
	}

	coldtmpl, err := filepathx.Glob("src/cold/*.tmpl")
	if err != nil {
		return err
	}

	warm := append(warmhtml, warmtmpl...)
	cold := append(coldhtml, coldtmpl...)

	for _, src := range append(warm, cold...) {
		if _, ok := ignoreFilenames[filepath.Base(src)]; ok {
			continue
		}
		doc, err := NewHTMLDocument(src)
		if err != nil {
			return err
		}
		s.docs = append(s.docs, &substructureDocument{Document: doc, Source: src})
	}
	for _, d := range s.docs {
		if p, _, ok := strings.Cut(d.Shortname(), "_"); ok {
			d.Parent = s.DocByShortname(p)
		}
	}

	jpg, err := filepathx.Glob("src/**/*.[jJ][pP][gG]")
	if err != nil {
		return err
	}

	jpeg, err := filepathx.Glob("src/**/*.[jJ][pP][eE][gG]")
	if err != nil {
		return err
	}

	var prev, next *galleryDocument
	for _, src := range append(jpg, jpeg...) {
		if _, ok := ignoreFilenames[filepath.Base(src)]; ok {
			continue
		}
		doc, err := NewGalleryDocument(src, s.cfg)
		if err != nil {
			return err
		}
		doc.Prev, prev = prev, doc
		s.docs = append(s.docs, &substructureDocument{Document: doc, Source: src})
	}
	for d := prev; d != nil; d = d.Prev {
		d.Next, next = next, d
	}

	static, err := filepathx.Glob("public/**/*")
	if err != nil {
		return err
	}

	for _, src := range static {
		if _, ok := ignoreFilenames[filepath.Base(src)]; ok {
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
		s.docs = append(s.docs, &substructureDocument{Document: doc, Source: src})
	}

	sort.Sort(s.docs)

	return err
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
	pad := pad()
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
		Description: s.cfg.Desc,
		Author: &feeds.Author{
			Name:  s.cfg.AuthorName,
			Email: s.cfg.AuthorEmail,
		},
		Link: &feeds.Link{Href: s.cfg.Domain.String()},
		Copyright: fmt.Sprintf(
			"Copyright %d–%d %s",
			s.cfg.Since,
			now.Year(),
			s.cfg.AuthorName,
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
	if err := ioutil.WriteFile("dist/feed.atom", []byte(atom), 0644); err != nil {
		return err
	}

	rss, err := feed.ToRss()
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile("dist/feed.rss", []byte(rss), 0644); err != nil {
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
				cfg.Domain.Host,
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

	if err := os.MkdirAll(filepath.Dir(dest), 0755); err != nil {
		return fmt.Errorf(
			"can't create %s directory `%s`: %w",
			galname,
			filepath.Dir(dest),
			err,
		)
	}

	bodyBytes, err := d.Build()
	if err != nil {
		return fmt.Errorf("can't build %s: %w", d.Shortname(), err)
	}
	if d.Layout() == "" {
		return os.WriteFile(dest, bodyBytes, 0644)
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

		"archives": s.archives,
		"img":      imgsFunc,
		"imgs":     imgsFunc,
		"posts":    s.posts,
		"video":    videoFunc,
		"videos":   videoFunc,
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

	// Hardcode that every document with a layout depends on style.css
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
