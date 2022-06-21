---
type: draft
date: 2022-06-17
---

# Dashes

There are three common dashes in written English:

- Hyphen : -
- En dash: –
- Em dash: —

## Usage

The hyphen is used to join words, prefixes, and/or suffixes together, like father-in-law.

The en dash is used to create ranges, like 2020–2030.

The em dash can be used as a more powerful parenthetical, a replacement for a colon, or to represent an interruption of a speaker.

## One Size Doesn't Fit All

A myth states the em dash is so named because it is "one `M` wide", and the en dash "one `N` wide". This is... almost true. The confusing correction is that an em dash is so named because it is `1em` wide. Incidentally, the en is defined as half of one em; much of my focus will be on the em.

### The Em
    
To talk about the em dash necessitates talk about the [em](https://en.wikipedia.org/wiki/Em_(typography)). The em (pronounced [/ɛm/](https://en.wiktionary.org/wiki/Appendix:English_pronunciation)) is an invisible bounding box every character is designed within[^"Bounding box" is strong; font designers are free to design outside it.] whose edge length is equal to the font height. Font height changes with font size, so the size of the em changes with it: `1em` is always one *character-heighth* long. It is a unitless measurement that allows us to speak in relatives.

When dealing with font sizes, it follows, speaking in ems is functionally equivalent to speaking in percentages. A font size in ems is a ratio between the desired font size and the current font size. The CSS declaration `font-size: 1em` is equivalent to `font-size: 100%`, while `1.5em` is equivalent to `150%`, and so on.

The em's eccentric offering is that it can be used outside the context of font sizes. It can set margins, widths, even blur radii based on the font size of the containing context. It can scale an image up and down with the text around it, whether the text's changing size is a product of assistive technologies or simply of a component being reused in multiple contexts.

The em came about in the times of the printing press, when every `M` type was large and perfectly square, thus convenient to use as a proportional measurement for a given font and size. When the `M` type eventually rebelled against being square, the modern em shifted its definition to the remaining constant, the height. Practically this was a non-change, but the updated definition allowed the `M` and the em to be decoupled—in the same way the SI second was decoupled from the Earth by defining it in terms of [radioactive decay](https://en.wikipedia.org/wiki/Caesium_standard).

#### The Em Dash

With the tailwind given by the definition of the em, the em dash is easily defined: a dash that is `1em` long. In a way, the em dash took up the mantle of the square `M` type. Because the em dash is `1em` long, and the em is defined as the height of the font, the em dash's bounding box is an exact square—in any font, at any size.

#### The Root Em

The `rem` or **root em** is a twist on the em. It is relative not to the parent context's font size, but to the root context's; in a web page, this is `<html>`. This allows one to escape relativity hell without reverting to absolutes:

```plain
html: 16pt              (default     = 16pt)
  → body: 1.5em         (16pt * 1.5  = 24pt)
    → .container: .75em (24pt *  .75 = 18pt)
      → h1: 4em         (18pt * 4    = 72pt)
      → p: 1em          (18pt * 1    = 18pt)
    → footer: 0.5rem    (16pt * 0.5  =  8pt)
      
key:
context: requested size (calculation = final size)
  → child element
```

### The Point
The point has had a tumultuous and unstandardized history both physically and digitally, but today the world has predominantly settled on `1/72in` per point. It's no coincidence if 72 dots per inch rings a bell for screen density. `1pt` = 1 dot.

That said, super high resolution displays have caused organizations to invent their own ways of scaling on behalf of softwares that expect 72 dpi displays, like Android's [density-independent pixels](https://developer.android.com/training/multiscreen/screendensities#TaskUseDP) or `dp`, and Windows's [effective pixels](https://docs.microsoft.com/en-us/windows/apps/design/layout/screen-sizes-and-breakpoints-for-responsive-design#effective-pixels-and-scale-factor) or `ep`/`epx`. The tumultuous history, as it happens, hasn't ended.

## Advanced Usage

Advanced usage of dashes unfortunately revolves around avoiding common issues.

### Em Dash Overuse

As a personal opinion, the em dash tends to be a trap to a new writer. It is sympathetic to stream-of-consciousness writing, wherein the writer is in "append-only" mode—thinking not of the structure of the sentence but of what new words may clarify it when tacked on. For people whose writings spring from internal monologues this seems a blessing: a way to avoid adding structure to organic thoughts.

Organic thoughts have their place in writing, but it turns out spending time perfecting sentence structure pays dividends to the reader. Long flowing thoughts connected by em dashes take mental effort to keep up with. Because the pre-dash thought may continue post-dash, the reader must keep its grammar and intention in an internal buffer, which—for example, if the author decides to describe something tangential or dive into an example to assist with understanding, hoping their explanation can let the rest of the sentence off the hook but ultimately making the reader do more work to understand it—can be hard to do.

Separately, the em dash is a large glyph and hard to go unnoticed. Its use is obvious even in peripheral vision, so its overuse becomes a stain on the page before reading it.

### Belittled Hyphens

When one side of a hyphen is a compound word like "ice cream break" the spaces visually overpower the hyphen: `pre-ice cream break` looks like `[pre-ice] cream break`, not `pre-[ice cream break]`.

One solution replaces spaces with hyphens (`pre-ice-cream-break`) or vice versa (`pre ice cream break`), becoming partially lossy and potentially ambiguous in the process.

\LaTeX users can replace the spaces with [thin spaces](https://en.wikipedia.org/wiki/Thin_space) (\LaTeX: `\thinspace`) to de-emphasize them relative to hyphens:

```
pre-ice\thinspace{}cream\thinspace{}break
```

Personally, I've found it best to restructure the sentence to avoid the problem:

```
before the ice cream break
```

It's crude, but for 90% of cases there's no loss of meaning.

### Broken Dashes

If a line soft-wraps immediately before or after a dash, what should its behavior be? It's accepted that a hyphen can appear as the last character on a line, whether it was already there or is being introduced to join a word split across lines. En and em dashes are a little harder to parse.

```plain
v left page edge             v right page edge
|                            |
I am looking at dates about 20
–30 days out.
```
```plain
|                            |
I'm looking at dates about 20–
30 days out.
```
```plain
|                            |
I am pleased they thought that
—however odd it may be.
```
```plain
|                            |
I'm pleased they thought that—
however odd it may be.
```

Strangely, the hyphen is the only dash that has an option here: the [non-breaking hyphen](https://en.wikipedia.org/wiki/Hyphen#Non-breaking_hyphens) or "hard hyphen" is a hyphen glyph that sends a signal to the word wrapper that a wrap must not happen after it.

### Drowned Dashes

When using [justified text](https://en.wikipedia.org/wiki/Typographic_alignment#Justified) elongated spaces can overpower dashes where non-elongated spaces would not have.

For example, see this simulation of left-aligned text:

```plain
v left page edge             v right page edge
|                            |
I said they were—wolf
playdates don't last forever!
```

And this simulation of justified text:

```plain
|                            |
I    said    they    were—wolf
playdates don't last forever!
```

In the justified simulation, "were—wolf" reads like the mythical creature "were-wolf" because the spaces are proportionally larger than the em dash.

One solution surrounds each parenthetical em dash with a space on either side, causing its total width to scale up at twice the rate of spaces:

```plain
|                            |
I   said  they  were  —  wolf
playdates don't last forever!
```

### Input Difficulties

Standard US keyboards support only the hyphen without more advanced knowledge, so many en and em dash users simply compose them from hyphens and spaces.

No one should ever fault that, but for those who are willing to commit something new to memory:

- macOS
	- En dash: `⌥ + -`
	- Em dash: `⌥ + Shift + -`
- Windows
	- En dash: `Alt + 0150`
	- Em dash: `Alt + 0151`

For all others, read on.

#### The Hyphen-Powered En Dash

The en dash looks similar enough to the hyphen in most fonts that people often settle for a single hyphen: `2020-2030`. \LaTeX requires two (`--`) for the en dash.

#### The Hyphen-Powered Em Dash

A common replacement for the em dash is a space on either end of two hyphens (`  --  `). Because \LaTeX reserves two hyphens for the en dash, it requires three (`---`) for the em.

### Output Difficulties

At the time of writing, twos.dev renders in a monospaced font for aesthetic reasons. Unfortunately, the glyphs for many dashes are virtually indistinguishable in monospaced fonts. To make the differences visible, I wrote [a bespoke preprocessing step](https://github.com/glacials/twos.dev/blob/main/cmd/transform/lengthen_dashes.go) to switch twos.dev's dashes—and dashes only—to variable width fonts.

I'd love for there to be a "monospace with exceptions" font that takes this chore out of my hands.

## Bonus Round: Esoteric Dashes

The hyphen, en dash, and em dash get all the love, but behind the scenes are a silent majority of dashes that don't often get to see the light of day.

The **swung dash** (⁓) is an elongated tilde used to stand in for a word being defined in a dictionary.[^http://wordnetweb.princeton.edu/perl/webwn?s=swung+dash\]

> boot (n)
>
> ex: Let me put on my other ⁓.

The **horizontal bar** is a way to introduce quotations.[^ ] Confusingly, its length is almost always identical to the em dash.

> ― O Miss Douce! Miss Kennedy protested. You horrid thing!
> 
> *James Joyce's Ulysses p. 335*

The **hyphen bullet** is a hyphen to be used in place of a bullet point.

```plain
- This is a hyphen
⁃ This is a hyphen bullet
```

The **figure dash** is a variant of the en dash having the same width as digits (which are uniformly wide in most fonts). It is meant for phone numbers and other numeric contexts where columnar alignment is required or pleasing.

```plain
Figure dash: 867‒5309
En dash:     867–5309
Hyphen:      867-5309
```

Lastly, my favorite: the **soft hyphen** is a zero-width, invisible character that (opposite to the hard hyphen) denotes a place the word wrapper is *welcome* to wrap. This can be used in the middle of a compound word or long line of inert code to provide a cleaner wrap.

```plain
v left page edge            v right page edge
|                           |
No soft hyphen:
Supercalifragilisticexpialid-
ocious

Soft hyphen after 'istic':
Supercalifragilistic-
expialidocious
```

How dashing!