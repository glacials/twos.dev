# twos.dev

This is the source code for my personal website. I post thoughts, hobbies, and
other random things here.

## Architecture

twos.dev is a low-tech website. It does not require JavaScript and it is
composed entirely of static files served by GitHub Pages.

These static files are built by a bespoke static website generator called
Winter, which is written in Go and embedded in this repository at `winter/`. See
more about Winter below.

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

## Winter

[![Go Reference](https://pkg.go.dev/badge/twos.dev/winter.svg)](https://pkg.go.dev/twos.dev/winter)

Winter is the bespoke static website generator that powers twos.dev. It can be
used to power your static website as well, as a CLI or Go library. Beware that
although it is technically website-agnostic, its feature set is shaped to serve
my peculiarities.

### Installation

```sh
# CLI
go install twos.dev/winter/cmd@latest
winter --help

# Go library
go get -u twos.dev/winter@latest
```

### Philosophy

The Winter CLI has three main actions:

```sh
winter build               # Build site once and stop
winter build --serve       # Build site continuously and serve results
winter freeze shortname... # Convert the given document(s) from warm to cold
```

#### Warm vs. Cold Documents

Winter's goals are to ease writing and editing of new content, and to harden
existing content against breakages. Normally these goals work against each
other. In Winter, the two types of content coexist and documents can flow easily
from easy-to-edit ("warm") to hard-to-break ("cold").

**Warm** content is unstable. It should be synchronized into the project
directory from an outside tool---which tool is up to the user. I use [iA
Writer](https://ia.net/writer) hooked up to an iCloud folder, with a Shortcut
and cron job to sync files into the repository. This makes creating content easy
from any platform in any state of mind, while still publishing the (unlinked)
content for preview, or for sending to a friend for review.

**Cold** content is stable. It must not be sourced from anywhere automatically,
because sync jobs are a great way to accidentally overwrite content. This is
finished content you want to last for years or decades without needing to
babysit its existence.

Winter is very early software, so please use with caution.
