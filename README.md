# twos.dev

This is the source code for my personal website. I post thoughts, hobbies, and
other random things here.

## Architecture

twos.dev is a low-tech website. It does not require JavaScript and it is
composed entirely of static files served by GitHub Pages.

These static files are built by a bespoke static website generator called
Winter, which is written in Go and embedded in this repository at `winter/`.

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

## Using Winter

Winter is written to be used with twos.dev, but if you want to use it elsewhere
you can:

```sh
# Use the `winter cmd`
go install twos.dev/winter/cmd@latest

# Use the twos.dev/winter library
go get -u twos.dev/winter@latest
```

[See documentation](https://pkg.go.dev/twos.dev/winter)
