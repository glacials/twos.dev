---
filename: meta.html
date: 2022-06-17
---

# How I Write

twos.dev has a unique advantage over anything I’ve ever worked on: it has just one developer, me, so I can design the developer experience in obnoxious and unintuitive ways to serve my peculiarities.

With this in mind, I decided to ruthlessly optimize it around two goals:

1. Make it easy to publish
2. URLs must never change

To fulfill these goals I divide my content into two types, warm and cold. Warm content is easy to create and edit, while cold content is hard to break.

## Warm Content

I previously wrote in my notes app while drafting, then manually migrated the draft to HTML when the piece started getting serious. From then until publish I would work directly in the HTML. This was surprisingly refreshing as it allowed me to have wild one-off customizations for individual posts, but it had three problems:

- Writing HTML has a higher context switch cost (sitting at my computer in my editor vs. scribbling quick drafts from my phone).
- Writing was slowed down by cruft like `<p>` and `<li>`
- The reading experience may flow differently on a web page than in my notes app, leading to sometimes large rewrites right after transferring

To solve these issues and work towards my first goal, I needed to close the gap in tooling between the draft phase and the publish phase.

I found [iA Writer](https://ia.net/writer) to replace my note-taking app for draft writing. It gets my thoughts out quickly, has helpful tools for e.g. reducing filler words, and looks nice to boot. Importantly, it stores writing as Markdown files.

### Shortcuts

iOS and macOS ship with a first-party app called [Shortcuts](https://apps.apple.com/us/app/shortcuts/id1462947752), which is a no-code, event-driven automation framework. Using Shortcuts, I set up an automation that triggers when I switch away from the iA Writer app. The automation adds a 1-2 line YAML frontmatter section to each document, then pushes it to the `src/warm` directory in the twos.dev Git repository by invoking [Working Copy](https://workingcopyapp.com).

### Preprocessing

On push, GitHub Actions builds the Markdown file into HTML with all the [right preprocessing](https://github.com/glacials/twos.dev/blob/main/cmd/build_document.go) required to make it look like twos.dev. I use a small subset of [`html/template`](https://pkg.go.dev/html/template) to allow myself to put things like MP4 screen recordings and dark-mode-aware images into Markdown. This is also where the frontmatter from before is stripped and parsed into metadata inserted elsewhere on the page.

Technically at this point the page is now published, just not linked to from anywhere. I can see how it looks in the context of twos.dev, and send the link to friends who have offered to review.

## Cold Content

There are 2-3 needs warm content doesn't cover: first, this odd pipeline that uses my phone for CI is not something I yet have high confidence in, so I'd like to limit its exposure.

Second, once in a while I need some wacky one-off piece of code code for a single post, like the CSS-only spectrum in [Anonymously Autistic](autism.html) or the variable-width font requirements of [Dashes](dashes.html). I prefer to write code in `$EDITOR`, not iA Writer. 

One great way to accomplish both of these needs is to simply remove the file from the iA Writer sync directory:

```sh
git mv src/{warm,cold}/DOCUMENT.md
```

From here, I may also take one final step of permanently converting the document to HTML, if I find it will be easier to write code like that than embedding it in the Markdown.

```sh
twos.dev build
git rm src/cold/DOCUMENT.md
cp dist/DOCUMENT.html src/cold/DOCUMENT.html
git add !$
```

This brings back the weight of editing prose in HTML, but by this point most of the writing is done.

#### Results

Allowing myself this escape hatch is freeing. I’m more encouraged to write interactive or otherwise bespoke components, and twos.dev becomes consistent by default but I can break that consistency when I need (e.g. my [CV](cv.html)'s bicolumnar layout)

## On URLs

I use my website as a test bed for interesting technology. As a side effect, content has been lost over the years as it becomes hard to migrate to newly-overhauled versions.

It so happens that by keeping my writing in these two formats, Markdown and raw HTML/CSS, I make it easy to accomplish my second goal: [cool URLs don’t change](https://www.w3.org/Provider/Style/URI). I’ll never migrate JavaScript frameworks and be too lazy to move things forward, or move to a database-backed writing system while being unable to simply `cp -r` these files into the `public` directory.

If I’ve done things right, this web page will be accessible at twos.dev/meta.html until the day I die [and then some](death.html).