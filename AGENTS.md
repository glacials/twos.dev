# Repository Guidelines

This repo hosts the twos.dev static website. Pages are authored in Markdown/HTML and rendered by the Winter generator into static files.

## Project Structure & Module Organization
- `src/`: Source content and templates.
  - `src/cold/` and `src/warm/`: Markdown articles/posts (`.md`).
  - `src/templates/`: Go-Template HTML partials/layouts (`.html.tmpl`, leading `_` for partials).
  - `src/img/`, `src/favicon/`: Source assets.
- `public/`: Hand-managed static files copied as-is (CSS, images, favicon, etc.).
- `dist/`: Build output (generated site). Do not edit.
- `Makefile`: Developer commands. `package.json`: Prettier + go-template plugin.

## Build, Test, and Development Commands
- `make serve`: Start Winter dev server with live rebuilds.
- `make build`: Build static site into `dist/` (honors `WINTER_ARGS`). Example: `WINTER_ARGS="--drafts" make build`.
- `make tools`: Install Winter CLI (`go install twos.dev/winter@latest`). Ensure Go 1.18+.
- Optional formatting: `npx prettier --write "src/**/*.{md,html,tmpl}"`.

## Coding Style & Naming Conventions
- Use Prettier 3 with `prettier-plugin-go-template` (configured in `package.json`).
- Prefer lowercase, hyphenated filenames for content (e.g., `my-post.md`).
- Templates: partials start with `_` (e.g., `_nav.html.tmpl`).
- Indentation: 2 spaces; wrap lines thoughtfully for readability.

## Testing Guidelines
- No formal test suite. Validate changes by running `make serve`, browsing the local site, and checking generated `dist/` diffs.
- Verify internal links, images, and template includes render as expected.

## Commit & Pull Request Guidelines
- Commit messages: short, imperative; optionally prefix with scope or file (e.g., `docs: fix link`, `index.html: Add Bluesky link`).
- PRs should include: concise description, rationale, screenshots for visual changes, and linked issues (if any).
- Keep changes focused; avoid editing generated files in `dist/`.

## Architecture Overview
- Winter builds static pages from `src/` and copies assets from `public/`. No runtime JavaScript is required; the site is served by GitHub Pages.
