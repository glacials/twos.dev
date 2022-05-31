---
filename: meta.html
date: 2022-05
---

# The twos.dev CI/CD Pipeline

```
             ┌─────┐
      ┌──────│ Me  │─────┐
      ▣      └─────┘     ▣
   ┌─────┐         ┌───────────┐
   │ Vim │  ┌──────│ iA Writer │
   ├─────┤  │      └───────────┘
┌──│ Git │  │            │
│  └─────┘  │            □
│           ■      ┌───────────┐
│      ┌────────┐  │ Shortcuts │
│   ┌──│ iCloud │  │  for iOS  │──┐
│   │  └────────┘  └───────────┘  │
│   │                             │
│   │           Shortcuts         │
│   │  ┌───────────────────────┐  │
│   ├─■│ Backfill frontmatter  │□─┤
│   │  ├───────────────────────┤  │
│   └─■│ Git add, commit, push │□─┘
│      └───────────────────────┘
│                  │
│                  ▣
│          ┌──────────────┐
│          │ Working Copy │
│          │   for iOS    │
│          │              │
│          └──────────────┘
│                  │
│                  ▣
│  ┌──────────────────────────────┐
└─▣│ github.com/glacials/twos.dev │──┐
   └──────────────────────────────┘  │
                                     │
            GitHub Actions           │
     ┌──────────────────────────┐    │
     │ Render Markdown to HTML  │▣───┤
     ├──────────────────────────┤    │
     │ Copy in high touch pages │▣───┤
     ├──────────────────────────┤    │
     │  Deploy to GitHub Pages  │▣───┘
     └──────────────────────────┘
                   │
                   ■
            ┌────────────┐
            │  twos.dev  │
            └────────────┘

   ┌──────────────────────────────────┐
   │ Key                              │
   │                                  │
   │ A ───■ B   Data flows from A → B │
   │ A ───□ B   A invokes B           │
   │                                  │
   └──────────────────────────────────┘
```

The twos.dev CI/CD pipeline has a unique advantage against anything I’ve ever worked on: it has just one user, me. I can therefore design it in obnoxious and unintuitive ways to serve my peculiarities.

My two biggest goals for this pipeline:

1. Easy to publish writing
2. URLs must never change

To fulfill these goals I divide my content into two types, warm and cold.

## Warm Content

I previously wrote in my note-taking app during the draft phase, then manually migrated the draft into an HTML file committed to the twos.dev repository (more on why later). From then until publishing, I would work directly in the HTML. This method had three problems that worked against my first goal:

- Once transferred to HTML, context switching to that piece of writing became much more effortful
- Writing was slowed down by paper cuts like typing `<p>` tags, using `gqj` to wrap lines, etc.
- A paragraph may flow differently on a web page than in my notes, leading to refactors at the time of transfer

To solve these issues and work towards my first goal, I needed to close the gap between the draft phase and publishing.

I found [iA Writer](https://ia.net/writer) to replace my note-taking app for draft writing. It gets my thoughts out quickly, has helpful tools for e.g. reducing filler words, and looks nice to boot. Importantly, it stores writing as Markdown files.

### Shortcuts

iOS and macOS devices ship with a first-party app called [Shortcuts](https://apps.apple.com/us/app/shortcuts/id1462947752). If you’re from Android like I am, Shortcuts is like Tasker if it were built into the operating system and had buy-in from third parties. They're like event-driven shell scripts for the masses, and they run on iOS natively.

Using Shortcuts, I set up an automation that triggers when I exit iA Writer. This automation runs a series of shortcuts that:

1. Adds default frontmatter to iA Writer documents in the `published/` directory without any
2.

behave as an API gateway into iA Writer. I can set up a shortcut to run after I switch away from a specific app. It can inspect a directory in iCloud Drive (let’s make a “Published” directory), and do… something with it.

### Working Copy

[Working Copy](https://apps.apple.com/us/app/working-copy-git-client/id896694807) is a Git client for iOS that exposes Shortcuts hooks. Using it, we can build a shortcut that does what we need:

1. (Working Copy) Pull from `twos.dev` remote
2. (Files) Get contents of folder `Published`
3. Repeat with each item in `Contents of Folder`
   1. (Working Copy) Write `Repeat Item` to `src` in `twos.dev`
4. (Working Copy) Commit `twos.dev` with `Automatic commit by iA Writer sync job`
5. (Working Copy) Push `twos.dev` to remote

([see the shortcut](https://www.icloud.com/shortcuts/6580819cd24041a1b7e093cf6cbe5888))

My iPhone runs this each time I switch away from iA Writer, so changes are immediately published. If I switch to another iPhone later, it will inherit the behavior without any action on my part.

We now have glue between iCloud Drive and GitHub Actions.

### Markdown → HTML

This step is straightforward. Using Stripe’s Markdoc library (more on why later), a GitHub Actions workflow renders the Markdown documents into HTML at build time.

## Cold Content

We’ve got our Markdown pipeline set up, but not everything fits neatly into Markdown. For example, [Anonymously Autistic](autism.html) uses CSS gradients to render a spectrum. [The Pop-In](thepopin.html) uses media queries to show a different image in dark mode.

To handle these edge cases, we could run superset of Markdown, such as with templating or Markdoc, but that brings up its own issues:

- I’ll inevitably rewrite this infrastructure some years later, and won’t want to rehandle every edge case I’ve ever handled
- When it’s time to write code I want to use `$EDITOR`, not iA Writer
- twos.dev has a low volume of content—I don’t want to write new templating code for small features that may only be used once

For these not-quite-Markdown situations, then, **the right option is to hardcode**. Render the Markdown to HTML once, then edit the HTML by hand and commit it. Chances are good I’ll never touch it again.

#### Implementation

To grease the wheels of hardcoding, we’ll set up our GitHub Actions workflow to explicitly allow it:

1. Render `src/*.md` files → `dist/*.html`
2. Copy `src/*.html` files → `dist/`, overwriting existing files

#### Results

Allowing myself this escape hatch is freeing. It has three effects:

- I’m more encouraged to write interactive or otherwise bespoke components, e.g. to prove a point about button animation UX
- Twos.dev is uniformly structured by default, but I’m allowed case-by-case to break that uniformity when I see fit (e.g. a [CV](cv.html) has a unique need for bicolumnar content)
- TODO

## On URLs

I use my website as a test bed for interesting technology. As a side effect, content has been lost over the years as it becomes hard to migrate to newly-overhauled versions.

It so happens that by keeping my writing in these two formats, Markdown and raw HTML/CSS, I make it easy to accomplish my second goal: [cool URLs don’t change](https://www.w3.org/Provider/Style/URI). I’ll never migrate JavaScript frameworks and be too lazy to move things forward, or move to a database-backed writing system while being unable to simply `cp -r` these files into the `public` directory.

If I’ve done things right, this web page will be accessible at twos.dev/meta.html until the day I die—[and then some](death.html).

### Extra Credit

TODO: Frontmatter œ
