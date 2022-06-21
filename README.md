# twos.dev

This is the source code for my personal website. I post thoughts, hobbies, and other
random things here.

## First run

### Dependencies

- Go 1.18+

### Starting a dev server

```sh
make serve
```

Files will be watched for changes and the server will restart or recompile
templates as needed.

## Architecture

twos.dev is a bespoke static website. The build process is captured in the twos.dev
binary, whose code is mostly in `cmd/`. The build steps are listed in
[`cmd/build_document.go`](./cmd/build_document.go).
