---
date: 2022-06-18
updated: 2023-02-20
filename: meta.html
preview: How does this website even work?
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

1. Writing HTML is a rough context switch (what was scribbling drafts on my phone became sitting at my computer in `$EDITOR`)
2. Writing was slowed down by cruft like `<p>` and `<li>`
3. The reading experience can flow differently on a web page than in a notes app, leading to large rewrites after transferring

To solve these issues and work towards my first goal, I needed tooling in the space between amorphous notes and strict HTML.

I first focused on changing my drafting app to something with more robust exporting
tools. I've jumped between [iA Writer](https://ia.net/writer) and
[Obsidian](https://obsidian.md), both of which let me get my thoughts out quickly, each
with their own strengths. Importantly, they store writing as Markdown. This solves `1`
and `2` above.

To solve `3`, we need to automate getting those rough drafts published to an unlisted
twos.dev page every time they change. This means we need programmatic access to the
drafts, but neither app has a third-party web API. Luckily, the editors can store the
files in the location of your choosing, including some natively supported cloud storage
services. We need to get at those files.

There are several ways to do this. The simplest would be to point your writing "vault"
directly at a cloud storage service with an API like Google Drive, then use something
like a cron-based GitHub Actions workflow to ensure your writing always gets published.

Since I'm still undergoing [my Apple experiment](apple.html), I've set myself a
constraint of only using iCloud Drive for cloud storage, which doesn't have an API. That
leaves two options.

### Shortcuts

iOS and macOS ship with
[Shortcuts](https://apps.apple.com/us/app/shortcuts/id1462947752), a no-code
event-driven automation app. Using Shortcuts, I set up an automation that triggers when
I switch away from the my writing app, and additionally on a cron. The automation adds a
1-2 line [Winter frontmatter](winter.html#frontmatter) section to each document if
needed, then pushes it to the `src/warm` directory in the twos.dev Git repository by
invoking [Working Copy](https://workingcopyapp.com). No interruption to UX; this happens
in the background.

I never thought a phone would be an integral part of my CI/CD pipeline, but here we are.
It actually hummed along in the background for several months until an iOS update
changed some minor behavior somewhere that broke it. I didn't have the motivation to fix
it because the debugging experience for Shortcuts isn't great, but otherwise this was a
surprisingly effective solution and it always triggered exactly when it was needed.

That leads us to our second option.

### `launchd`

There's no network API into iCloud Drive, but there is a local one: the filesystem. If
we control a Mac that will be online enough to satisfy our synchronization needs, which
are not necessarily 24/7, we can snatch the files from there.

Enter
[launchd](https://developer.apple.com/library/archive/documentation/MacOSX/Conceptual/BPSystemStartup/Chapters/ScheduledJobs.html#//apple_ref/doc/uid/10000172i-CH1-SW2),
a hidden gem of macOS. It's a few things, but for our purposes it's basically cron
purpose-built for laptops. It intelligently handles sleep, network loss, and ([charge
state](https://support.apple.com/guide/mac-help/what-is-power-nap-mh40773/mac)). The
specifics aren't important; suffice it to say It Just Works.

As long as we install this on a machine or machines that get enough use—or that are
[Caffeinated](https://caffeinated.app) enough—this will be a good enough pipeline. We
have our draft synchronization.

### Preprocessing

On push, GitHub Actions builds the Markdown file into HTML with all the [right preprocessing](https://github.com/glacials/twos.dev/blob/c59cc1a/winter/document.go#L408-L436) required to make it look like twos.dev. I use a small subset of [`html/template`](https://pkg.go.dev/html/template) to allow myself to put things like MP4 screen recordings and dark-mode-aware images into Markdown. This is also where the frontmatter from before is stripped and parsed into metadata inserted elsewhere on the page.

Technically at this point the page is now published, just not linked to from anywhere. I can see how it looks in the context of twos.dev, and send the link to friends who have offered to review.

## Cold Content

There are 2-3 needs warm content doesn't cover: first, this odd pipeline that uses my phone for CI is not something I have high confidence in, so I'd like to limit its blast radius to content that benefits from it—content I'm actively working on—and keep safe the innocent bystanders that are historical content.

Second, once in a while I need some wacky one-off piece of code for a single post, like the CSS-only spectrum in [Anonymously Autistic](autism.html) or the variable-width font requirements of [Advanced Dashes](dashes.html). iA Writer is great for prose, but when it's time to write code I need to be back in `$EDITOR`.

Cold content exists to serve these needs. At the most basic level, cold content is just warm content that has been moved to a directory outside the reach of the warm pipeline. This is content I only touch with `$EDITOR`, like a normal code repository.

To fully turn warm content cold, there is a small amount of bookkeeping.

### Winter

To handle this bookkeeping along with the other build pipeline steps, I wrote the [Winter](https://twos.dev/winter) CLI. Winter can handle the process of converting warm content to cold like so:

```sh
winter freeze DOCUMENT
```

which more or less performs a `git mv src/{cold,warm}/DOCUMENT.md` plus chips.

Winter also handles all the static generation. Writing my own generator was an
important step in building longevity into twos.dev because I often switch tools
when I discover a small feature I want but don't have (put simply: "oo shiny")
and that switching of tools is what eventually erodes the content.

Many static site generators are open source, but I'll often prefer to switch
tools than learn how to contribute if the barrier to entry contribution is
high, such as a language I don't like writing or an overly complex architecture.

Writing the generator myself guarantees that barrier is minimal. Writing it in
Go is a guarantee the code will not rot quickly or fall into dependency hell
after a few years without attention.

### Post-freeze

After freezing my content I may also take one final step of permanently
rendering it to HTML if I've built more than a little code into it. Sure I can
always embed HTML in Markdown, but taking this extra step gives me all my editor
doodads and hardens the file against future changes to preprocessing.

```sh
winter build
mv dist/DOCUMENT.html src/cold/DOCUMENT.html
git rm src/cold/DOCUMENT.md
git add src/cold/DOCUMENT.html
```

This brings back the weight of editing prose in HTML, but these cases are the minority and at this point most of the writing is done.

#### Results

Allowing myself this escape hatch is freeing. I’m more encouraged to write interactive or otherwise bespoke components, and twos.dev becomes consistent by default but in a way I can break when I need (e.g. for a [CV](cv.html)'s bicolumnar layout).

## On URLs Not Changing

I use my website as a test bed for interesting technology. As a side effect, content has been lost over the years as it becomes hard to migrate to overhauled versions.

It so happens that by keeping my writing in Markdown or raw HTML/CSS, I make it easy to accomplish my second goal: [cool URLs don’t change](https://www.w3.org/Provider/Style/URI). I've removed the possibility that changing frontend or backend frameworks will leave my documents behind; Markdown is not going anywhere, and all the HTML files need to look decent is a header and footer. In the worst case scenario I'll simply `cp -r dist` to my next crazy system's web root.

If I’ve done things right, this web page will be accessible at twos.dev/meta.html for at least as long as I live.
