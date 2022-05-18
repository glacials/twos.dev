---
filename: meta.html
date: 2022-05
---

# My Phone is Part of this Website’s CI/CD
I have two goals for this website:

1. It should be easy to write things
2. It should be hard to break things

## It Should Be Easy to Write Things

I use [iA Writer](https://ia.net/writer) for writing drafts. It lets me get thoughts out quickly, without mucking about in code. However, I’d like not to have to do extra work to move them from there to the internet.

That’s a straightforward task on its own, but being easy to write sometimes means having the ability to drop into code. For example, [Anonymously Autistic](autism.html) uses CSS gradients to render a spectrum. [The Pop-In](thepopin.html) uses media queries to show a different image to readers in dark mode. I need a way to keep things in iA Writer by default, but allow moving into my code editor if I need to get fancy.

## It Should Be Hard to Break Things

I use my website as a test bed for interesting technology. As a side effect, content has been lost over the years as it becomes hard to migrate to newly-overhauled versions.

I’d like for this never to happen again. [Cool URLs don’t change](https://www.w3.org/Provider/Style/URI). If I’ve done things right, this web page will be accessible at twos.dev/meta.html until the day I die [and beyond](death.html).

## Implementation 

Frameworks are out of the question. I do not pretend I can stop my future self from another overhaul, or from dropping content that’s hard to migrate through one. Instead I must prepare for it.

Abstracting content far from structure (e.g. Markdown) does not satisfy my requirement for dropping into code once in a while.

### Biphasic Content

To accommodate my needs, this website has two modes of content.

#### Warm Content

Warm content is easy to write. I write new, in-progress, or simple content in iA Writer. iA Writer saves this content to iCloud in a directory of its own, natively in Markdown.

On my iPhone, I have a Shortcuts automation that triggers daily. This automation copies these files from iCloud into [Working Copy](https://workingcopyapp.com/)’s local clone of this website, which adds, commits, and pushes them.

A GitHub Actions workflow triggers that builds the Markdown files into HTML files using [Markdoc](https://markdoc.io/), and places it between a header and footer using a few lines of Node and bespoke templating to get, for example, the `<title>` tag right.

##### Aside: Why iPhone

As I get older, I worry less about things working now and more about them working in 10 years. I can set up a cron job on [a machine that’s persistently on](apple.html#iMac), but when I inevitably replace it I’m going to forget—or forget how—to set it up again. I may not even replace that machine, but instead remove it. Discovering that missing cron is a ripe opportunity for me to overhaul things again, which is not something I want to encourage.

But the process of migrating from one iPhone to another is seamless. The shortcuts and automations come along automatically. And I will always have a phone.

There’s something magical about having your “server at home” be in your pocket.

#### Cold Content

Cold content is hard to break. I render old, complete, and complex content into HTML once, sans header and footer, and commit the result. 

When a piece of cold content and a piece of warm content conflict, the cold content wins—I cannot update the content with iA Writer anymore, because it can’t do the job sufficiently for complex content.