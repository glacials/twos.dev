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

- Go 1.23+

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

## Deploying

Push `main` normally:

```sh
git push
```

The tracked GitHub Actions workflow runs on every push, checks out the existing
`gh-pages` branch as `dist`, and runs `winter publish --allow-dirty`. Publishing
is incremental and append-only: unchanged generated images and static files are
reused, changed files are overlaid, and removed source files do not silently
remove old public URLs. Winter checks every local link and resource before it
commits and pushes the deployment branch.

This workflow is stored in the repository, so the automatic deployment behavior
works on every clone without installing a local Git hook. If an immediate local
deployment is useful, install Winter v0.6.1 or newer and run `winter publish`;
that is optional, not part of the normal deployment routine.

A weekly and manually dispatchable audit performs a clean build and compares it
with the persistent deployment. It is intentionally separate from normal pushes
because regenerating every image is slow.
