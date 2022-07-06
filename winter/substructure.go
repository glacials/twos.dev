package winter

import (
	"bytes"
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
	ignoreFiles = map[string]struct{}{
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

// NewSubstructure returns a substructure with the given configuration. Upon
// initialization, a substructure is the result of a discovery phase of content
// on the filesystem. Further calls are needed to build the full graph of
// content and render it to HTML.
func NewSubstructure(cfg Config) (*Substructure, error) {
	s := Substructure{cfg: cfg}
	md, err := filepathx.Glob("src/**/*.md")
	if err != nil {
		return nil, err
	}

	for _, src := range md {
		if _, ok := ignoreFiles[filepath.Base(src)]; ok {
			continue
		}
		doc, err := NewMarkdownDocument(src)
		if err != nil {
			return nil, err
		}
		s.docs = append(s.docs, doc)

	}

	coldhtml, err := filepath.Glob("src/cold/*.html")
	if err != nil {
		return nil, err
	}

	warmhtml, err := filepath.Glob("src/warm/*.html")
	if err != nil {
		return nil, err
	}

	coldtmpl, err := filepath.Glob("src/cold/*.tmpl")
	if err != nil {
		return nil, err
	}

	warmtmpl, err := filepath.Glob("src/warm/*.tmpl")
	if err != nil {
		return nil, err
	}

	html := append(coldhtml, warmhtml...)
	tmpl := append(coldtmpl, warmtmpl...)

	for _, src := range append(html, tmpl...) {
		if _, ok := ignoreFiles[filepath.Base(src)]; ok {
			continue
		}
		doc, err := NewHTMLDocument(src)
		if err != nil {
			return nil, err
		}
		s.docs = append(s.docs, doc)
	}

	sort.Sort(s.docs)

	return &s, err
}

func (s *Substructure) DocByShortname(shortname string) *document {
	for _, d := range s.docs {
		if d.Shortname == shortname {
			return d
		}
	}
	return nil
}

func (s *Substructure) DocBySrc(path string) *document {
	for _, d := range s.docs {
		if d.SrcPath == path {
			return d
		}
	}
	return nil
}

func (s *Substructure) posts() (docs documents) {
	for _, d := range s.docs {
		if d.Kind == post {
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
		body, err := post.build()
		if err != nil {
			return err
		}
		feed.Items = append(feed.Items, &feeds.Item{
			Id:          post.Shortname,
			Title:       post.Title,
			Author:      feed.Author,
			Content:     string(body),
			Description: string(body),
			Link: &feeds.Link{
				Href: fmt.Sprintf("%s/%s.html", feed.Link.Href, post.Shortname),
			},
			Created: post.CreatedAt,
			Updated: post.UpdatedAt,
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

func (s *Substructure) setlinks() error {
	for _, out := range s.docs {
		linksout, err := out.linksout()
		if err != nil {
			return err
		}

		for _, l := range linksout {
			l = strings.TrimSuffix(l, filepath.Ext(l))
			if in := s.DocByShortname(l); in != nil {
				out.outgoing = append(out.outgoing, in)
				in.incoming = append(in.incoming, out)
			}
		}
	}

	return nil
}

// ExecuteAll builds all documents known to the substructure, as well as any
// site-scoped non-documents such as RSS feeds.
func (s *Substructure) ExecuteAll(dist string) error {
	for _, d := range s.docs {
		if err := s.Execute(d, dist); err != nil {
			return err
		}
	}

	s.writefeed()

	return nil
}

// Execute builds the given document into the given directory.
func (s *Substructure) Execute(d *document, dist string) error {
	t := template.New("")
	if err := loadTemplates(t); err != nil {
		return err
	}

	imgsFunc, err := imgs(d.Shortname)
	if err != nil {
		return err
	}

	videoFunc, err := videos(d.Shortname)
	if err != nil {
		return err
	}

	postsFunc := s.posts

	_ = t.Funcs(template.FuncMap{
		"img":    imgsFunc,
		"imgs":   imgsFunc,
		"video":  videoFunc,
		"videos": videoFunc,
		"posts":  postsFunc,
	})

	b, err := d.build()
	if err != nil {
		return err
	}

	if _, err := t.New("body").Parse(string(b)); err != nil {
		return fmt.Errorf("can't parse %s: %w", d.SrcPath, err)
	}

	var buf bytes.Buffer
	err = t.Lookup("text_document").
		Execute(&buf, templateVars{d, s, time.Now()})
	if err != nil {
		return fmt.Errorf("can't execute document `%s`: %w", d.Shortname, err)
	}

	path := filepath.Join(dist, d.Shortname+".html")
	if os.WriteFile(path, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("can't write document `%s`: %w", d.Shortname, err)
	}

	return nil
}
