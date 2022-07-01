package winter

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/gorilla/feeds"
	"github.com/yargevad/filepathx"
)

// substructure is a graph of documents on the website, generated with read-only
// operations. It will later be fed into a renderer.
type substructure struct {
	docs []*Document
	t    *template.Template
}

// Discover is the first step of the static website build process. It crawls the
// website root looking for files that need to be wrapped up in the build
// process, such as Markdown files, templates, and static assets.
//
// It then learns what it can about each file and stores the information in
// a substructure.
func Discover(cfg Config) (s substructure, err error) {
	md, err := filepathx.Glob("**/*.md")
	if err != nil {
		return
	}

	for _, m := range md {
		doc, err := fromMarkdown(m)
		if err != nil {
			return substructure{}, err
		}
		s.docs = append(s.docs, doc)
	}

	html, err := filepathx.Glob("**/*.html")
	if err != nil {
		return
	}

	for _, h := range html {
		doc, err := fromHTML(h)
		if err != nil {
			return substructure{}, err
		}
		s.docs = append(s.docs, doc)
	}

	return
}

func (s substructure) Get(shortname string) *Document {
	for _, d := range s.docs {
		if d.Shortname() == shortname {
			return d
		}
	}
	return nil
}

func (s substructure) posts() (u []*Document) {
	for _, d := range s.docs {
		if d.meta.kind == post {
			u = append(u, d)
		}
	}
	return
}

func (s substructure) writefeed(cfg Config) error {
	now := time.Now()
	feed := feeds.Feed{
		Title:       cfg.Name,
		Description: cfg.Desc,
		Author: &feeds.Author{
			Name:  cfg.AuthorName,
			Email: cfg.AuthorEmail,
		},
		Link: &feeds.Link{Href: cfg.Domain.String()},
		Copyright: fmt.Sprintf(
			"Copyright %d–%d %s",
			cfg.Since,
			now.Year(),
			cfg.AuthorName,
		),
		Items: []*feeds.Item{},

		Created: now,
		Updated: now,
	}

	for _, post := range s.posts() {
		body, err := post.render()
		if err != nil {
			return err
		}
		title, err := post.Title()
		if err != nil {
			return err
		}
		feed.Items = append(feed.Items, &feeds.Item{
			Id:          post.Shortname(),
			Title:       title,
			Author:      feed.Author,
			Content:     string(body),
			Description: string(body),
			Link: &feeds.Link{
				Href: fmt.Sprintf("%s/%s.html", feed.Link.Href, post.Shortname()),
			},
			Created: post.meta.createdAt,
			Updated: post.meta.updatedAt,
		})
	}

	return nil
}

func (s substructure) setlinks() error {
	for _, out := range s.docs {
		linksout, err := out.linksout()
		if err != nil {
			return err
		}

		for _, l := range linksout {
			l = strings.TrimSuffix(l, filepath.Ext(l))
			if in := s.Get(l); in != nil {
				out.outgoing = append(out.outgoing, in)
				in.incoming = append(in.incoming, out)
			}
		}
	}

	return nil
}

func (s substructure) Execute(d *Document) error {
	if s.t == nil {
		s.t = &template.Template{}
	}
	if err := loadTemplates(s.t); err != nil {
		return err
	}

	for _, d := range s.docs {
		b, err := d.render()
		if err != nil {
			return err
		}
		s.t.New("body").Parse(string(b))
		var buf bytes.Buffer

		if err := s.t.Lookup("text_document").Execute(&buf, templateVars{d, &s}); err != nil {
			return fmt.Errorf(
				"can't execute document `%s`: %w",
				d.Shortname(),
				err,
			)
		}

		path := filepath.Join("dist", d.Shortname()+".html")
		if ioutil.WriteFile(path, buf.Bytes(), 0644); err != nil {
			return fmt.Errorf(
				"can't write document `%s`: %w",
				d.Shortname(),
				err,
			)
		}
	}

	return nil
}
