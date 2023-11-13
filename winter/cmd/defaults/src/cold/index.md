---
type: page
---

# New Winter Website

If you're seeing this as a web page in your browser, your Winter setup works!

This file is located at `src/cold/index.md`. Edit it to your liking and it will
be your homepage.

To create a new page, make a new Markdown or HTML file in `src/cold` and give it
some content. If that page has type `post`, it will be displayed below.

<!--
  To syntax highlight template calls like below, try checking if your editor or
  a plugin has support for .tmpl files. Winter will pick up .html.tmpl files and
  .md.tmpl files the same as their respective counterparts.
-->

{{ range posts }}

- {{ with .Category }}{{ . }}: {{ end }}<a href="{{ .WebPath }}">{{ .Title }}</a> {{ if not .CreatedAt.IsZero }}({{ .CreatedAt.Format "2006 January" }}){{ end }}

{{ else }}

_No posts detected._

{{ end }}

To give a document a type, prepend it with a frontmatter section, written in
YAML and surrounded by `---`, like so:

```markdown
---
type: post
---

# My First Post
```

Any document with no type specified is automatically of type `draft`.
