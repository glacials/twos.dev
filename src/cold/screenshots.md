---
date: 2024-02-04
filename: screenshots.html
type: draft
---

# Aligning macOS Screenshots With CSS

Today, I found myself writing a post for this site that included some screenshots from macOS.

Although my CSS skills are flawless^[ahahahahahahahahahahaha] the macOS screenshots looked a little timid:

![A screenshot of a twos.dev post embedding a macOS screenshot. The screenshot is noticeably narrower than the rest of the post above and below it.](/img/screenshots-1-dark.png)

_The screenshot is skinnier than the rest of the content._

This is because Apple adds a drop shadow to any window you screenshot with ⌘⇧4 → Space.

![A screenshot of the same content as before, but with developer tools enabled and showing that the screenshot is much wider than it looks to be, to account for the drop shadow embedded in it.](/img/screenshots-2-dark.png)

I like this behavior.
It makes screenshots look like they belong,
whether you're sharing them in a blog or an iMessage conversation.

However, it of course makes the CSS a bit harder.
How can we expand the actual window part of the screenshot to be the same width as the content,
and let the drop shadows bleed over into the margins?

## Option 1: Screenshot without Shadows

We can disable the screenshot drop shadow forever system-wide with a simple command:

```sh
defaults write com.apple.screencapture disable-shadow true
killall SystemUIServer # Or reboot
```

There are a few things I didn't like about this solution:

1. This would leave some older posts that use macOS screenshots behind,
   or require me to manually edit or recreate them.
