---
date: 2022-07-07
filename: winter.html
toc: true
type: draft
updated: 2022-07-13
---

<meta
    name="go-import"
    content="twos.dev
             git https://github.com/glacials/twos.dev"
/>
<meta
    name="go-source"
    content="twos.dev
             https://github.com/glacials/twos.dev
             https://github.com/glacials/twos.dev/tree/main{/dir}
             https://github.com/glacials/twos.dev/blob/main{/dir}/{file}#L{line}"
/>

# Winter

[![Go Reference](https://pkg.go.dev/badge/twos.dev/winter.svg)](https://pkg.go.dev/twos.dev/winter)

_Warning: Winter is in early alpha and so is this documentation. You may find
inconsistencies or pieces of code present in twos.dev that have not yet been
migrated to Winter._

_Winter_ is the bespoke static website generator that powers twos.dev. It can
power your static website as well, either as a CLI tool or Go library. Winter is
strongly opinionated to serve my pecularities around building twos.dev.

## Installation {#installation}

```sh
# CLI
go install twos.dev/winter/cmd@latest
winter --help

# Go library
go get -u twos.dev/winter@latest
```

## Documentation {#documentation}

This section details the `winter` CLI documentation. For Go library
documentation, see
[pkg.go.dev/twos.dev/winter](https://pkg.go.dev/twos.dev/winter).

Winter's design goals are to allow content to be both easy to edit and hard to
break. These goals work against each other by default, so Winter splits content
into two modes: **warm** and **cold**.

**Warm** content should be synchronized into the project directory from an
outside tool of your choice. Generally this is a writing program such as [iA
Writer](https://ia.net/writer) hooked up to a shared or cloud directory;
simplistically, it can be a cron job that runs `rsync` followed by a Git add and
push.

**Cold** content should not be touched by automated tools. This content must be
preseved for years or decades, so less exposed surface area is better. When a
piece of warm content is stable, it can be "frozen" into cold content using
`winter freeze`.

### Directory Layout {#layout}

- `./src`—Content to be built
  - `./src/cold`—Stable content
    - `*.md`—Markdown files
    - `*.html`—HTML files
    - `*.html.tmpl`—HTML files compatible with Go's [`text/template`](https://pkg.go.dev/text/template) package
  - `./src/warm`—Unstable content
    - `*.md`—Markdown files
    - `*.html`—HTML files
    - `*.html.tmpl`—HTML files compatible with Go's [`text/template`](https://pkg.go.dev/text/template) package
  - `./src/templates`—Reusable content
    - `text_document.html.tmpl`—Default page container
    - `imgcontainer.html.tmpl`—Gallery container
    - `*.html`—HTML templates
    - `*.html.tmpl`—HTML templates compatible with Go's [`text/template`](https://pkg.go.dev/text/template) package
  - `./src/img`—Gallery images
    - `...`—Any directory structure
- `./public`—Static files copied directly to the build directory
- `./dist`—Build directory

### Commands {#commands}

#### `winter build` {#build}

Usage: `winter build [--serve]`

Build all content into `dist`. When `--serve` is passed, a fileserver is stood
up afterwards pointing to `dist`, content is continually rebuilt as it changes,
and the browser automatically refreshes.

#### `winter freeze` {#freeze}

Usage: `winter freeze shortname...`

Freeze all arguments, specified by shortname. This moves the files from
`src/warm` to `src/cold` and reflects the change in Git.

### Documents {#documents}

A document is an HTML file or a Markdown file with optional frontmatter. This is
an example document called `example.md`:

```markdown
---
date: 2022-07-07
filename: example.html
type: post
---

# My Example Document

This is an example document.
```

#### Frontmatter {#frontmatter}

Frontmatter is specified in YAML. All fields are optional.

##### `filename` {#filename}

Filename specifies the desired final location of the built file in `dist`. This
must end in `.html` and must not be in a subdirectory.

When not set, the filename of the document minus extension is used in place. For
example, `envy.html.tmpl` and `envy.md` would both become `envy.html`. In
templates, the extension is stripped and the remainder is coalesced to `{{"{{"}} .Shortname }}`.

##### `date` {#date}

The publish date of the document as a Go
[`time.Time`](https://pkg.go.dev/time#Time). Coalesces to `{{ .CreatedAt }}` in
templates.

Templates can format the time using Go's [`func (time.Time) Format`](https://pkg.go.dev/time#Time.Format) function, which accepts a string
of the reference time `01/02 03:04:05PM '06 -0700`. For example, for a document
dated 2022-07-08:

```template
{{"{{"}} .CreatedAt.Format "2006 January" }} <!-- Renders 2022 July  -->
{{"{{"}} .CreatedAt.Format "2006-01-02" }}   <!-- Renders 2022-07-08 -->
```

Use `{{"{{"}} .CreatedAt.IsZero }}` to see if the date was not set. You can use this
to hide unset dates:

```template
{{"{{"}} if not .CreatedAt.IsZero }}
  published {{"{{"}} .CreatedAt.Format "2006 January" }}
{{"{{"}} end }}
```

##### `updated` {#updated}

The date the document was last meaningfully updated, if any, as a Go
[`time.Time`](https://pkg.go.dev/time). Coalesces to `{{"{{"}} .UpdatedAt }}` in
templates.

Use `{{"{{"}} .UpdatedAt.IsZero }}` to see if the date was not set. You can use this
to hide unset dates:

```template
{{"{{"}} if not .CreatedAt.IsZero }}
  published {{"{{"}} .CreatedAt.Format "2006 January" }}
  {{"{{"}} if not .UpdatedAt.IsZero }}
    / last updated {{"{{"}} .UpdatedAt.Format "2006 January" }}
  {{"{{"}} end }}
{{"{{"}} end }}
```

Renders:

```html
published 2022 July / last updated 2022 August
```

##### `toc` {#tocprop}

Whether to render a table of contents (default `false`). If `true`, the table of
contents will be rendered just before the first level 2 header (`<h2>` in HTML,
`##` in Markdown) and will list all level 2 and 3 headers. See the [top of this
page](#toc) for an example.

##### `type` {#type}

The kind of document. Possible values are `post`, `page`, `draft`.

`post` documents are programmatically included in template functions. `page` and
`draft` documents have no programmatic action taken on them.

#### Templates {#templates}

Templates use the [`text/template`](https://pkg.go.dev/text/template) Go
library. Available variables are in flux so are not yet documented, but can be
viewed in [`winter/template.go`](https://github.com/glacials/twos.dev/blob/main/winter/template.go).

Template functions available are:

##### `posts` {#posts}

Returns a list of all documents with type `post`, from most to least recent.

##### `archives` {#archives}

Returns a list of `archive` types ordered from most to least recent. An
`archive` is a struct with two fields, `.Year` (integer) and `.Documents` (array
of `document`s). This allows you to display posts sectioned by year.

##### `img` {#img}

_Alias: `imgs`_

Render an image or images with optional alt text and a caption.

Usage: `{{"{{"}} img caption <imgname alt>... }}`

For example, to render one image with a caption and alt text:

```template
---
filename: mypage.html
---

{{"{{"}} img
   "A caption."
   "myimg"
   "Descriptive alt text of what the image is of, for assistive tech"
}}
```

Images must be present in `public/img` in the form:

```plain
pageshortname-imageshortname[-<light|dark>].<png|jpg>
```

When the `img` function is called from a template, `public/img` is searched for
an image named in this format. In the example above, one possible image name
that would be found is `public/img/mypage-myimg.jpg`.

If both a `-light` and a `-dark` image exist, the correct one will be used for
the user's dark mode preference.

Any number of images can be rendered next to each other with one caption
beneath the group:

```template
---
filename: example.html
---

{{"{{"}} imgs
   "A pair of images."
   "imagename1"
   "Descriptive text of imagename1"
   "imagename2"
   "Descriptive text of imagename2"
}}
```

##### `video` {#video}

_Alias: `videos`_

Behaves exactly as `img` but searches for `mp4`/`mov` files instead. Note that
most browsers do not currently support light or dark mode variations of
`<video>` tags.
