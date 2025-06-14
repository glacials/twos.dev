# twos.dev

This is the source code for my personal website. I post thoughts, hobbies, and
other random things here.

## Architecture

twos.dev is a low-tech website. It does not require JavaScript and it is
composed entirely of static files served by GitHub Pages.

These static files are built by
[Winter](https://github.com/glacials/winter).

## First run

### Dependencies

- Go 1.18+

### Starting a dev server

```sh
go install
go install github.com/mitranim/gow@latest
make serve
```

Files will be watched for changes. Changes to documents or graphics will
automatically trigger the right transformations. Changes to the generator
itself will trigger a program recompilation and restart. In both cases, a
WebSocket connection on the local page will listen for the change and trigger a
refresh automatically.

## Winter

[![Go Reference](https://pkg.go.dev/badge/twos.dev/winter.svg)](https://pkg.go.dev/twos.dev/winter)

Winter is the bespoke static website generator that powers twos.dev. It can be
used to power your static website as well, as a CLI or Go library. See the
[winter README](https://github.com/glacials/winter) for details.
