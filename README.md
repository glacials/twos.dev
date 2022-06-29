# twos.dev

This is the source code for my personal website. I post thoughts, hobbies, and
other random things here.

## Architecture

twos.dev is a low-tech website. It does not require JavaScript and it is
composed entirely of static files served by GitHub Pages.

These static files are built by a bespoke static website generator called
Winter, which is written in Go and embedded in this repository at `winter/`.
Winter executes a list of transformation functions on each document in order,
each performing a simple task like converting Markdown to HTML or extracting
metadata for later use.

The transformations built into Winter are listed in
[`winter/cmd/build_document.go`](./winter/cmd/build_document.go).

## First run

### Dependencies

- Go 1.18+

### Starting a dev server

```sh
make serve
```

Files will be watched for changes. Changes to documents or graphics will
automatically trigger the right transformations. Changes to the builder
itself will trigger a program recompilation and restart. In both cases, a
WebSocket connection on the local page will listen for the change and trigger a
refresh automatically.

#### Debugging Transformations

If your documents are coming out of the transformation gauntlet wrong and you
don't know which transformation is misbehaving, instead run:

```sh
make debug
```

which will dump the document state after each transformation to

```
dist/debug/DOCUMENT/XX_TRANSFORMATION.html
```

where `DOCUMENT` is the filename of the source document (e.g. `guide.md`),
`TRANSFORMATION` is the name of the transformation function, and `XX` is a
two-digit number representing the order of the transformation among all
transformations.
