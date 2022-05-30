# twos.dev

This is the source code for my personal website. I post thoughts, hobbies, and other
random things here.

## First run

### Dependencies

- Go 1.18+

### Starting a dev server

```sh
go install
twos.dev serve
```

## Architecture

twos.dev is a bespoke static website. The build process is captured in the twos.dev
binary, whose code is mostly in `cmd/`. The build steps include:

- Create gallery containers for photos
- Create thumbnails from photos & replace non-gallery `<img>` tags accordingly
- Build Markdown into HTML
- Copy in static assets
- Clean up small formatting issues like correcting between `--` & `&mdash;`

`twos.dev serve` continuously builds while acting as a file server for the build
directory.
