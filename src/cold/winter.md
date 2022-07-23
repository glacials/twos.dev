---
date: 2022-07-07
filename: winter.html
toc: true
type: page
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

_Warning: Winter is in early alpha and so is this documentation. You may find
inconsistencies or pieces of code present in twos.dev that have not yet been
migrated to Winter._

# Winter

[![Go Reference](https://pkg.go.dev/badge/twos.dev/winter.svg)](https://pkg.go.dev/twos.dev/winter) [GitHub Repository](https://github.com/glacials/twos.dev#winter)

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

### Directory Structure {#layout}

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

#### `winter init` {#init}

Usage: `winter init`

Initialize the current directory for use with Winter. The Winter directory
structured detailed above will be created, and default starting templates will
be populated so that you have a working `index.html` listing posts.

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

A document is an HTML file or a Markdown file with optional frontmatter. The
first level 1 heading (`<h1>` in HTML or `#` in Markdown) will be used as the
document title.

This is an example document called `example.md`:

```markdown
---
date: 2022-07-07
filename: example.html
type: post
---

# My Example Document

This is an example document.
```

### Frontmatter {#frontmatter}

Frontmatter is specified in YAML. All fields are optional.

```markdown
---
filename: example.html
date: 2022-07-07
updated: { { .Now.Format "2006-01-02" } }

category: arbitrary string
toc: true|false
type: post|page|draft
---

# The Thing About Icebergs

...
```

#### `filename` {#filename}

Filename specifies the desired final location of the built file in `dist`. This
must end in `.html` (even if the source document is Markdown) and must not be in
a subdirectory. Winter enforces this because if you later move off Winter, web
paths that end in `.html` and do not use subdirectories will be the easiest to
migrate.

When not set, the filename of the source document minus extension is used in
place. For example, `envy.html.tmpl` and `envy.md` would both become `envy.html`
(though if two source files would produce the same destination file, Winter will
error). The result can be accessed using the [`{{"{{"}} .Shortname }}`](#shortname) template var.

#### `date` {#date}

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

#### `updated` {#updated}

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

#### `toc` {#tocprop}

Whether to render a table of contents (default `false`). If `true`, the table of
contents will be rendered just before the first level 2 header (`<h2>` in HTML,
`##` in Markdown) and will list all level 2, 3, and 4 header in nested `<ul>`s.
See the [top of this page](#toc) for an example.

#### `type` {#type}

The kind of document. Possible values are `post`, `page`, `draft`.

`post` documents are programmatically included in template functions. `page` and
`draft` documents have no programmatic action taken on them.

#### `category` {#category}

The category of the document. Accepts any string. This is exposed to templates
via the `{{"{{"}} .Category }}` field. It has no other effect.

### Templates {#templates}

Templates use the [`text/template`](https://pkg.go.dev/text/template) Go
library.

#### Document Fields {#fields}

Document fields are available on any document.

##### `{{"{{"}} .Category }}` {#category}

_Type: `string`_

The value specified by the frontmatter [`category`](#category) field. This can
be any arbitrary string specified by the document and is not used internally by
Winter.

##### `{{"{{"}} .Dest }}` {#dest}

_Type: `string`_

The path of the document, relative to the web root.

##### `{{"{{"}} .IsDraft }}` {#isdraft}

_Type: `bool`_

Whether the document is a draft (i.e. has frontmatter specifying `type: draft`).

##### `{{"{{"}} .IsPost }}` {#ispost}

_Type: `bool`_

Whether the document is a post (i.e. has frontmatter specifying `type: post`).

##### `{{"{{"}} .Title }}` {#title}

_Type: `string`_

The value of the document's first level 1 heading (`<h1>` for HTML or `#` for Markdown).

##### `{{"{{"}} .CreatedAt }}` {#createdat}

_Type: [`time.Time`](https://pkg.go.dev/time#Time)_

The parsed date specified by the frontmatter [`date`](#date) field.

##### `{{"{{"}} .UpdatedAt }}` {#updatedat}

_Type: [`time.Time`](https://pkg.go.dev/time#Time)_

The parsed date specified by the frontmatter [`updated`](#updated) field.

#### Functions {#functions}

##### `posts` {#posts}

Usage: `{{"{{"}} range posts }} ... {{"{{"}} end }}`

Returns a list of all documents with type `post`, from most to least recent.

See [Document Fields](#fields) for a list of fields available to documents.

##### `archives` {#archives}

Usage: `{{"{{"}} range archives }}{{"{{"}} .Year }}: {{"{{"}} range .Documents }} ... {{"{{"}} end }}{{"{{"}} end }}`

Returns a list of `archive` types ordered from most to least recent. An
`archive` has two fields, `.Year` (integer) and `.Documents` (array of
documents). This allows you to display posts sectioned by year.

See [Document Fields](#fields) for a list of fields available to documents.

##### `img` {#img}

_Alias: `imgs`_

Usage: `{{"{{"}} img[s] <caption> <imageshortname alttext>... }}`

Render an image or images with a single caption. Image files must be present in
this format:

```plain
public/img/<pageshortname>-<imageshortname>[-<light|dark>].<png|jpg|jpeg>
```

For example, to render one image with a caption and alt text:

```template
<!-- mypage.md -->

{{"{{"}} img
   "A caption."
   "test1"
   "Descriptive alt text of what the image is of, for assistive tech"
}}
```

Result:

{{ img
   "A caption."
   "test1"
   "Descriptive alt text of what the image is of, for assistive tech"
}}

In this example, the image file must hold one or more of these forms:

- `public/img/mypage-test1.jpg`
- `public/img/mypage-test1-light.jpg`
- `public/img/mypage-test1-dark.jpg`
- `public/img/mypage-test1.jpeg`
- `public/img/mypage-test1-light.jpeg`
- `public/img/mypage-test1-dark.jpeg`
- `public/img/mypage-test1.png`
- `public/img/mypage-test1-light.png`
- `public/img/mypage-test1-dark.png`

If `-light` and/or a `-dark` variants exist, they will be used when the user is
in the respective dark mode setting.

Any number of images can be rendered together with one caption beneath the
group by passing multiple images and alt texts. They will appear next to each
other when the page width allows it, or stacked vertically otherwise.

```template
---
filename: example.html
---

{{"{{"}} imgs
   "A pair of images."
   "test1"
   "Descriptive text of test1"
   "test1"
   "Descriptive text of test2"
}}
```

Result:

{{ imgs
   "A pair of images."
   "test1"
   "Descriptive text of imagename1"
   "test2"
   "Descriptive text of imagename2"
}}

##### `video` {#video}

_Alias: `videos`_

Usage: `{{"{{"}} video[s] <caption> <videoshortname alttext>... }}`

Behaves exactly as `img` but searches for `mp4`/`mov` files instead and renders
them in `<video>` tags. Note that most browsers do not currently support light
or dark mode variations for videos, so the wrong variant may be displayed.
