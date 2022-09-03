---
date: 2022-06-18
filename: meta.html
type: post
---

# Inside twos.dev

There's a unique advantage in working on twos.dev over anything I’ve worked on: it has just one developer, me. I can therefore design the developer experience in obnoxious and unintuitive ways to serve my peculiarities.

With this in mind, I decided to ruthlessly optimize for two goals:

1. It must be easy to publish writing
2. As long as I'm alive, URLs must not change

In some ways these goals work against each other; easy creation usually means easy destruction. So I divide my content into two types, warm and cold. Warm content is easy to create and edit, while cold content is hard to break.

## Warm Content

I used to write drafts in my notes app then migrate to HTML when the piece became less amorphous; from then on I would write and edit directly in HTML. This was surprisingly refreshing as it allowed me to have wild one-off customizations for individual posts, but it had three problems:

- Writing HTML is a rough context switch (what was scribbling drafts on my phone became sitting at my computer in `$EDITOR`)
- Writing was slowed down by cruft like `<p>` and `<li>`
- The reading experience can flow differently on a web page than in a notes app, leading to large rewrites after transferring

To solve these issues and work towards my first goal, I needed tooling in the space between amorphous notes and strict HTML.

I first focused on changing my drafting app to something with more robust exporting tools. I found [iA Writer](https://ia.net/writer), which has been great to get my thoughts out quickly. It looks nice and has tools for reducing filler words and cliches. Importantly, it stores writing as Markdown files in an online drive of my choosing.

### Shortcuts

iOS and macOS ship with [Shortcuts](https://apps.apple.com/us/app/shortcuts/id1462947752), a no-code event-driven automation app. Using Shortcuts, I set up an automation that triggers when I switch away from the iA Writer app. The automation adds a 1-2 line YAML frontmatter section to each document, then pushes it to the `src/warm` directory in the twos.dev Git repository by invoking [Working Copy](https://workingcopyapp.com).

### Preprocessing

On push, GitHub Actions builds the Markdown file into HTML with all the [right preprocessing](https://github.com/glacials/twos.dev/blob/main/cmd/build_document.go) required to make it look like twos.dev. I use a small subset of [`html/template`](https://pkg.go.dev/html/template) to allow myself to put things like MP4 screen recordings and dark-mode-aware images into Markdown. This is also where the frontmatter from before is stripped and parsed into metadata inserted elsewhere on the page.

Technically at this point the page is now published, just not linked to from anywhere. I can see how it looks in the context of twos.dev, and send the link to friends who have offered to review.

## Cold Content

There are 2-3 needs warm content doesn't cover: first, this odd pipeline that uses my phone for CI is not something I have high confidence in, so I'd like to limit its blast radius to content that benefits from it---content I'm actively working on---and keep safe the innocent bystanders that are historical content.

Second, once in a while I need some wacky one-off piece of code code for a single post, like the CSS-only spectrum in [Anonymously Autistic](autism.html) or the variable-width font requirements of [Advanced Dashes](dashes.html). iA Writer is great for prose, but when it's time to write code I need to be back in `$EDITOR`.

To accomplish these two needs I simply migrate the file out of iA Writer and into a plain directory meant for these pages that have "graduated":

```sh
git mv src/{warm,cold}/DOCUMENT.md
```

From here I may also take one final step of permanently converting the document to HTML if I find myself building more than a little code into it. I can always embed HTML in the Markdown, but taking this extra step gives me all my editor doodads and hardens the file against future changes to preprocessing.

```sh
winter build
mv dist/DOCUMENT.html src/cold/DOCUMENT.html
git rm src/cold/DOCUMENT.md
git add src/cold/DOCUMENT.html
```

_([Winter](https://twos.dev/winter) is the bespoke CLI that builds twos.dev.)_

This brings back the weight of editing prose in HTML, but these cases are the minority and at this point most of the writing is done anyway.

#### Results

Allowing myself this escape hatch is freeing. I’m more encouraged to write interactive or otherwise bespoke components, and twos.dev becomes consistent by default but I can break that consistency when I need (e.g. for a [CV](cv.html)'s bicolumnar layout).

## On URLs Not Changing

I use my website as a test bed for interesting technology. As a side effect, content has been lost over the years as it becomes hard to migrate to overhauled versions.

It so happens that by keeping my writing in Markdown or raw HTML/CSS, I make it easy to accomplish my second goal: [cool URLs don’t change](https://www.w3.org/Provider/Style/URI). I've removed the possibility that changing frontend or backend frameworks will leave my documents behind; Markdown is not going anywhere, and all the HTML files need to look decent is a header and footer. In the worst case scenario I'll simply `cp -r dist` to my next crazy system's web root.

If I’ve done things right, this web page will be accessible at twos.dev/meta.html for at least as long as I live.
