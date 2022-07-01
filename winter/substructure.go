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
	docs []*document
	t    *template.Template
}

// Discover is the first step of the static website build process. It crawls the
// website root looking for files that need to be wrapped up in the build
// process, such as Markdown files, templates, and static assets.
//
// It then learns what it can about each file and stores the information in
// a substructure.
func discover() (s substructure, err error) {
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
		doc, err := htmlDoc(h)
		if err != nil {
			return substructure{}, err
		}
		s.docs = append(s.docs, doc)
	}

	return
}

func (s substructure) get(shortname string) *document {
	for _, d := range s.docs {
		if d.meta.shortname == shortname {
			return d
		}
	}
	return nil
}

func (s substructure) posts() (u []*document) {
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
		feed.Items = append(feed.Items, &feeds.Item{
			Id:          post.meta.shortname,
			Title:       post.meta.title,
			Author:      feed.Author,
			Content:     string(body),
			Description: string(body),
			Link: &feeds.Link{
				Href: fmt.Sprintf("%s/%s.html", feed.Link.Href, post.meta.shortname),
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
			if in := s.get(l); in != nil {
				out.outgoing = append(out.outgoing, in)
				in.incoming = append(in.incoming, out)
			}
		}
	}

	return nil
}

func (s substructure) parent(d *document) *document {
	p, _, ok := strings.Cut(d.meta.shortname, "_")
	if !ok {
		return nil
	}

	return s.get(p)
}

func (s substructure) execute(d *document) error {
	partials, err := filepath.Glob(fmt.Sprintf("src/templates/_*.html.tmpl"))
	if err != nil {
		return fmt.Errorf("can't glob for partials: %w", err)
	}

	for _, partial := range partials {
		name := filepath.Base(partial)
		name = strings.TrimPrefix(name, "_")
		name = strings.TrimSuffix(name, ".html.tmpl")
		p := s.t.New(name)

		s, err := ioutil.ReadFile(partial)
		if err != nil {
			return fmt.Errorf(
				"can't read partial `%s`: %w",
				partial,
				err,
			)
		}

		if _, err := p.Parse(string(s)); err != nil {
			return fmt.Errorf(
				"can't parse partial `%s`: %w",
				partial,
				err,
			)
		}
	}

	for _, d := range s.docs {
		b, err := d.render()
		if err != nil {
			return err
		}
		s.t.New("body").Parse(string(b))
		var buf bytes.Buffer

		if err := s.t.Execute(&buf, templateVars{d, &s}); err != nil {
			return fmt.Errorf(
				"can't execute document `%s`: %w",
				d.meta.shortname,
				err,
			)
		}

		path := filepath.Join("dist", d.meta.shortname+".html")
		if ioutil.WriteFile(path, buf.Bytes(), 0644); err != nil {
			return fmt.Errorf(
				"can't write document `%s`: %w",
				d.meta.shortname,
				err,
			)
		}
	}

	return nil
}
