---
filename: meta.html
date: 2022-05
---

# The Bespoke Twos Deployment Pipeline
The Twos CD pipeline has a unique advantage against anything I’ve ever worked on: it has just one user, me, so can be designed in what would otherwise be obnoxious and unintuitive ways to serve my peculiarities.

My two biggest goals for this pipeline:

1. Encourage myself to write
2. URLs must never change

To fulfill these goals, the Twos CD pipeline allows for two content types:

- Warm content
- Cold content

## On Writing

I use [iA Writer](https://ia.net/writer) to write drafts. It gets my thoughts out quickly, without mucking about in formatting or DOM. If I were writing for a personal diary, things would end there.

Because I like to publish, it's a problem that these drafts too often stay drafts.

Luckily, iA Writer stores documents as plaintext Markdown files. Normally, shipping these files to a Markdown parser would be a sufficient general solution.

But there are two peculiarities I’m allowing myself to indulge.

### Automatic Publishing

Because I have trouble pulling the publish trigger, I’ve decided upfront that I need to make that trigger 10x easier to pull. For this case, a good path towards that is to pull the publishing process out of Git and into iA Writer itself, to allow easy publishing from even my phone.

### Mechanics

iA Writer on macOS has some automation tools, but not a lot. It has the ability to store documents in services like Dropbox and Google Drive, but for [unrelated peculiarities](apple.html) I limit myself to iCloud Drive for those purposes at the moment. Unfortunately, iCloud Drive has no native API or other way to reach files from within GitHub Actions. (TODO: Perhaps there’s a way to sign into an Apple ID on a GitHub Actions Mac? Perhaps not due to 2FA, or maybe an app-specific password would work?)

I keep infrastructure off my own devices when I can—I wipe my drives a lot—but it seems there is no way to avoid it if I need to ship files from iCloud Drive to GitHub Actions.

What I can do is automate this in a way that is resistant to my drive-wiping habits. And so far, the most resistant device I own to this, the only one I’m comfortable directly transferring to a replacement device using its migration tool, is my iPhone.

#### Shortcuts

Starting with iOS 13, all iPhones ship with an Apple app called [Shortcuts](https://apps.apple.com/us/app/shortcuts/id1462947752). If you’re from Android like I am, Shortcuts is like Tasker, if it were built into the operating system and had buy-in from every third-party app. For example, I use a shortcut to turn on my coffee machine every morning when I take my phone off its charger.

Shortcuts can therefore behave as an API gateway into iCloud Drive. I can set up a shortcut to run like a cron job, inspect a specific directory in iCloud Drive (let’s make a “Published” directory), and do… something with it.

#### Working Copy

[Working Copy](https://apps.apple.com/us/app/working-copy-git-client/id896694807) is a Git client for iOS that exposes Shortcuts hooks. Using it, we can build a shortcut that does what we need:

1. (Working Copy) Pull from `twos.dev` remote
2. (Files) Get contents of folder `Published`
3. Repeat with each item in `Contents of Folder`
    1. (Working Copy) Write `Repeat Item` to `src` in `twos.dev`
4. (Working Copy) Commit `twos.dev` with `Automatic commit by iA Writer sync job`
5. (Working Copy) Push `twos.dev` to remote

([see the shortcut](https://www.icloud.com/shortcuts/6580819cd24041a1b7e093cf6cbe5888))

My iPhone runs this once daily at sunrise. If I switch to another iPhone later, it will inherit the behavior without any action on my part.

We now have glue between iCloud Drive and GitHub Actions.

#### Markdown → HTML

This step is straightforward. Using Stripe’s Markdoc library (more on why later), a GitHub Actions workflow renders the Markdown documents into HTML at build time.

### Custom Content

We’ve got our Markdown pipeline set up, but not everything fits neatly into Markdown. For example, [Anonymously Autistic](autism.html) uses CSS gradients to render a spectrum. [The Pop-In](thepopin.html) uses media queries to show a different image in dark mode.

To handle these edge cases, we could run superset of Markdown, such as with templating or Markdoc, but that brings up its own issues:

- I’ll inevitably rewrite this infrastructure some years later, and won’t want to rehandle every edge case I’ve ever handled
- When it’s time to write code I want to use `$EDITOR`, not iA Writer
- twos.dev has a low volume of content—I don’t want to write new templating code for small features that may only be used once

For these not-quite-Markdown situations, then, **the right option  is to hardcode**. Render the Markdown to HTML once, then edit the HTML by hand and commit it. Chances are good I’ll never touch it again.

#### Implementation

To grease the wheels of hardcoding, we’ll set up our GitHub Actions workflow to explicitly allow it:

1. Render `src/*.md` files → `dist/*.html`
2. Copy `src/*.html` files → `dist/`, overwriting existing files

#### Results

Allowing myself this escape hatch is freeing. It has three effects:

- I’m more encouraged to write interactive or otherwise bespoke components, e.g. to prove a point about button animation UX
- Twos.dev is uniformly structured by default, but I’m allowed  case-by-case to break that uniformity when I see fit (e.g. a [CV](cv.html) has a unique need for bicolumnar content)
- TODO

## On URLs

I use my website as a test bed for interesting technology. As a side effect, content has been lost over the years as it becomes hard to migrate to newly-overhauled versions.

It so happens that by keeping my writing in these two formats, Markdown and raw HTML/CSS, I make it easy to accomplish my second goal: [cool URLs don’t change](https://www.w3.org/Provider/Style/URI). I’ll never migrate JavaScript frameworks and be too lazy to move things forward, or move to a database-backed writing system while being unable to simply `cp -r` these files into the `public` directory.

If I’ve done things right, this web page will be accessible at twos.dev/meta.html until the day I die—[and then some](death.html).

### Extra Credit

TODO: Frontmatter 