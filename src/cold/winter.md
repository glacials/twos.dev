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

[![Go Reference](https://pkg.go.dev/badge/twos.dev/winter.svg)](https://pkg.go.dev/twos.dev/winter) [{{ icon
    "github"
    "GitHub logo"
}}](https://github.com/glacials/twos.dev#winter)

_Winter_ is the bespoke static website generator that powers twos.dev. It can
power your static website as well, either as a CLI tool or Go library. Winter is
strongly opinionated to serve my pecularities around building twos.dev.

On the surface, Winter is similar to other static website
generators. It can take a directory of files of varying content
types—Markdown, Org, or HTML—apply some templating, and render them to
HTML files ready for deployment to a static web server.

Winter is different because of its unique design goals.

Winter's design goals are to allow content to be both easy to edit and hard to
break. These goals work against each other by default, so Winter splits content
into two modes: **warm** and **cold**.

**Warm** content should be synchronized into the project directory from an
outside tool of your choice. Generally this is a writing program such as [iA
Writer](https://ia.net/writer) hooked up to a shared or cloud directory;
simplistically, it can be a cron job that runs `rsync`.

**Cold** content should not be touched by automated tools. This content must be
preseved for years or decades, so less exposed surface area is better. When a
piece of warm content is stable, it can be "frozen" into cold content using
`winter freeze`.

## Installation {#installation}

```sh
# CLI
go install twos.dev/winter@latest
winter --help

# Go library
go get -u twos.dev/winter@latest
```

## Documentation {#documentation}

This section details the `winter` CLI documentation. For Go library
documentation, see
[pkg.go.dev/twos.dev/winter](https://pkg.go.dev/twos.dev/winter).

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

Pass `--help` to any command to see detailed usage and behavior.

#### `winter init` {#init}

Usage: `winter init`

Initialize the current directory for use with Winter. The Winter directory
structure detailed above will be created, and default starting templates will
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

A document is an HTML, Markdown, or Org file with optional
frontmatter. The first level 1 heading (`<h1>` in HTML, `#` in
Markdown, `*` in Org) will be used as the document title.

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

Frontmatter for HTML and Markdown documents is specified in YAML.
Frontmatter for Org files is specified using Org keywords of
equivalent names (in whatever case you choose). All fields are
optional.

The available frontmatter fields for HTML and Markdown are:

```yaml
filename: example.html
date: 2022-07-07
updated: 2022-11-10
category: arbitrary string
toc: true|false
type: post|page|draft
```

The same fields are available in Org files:

```org
#+FILENAME: example.html
#+DATE: 2022-07-07
#+UPDATED: 2022-11-10
#+CATEGORY: arbitrary string
#+TOC: true|false
#+TYPE: post|page|draft
```

See below for details of each.

#### `filename` {#filename}

Filename specifies the desired final location of the built file in
`dist`. This must end in `.html` no matter the source document type
and must not be in a subdirectory. Winter enforces this because if you
later move off Winter, web paths that end in `.html` and do not use
subdirectories will be the easiest to migrate.

When not set the filename of the source document is used, with any
extensions replaced with `.html`. For example, `envy.html.tmpl` and
`envy.md` would both become `envy.html` (though if two source files
would produce the same destination file, Winter will error).

A document's **shortname** is defined as its extensionless filename.
The shortname is accessible to templates using the [`{{"{{"}}
.Shortname }}`](#shortname) template variable. The filename is not
accessible to templates, but can be calculated by appending `.html` to
the shortname.

#### `date` {#date}

Date is the publish date of the document written as `YYYY-MM-DD`. It
is available to templates as a Go
[`time.Time`](https://pkg.go.dev/time#Time) using [`{{"{{"}}
.CreatedAt }}`](#createdat).

When not set the document will not have a publish date attached to it.

#### `updated` {#updated}

Updated is the date the document was last meaningfully updated, if
any, written as `YYYY-MM-DD`. It is available to templates as a Go
[`time.Time`](https://pkg.go.dev/time#Time) using [`{{"{{"}}
.UpdatedAt }}`](#updatedat).

When not set the document will not have an update date attached to it.

#### `toc` {#tocprop}

Whether to render a table of contents (default `false`). If `true`,
the table of contents will be rendered just before the first level 2
header (`<h2>` in HTML, `##` in Markdown, `**` in Org) and will list
all level 2, 3, and 4 headers. See the [top of this page](#toc) for an
example.

#### `type` {#type}

The kind of document. Possible values are `post`, `page`, `draft`.

`post` documents are programmatically included in template functions.
`page` and `draft` documents have no programmatic action taken on
them and will not be discoverable unless linked to.

The default template set provided by `winter init` gives each a
slightly different visual treatment:

- Posts have a publish date and update date at the top and bottom of
  the page
- Pages have a publish date and update date at the bottom of the page
- Drafts behave like posts but have a large "DRAFT" banner at the top
  of the page

Beyond this, the types are only different in semantics.

#### `category` {#category}

The category of the document. Accepts any string. This is exposed to templates
via the [`{{"{{"}} .Category }}`](#cat) field.

The default template set provided by `winter init` gives a minor
visual treatment to the listing and display of documents with
categories.

Otherwise, the category is semantic. There is no programmatic access
to documents by category.

### Templates {#templates}

Templates use the [`text/template`](https://pkg.go.dev/text/template) Go
library.

#### Document Fields {#fields}

The following fields are available to templates rendering documents.

##### `{{"{{"}} .Category }}` {#cat}

_Type: `string`_

The value specified by the frontmatter [`category`](#category) field.
This is an arbitrary string.

##### `{{"{{"}} .Dest }}` {#dest}

_Type: `string`_

The filesystem path of the document after having been rendered,
relative to the web root.

##### `{{"{{"}} .IsDraft }}` {#isdraft}

_Type: `bool`_

Whether the document is a draft (i.e. has frontmatter specifying `type: draft`).

##### `{{"{{"}} .IsPost }}` {#ispost}

_Type: `bool`_

Whether the document is a post (i.e. has frontmatter specifying `type: post`).

##### `{{"{{"}} .Shortname }}` {#shortname}

_Type: `string`_

The shortname of the document. This is the extensionless
[filename](#filename) frontmatter attribute if set, or the source
document's extensionless filename otherwise.

##### `{{"{{"}} .Title }}` {#title}

_Type: `string`_

The value of the document's first level 1 heading (`<h1>` in HTML, `#` in Markdown, `*` in Org).

##### `{{"{{"}} .CreatedAt }}` {#createdat}

_Type: [`time.Time`](https://pkg.go.dev/time#Time)_

The publish date of the document as specified by the frontmatter
[`date`](#date) attribute.

This value can be formatted using Go's [`func (time.Time)
Format`](https://pkg.go.dev/time#Time.Format) function:

```template
{{"{{"}} .CreatedAt.Format "2006 January" }} <!-- Renders 2022 July  -->
{{"{{"}} .CreatedAt.Format "2006-01-02" }}   <!-- Renders 2022-07-08 -->
```

Use `{{"{{"}} .CreatedAt.IsZero }}` to see if the date was not set.
You can use this to hide unset dates:

```template
{{"{{"}} if not .CreatedAt.IsZero }}
  published {{"{{"}} .CreatedAt.Format "2006 January" }}
{{"{{"}} end }}
```

##### `{{"{{"}} .UpdatedAt }}` {#updatedat}

_Type: [`time.Time`](https://pkg.go.dev/time#Time)_

The date the document was last meaningfully updated, as specified by
the frontmatter [`updated`](#updated) attribute.

This value behaves identically to [`{{"{{"}} .CreatedAt
}}`](#createdat). For details on dealing with it, see that
documentation.

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

If `-light` and/or `-dark` variants exist, they will be used when the user is
in the respective dark mode setting.

Any number of images can be rendered together with one caption beneath the
group by passing multiple images and alt texts. They will appear next to each
other when the page width allows it, or stacked vertically otherwise.

```template
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

Behaves exactly as `img` but searches for `mp4`/`mov` files instead
and renders them in `<video>` tags.

Note that most browsers do not currently support varying a video
source by light/dark mode setting, so the wrong variant may be
displayed.