2. I will 100% forget that I did this,
   and my next system will just reintroduce the problem.
   I could add it to my system setup scripts in [dot files](https://github.com/glacials/dotfiles),
   but that's just another thing that could break later.
3. As I said, I enjoy the drop shadow.

## Option 2: Manually Remove Shadows

I could opt to remove the shadows from each screenshot I take before adding it to a blog post,
either manually or by adding some workflow with `ffmpeg` or similar.

This has nearly the same problems as above.
It leaves older posts behind,
and adds another moving part that can break in the future,
and removes something I like from the screenshots.

## Option 3: Make macOS Screenshots Bigger

The option I'll go for is to write some CSS to make the macOS screenshots bigger on the page.
Here are our requirements:

1. Only macOS screenshots must be bigger; other images should be unaffected.
2. The borders of the screenshotted windows should align visually with the content above and below them.
3. In accordance with my [goals for twos.dev](meta.html) and [Winter](https://twos.dev/winter),
   it should be compatible with Markdown,
   not require HTML,
   and gracefully degrade into invisible non-syntax when viewed in other contexts.
   This ensures that if I change stylesheets,
   static site generators,
   and/or websites,
   I don't have to comb through all my previous posts and fix a bunch of broken-looking stuff.
   At that point it's okay if the shadow correction doesn't work anymore,
   but it's not okay if my posts have some weird [[REMOVE_SHADOWS(...)]] syntax or somesuch.

### Step 1: Move Widths to the Children

The first issue is that even when we get a way to make the image bigger,
it is still bound by the width of its container.

![A screenshot of the same content as before, but with developer tools enabled and showing that the post container spreads the entire height of the page, preventing anything inside it from sticking out on the left or right sides.](/img/screenshots-2-dark.png)

_The post `article` (blue) and its auto margins (orange)._

We need to break out of the `article` for macOS screenshots only.

We could stop the `article` just before the screenshot,
then display the image,
then start a new `article` afterwards;
however this violates the semantics of the article,
whose job is to contain one piece of content.

Even if it looks perfect,
accessibility tools may trip over it.
There's also no doubt that in 1 or 5 or 10 years future-me will come along and accidentally (and silently) break it.

Instead, let's kill the max-width on the container and move it to the children:

```diff
:root {
  --container-width: 33rem;
}

- article {
+ article > * {
  max-width: var(--container-width);
}
```

Our page should look the same,
but each immediate child of `article` is now responsible for its own width.
Depending on how you have your `article` aligned,
you may need to port some more code:

```diff
article > * {
+ margin: auto;
  max-width: var(--container-width);
}
```

You may opt to manually add a class to each element you place inside `article` instead.
I'm a forgetful person and I'm making a small website for myself only.
I like being able to build features once then forget about them,
so I reach for heuristics that will be good enough to free up my brain to worry about other things.

Now we can vary our widths per element.
Let's make sure only the right ones vary.

### Step 2: Pick Out Only macOS Screenshots

Our goal now is something similar to this:

```css
article > * {
  margin: auto;
  max-width: var(--container-width);
}

article > img.macos {
  max-width: 40rem;
}
```

Technically these images are in their own paragraphs,
so we need to complicate this a little to:

```css
article > * {
  margin: auto;
  max-width: var(--container-width);
}

article > p:has(img.macos) {
  max-width: 40rem;
}
```

But there are still two problems.

The first issue is rule 3 prevents us from tagging macOS images with their own class,
because we want to keep all that ugly HTML out of our futureproofed Markdown.
We need another way.

The second issue I'll come back to.

Markdown is pretty limited in the control we have over the `<img>` elements it produces,
with good reason.
Without any Markdown extensions,
and without mucking up our futureproof Markdown,
the only two attributes that get sent from Markdown to HTML are the `alt` and the `src` attributes.

We could hijack the `alt` attribute for our purposes:

```markdown
![macOS](/img/screenshot.png)
```

produces:

```html
<img alt="macos" src="/img/screenshot.png" />
```

which we can select using:

```css
article > p:has(img[alt="macos"]) {
  max-width: 40rem;
}
```

but this completely destroys the ability for assistive technologies like screen readers operate on our images,
which is not an acceptable solution.

The only other attribute we can impact from clean Markdown is `src`,
but we need that to find our image!

Luckily, CSS has a seldom-used selector for this:

```css
p:has(img[attr*="macos"]) {
  max-width: 40rem;
}
```

the [`[attr*=value]`](https://developer.mozilla.org/en-US/docs/Web/CSS/Attribute_selectors) selector matches any element whose `attr` attribute contains at least one occurrence of `value`.

All we have to do, then, is add a sort of tag to our macOS screenshots.
In fact, if you exclusively share screenshots of entire windows
(and never, say, subportions)
you could forego any special naming conventions and just match against what macOS names its screenshots by default:

```css
p:has(img[src*="Screenshot"]) {
  max-width: 40rem;
}
```

I share multiple types of screenshots, though.
[Winter](https://twos.dev/winter.html) also already supports a rudimentary tagging system for post images: ending a pre-extension filename in `-dark` or `-light` will select the correct image for the user's color scheme preference.

We can extend this simple system and say that any image with the word `macos` in it will be formatted as a macOS window screenshot:

```plain
my-screenshot-macos-dark.png
my-screenshot-macos-light.png
```

Our macOS screenshots can now be wider,
but we still have to account for the `<p>` around the image:

```css
p:has(img[src*="macos"]) {
  max-width: 40rem;
}
```

![A screenshot of the same content as earlier, but the macOS screenshot visually lines up with the text, mostly.](/img/screenshots-4-dark.png)

Now we've got our macOS screenshots _roughly_ lining up with the rest of the content.
That brings us to our last problem.

### Step 3: Figure Out How Much Wider macOS Screenshots Should Be

Now we've improved our example CSS to pick out only macOS screenshots:

```css
article > * {
  margin: auto;
  max-width: var(--container-width);
}

article > p:has(img.macos) {
  max-width: 40rem;
}
```

There's one remaining issue,
which is that `40rem` is just a number I pulled out of thin air.
It won't hold water if we use screenshots of different sizes or aspect ratios:

![A screenshot of the same dang content, but with a newer macOS screenshot below the older one, of a different shape and size. It doesn't line up with the text.](/img/screenshots-4-dark.png)

Let's try to be more precise with this.
What are we actually trying to achieve,
in precise terms?

We want the edges of the window in the screenshot to line up with some content.
The edges of the window can be found at the edges of the image,
inset by some amount.

With some quick pixel peeping in [Acorn](https://apps.apple.com/us/app/acorn-7/id1547371478?mt=12),
we can see that the macOS window screenshots use the same number of pixels for the shadow every time,
whether the window is small or large, tall or wide.

![A screenshot of a zoomed-in screenshot, with a select box dragged around the drop shadow portion. A status bar below says it is 112px wide.](/img/screenshots-8-dark.png)

_112px._

It's the same every time. 112px on the left and right sides.
But then I tried it on my external monitor:

![A screenshot of a zoomed-in screenshot, with a select box dragged around the drop shadow portion. A status bar below says it is 112px wide.](/img/screenshots-7-dark.png)

_56px._

It's different!
In fact it's _exactly_ half: 56px vs. 112px.

Looking at the image metadata,
Acorn reports the screenshot from this external monitor at 2560x1440 is 72 DPI,
while the original, from my MacBook Pro's 13" Retina display at 2560x1600,
is 144 DPI.
Exactly double.

There's something here.
If we can find a way to get the DPI of the image,
we can use some CSS math to figure out how much space to give it:

```css
article > p:has(img[src*="macos"]) {
  max-width: calc(var(--container-width) + calc(DPI_HERE * 2));
}
```

There's one last problem.
We're setting the right width for images at their full size,
but the page often doesn't display them at full size—they have to shrink to reasonably fit in it.

So the shadow, although it takes up `x` pixels in reality,
when the image is shrunk only takes up `y` pixels.
We need our container to also shrink by only `y` pixels.

The key insight is that the shadow is embedded in the image file itself,
so it scales proportionally with the image.
If a screenshot is 2000px wide and has a 112px shadow on each side,
the shadow is 5.6% of the image width.
When the image scales down to fit the page,
the shadow also scales down by that same percentage.

So instead of adding a fixed pixel value,
we can scale the container width by the percentage that the shadow occupies.
Since we know the shadow is 5.6% on either side,
we can scale the container width up by `5.6% * 2 = 11.2%`.

```css
article > * {
  margin: auto;
  max-width: var(--container-width);
}

article > p:has(img[src*="macos"]) {
  max-width: calc(var(--container-width) * 1.112);
}
```

This works at any display size because both the image and its embedded shadow scale together.

![A screenshot of the final result: macOS screenshots align perfectly with the text content, with shadows bleeding into the margins.](/img/screenshots-5-dark.png)

_Now they line up perfectly, regardless of the screenshot size or page zoom level._
