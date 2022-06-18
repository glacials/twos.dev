---
date: 2022-06-17
---

# Dashes over the Years

There are three common dashes in English written language:

- Hyphen (-)
- En dash (–)
- Em dash (—)

The en dash and em dash are so named after their widths, which are equal to the widths of `N` and `M` respectively. "Hyphen" comes from ὑφέν (huphén) meaning "together". [^https://en.wiktionary.org/wiki/en_dash][^https://en.wiktionary.org/wiki/em_dash][^https://en.wiktionary.org/wiki/hyphen].

## Usage

The hyphen is used to hyphenate words or otherwise join parts of a whole, like father-in-law.

The en dash is used to create ranges, like 2020–2030.

The em dash is used like a parenthetical—to separate one thought from another—when the thought needs to stand out (more than a typical parenthetical).

### A Note About the Em Dash

Em dashes tend to be a trap when a writer first learns to jive with them. They cause a piece of text to flow like speech, without such sudden endings as periods. For those whose writings spring from internal monologues, this seems to be a blessing, a way to avoid molding thoughts into a comparatively non-fluid format. This is true, but it can backfire if overused.

For readers, long flowing thoughts connected by em dashes take mental effort to keep up with. The reader must keep the grammar and meaning of the current sentence in an internal buffer until the very end, which -- for example, if the author decides to describe something tangential or dive into an example to assist with understanding, hoping their explanation can let the rest of the sentence off the hook but ultimately making the reader do more work to understand it -- can be hard to do. Additionally, the em dash is a large glyph and its impact to the sentence structure is even larger. Its use is obvious even through squinting eyes, so its overuse becomes a stain on the page even before reading it.

## Hidden Complexities

Within each usage pattern above are small issues.

### Spaces Overpower Hyphens

When one side of a hyphen is a compound word like "ice cream break", its spaces visually overpower the hyphen: `pre-ice cream break` looks like `[pre-ice] cream break`, not `pre-[ice cream break]`.

One solution replaces the spaces with hyphens (`pre-ice-cream-break`), but this removes the ability to distinguish the original parts. Another removes the hyphens altogether (`pre ice cream break`), which has a similar problem but may be better for specific situations.

Those using LaTeX can replace the spaces with [thin spaces](https://en.wikipedia.org/wiki/Thin_space) (\LaTeX: `\thinspace`) to de-emphasize their separation:

```
pre-ice\thinspace{}cream\thinspace{}break
```

Personally, I've found it best to restructure the sentence to avoid the problem:

```
before the ice cream break
```

It's crude, but for 90% of cases there's no loss of meaning.

### Line Breaks Interfering with Dashes

If a line soft-wraps immediately before or after a dash, what should its behavior be? It's accepted that a hyphen can appear as the last character on a line, whether it was already there or is being used to join two parts of a non-hyphenated word. En and em dashes are a little harder to parse.

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

Generally it's recommended to either give them the same treatment as hyphens or to treat their left and right neighbors as one inseparable "word", as backwards as it may sound considering our treatment of hyphens.

### Justified Text Messes with Proportions

Whatever dash one is using, when using justified text (where spaces are elongated to make both the left and right edges of a page straight vertical lines) elongated spaces can overpower dashes where non-elongated spaces would not have.

For example, see this simulation of left-aligned text:

```plain
v left page edge             v right page edge
|                            |
I said they were—wolf
playdates don't last forever!
```

versus this simulation of justified text:

```plain
v left page edge             v right page edge
|                            |
I    said    they    were—wolf
playdates don't last forever!
```

wherein "were-wolf" reads hyphenated because the spaces are proportionally larger than the em dash, making it appear as a hyphen.

One solution surrounds parenthetical em dashes with a space on either side, causing its total width to scale up at twice the rate of spaces:

```plain
v left page edge             v right page edge
|                            |
I   said   they  were  —  wolf
playdates don't last forever!
```


### Input Difficulties

Standard US keyboards support only the hyphen without more advanced knowledge, so many en and em dash users simply compose them from hyphens and spaces.

#### The Hyphen-Powered En Dash

The en dash looks similar enough to the hyphen in most fonts that people often settle for a single hyphen. \LaTeX requires two (`--`) for the en dash, however.

#### The Hyphen-Powered Em Dash

A common replacement for the em dash is a space on either end of two hyphens (`  --  `) or simply two hyphens (`--`). Because \LaTeX reserves two hyphens (`--`) for the en dash, it requires three (`---`) for the em.

### Output Difficulties

At the time of writing, twos.dev renders in a monospace font. Monospace fonts have only a single glyph used for all hyphen types, so it becomes impossible to differentiate. Currently, I use a bespoke text preprocessor  [insert link once pushed] to replace em dashes with a multi-glyph hyphen-based em dash, but I'm not satisfied with this solution. Because I use monospace for aesthetic only reasons, I'm curious to find or build a monospace-with-exceptions font that allows for a small number of double-width glyphs.

## Debate

Every usage rule above is in some level of turmoil [better word?] among writers, both within the rules and between them.

## Esoteric Dashes

The **swung dash** (⁓) is an elongated tilde used to stand in for a word being defined in a dictionary. [^http://wordnetweb.princeton.edu/perl/webwn?s=swung+dash&sub=Search+WordNet&o2=&o0=1&o8=1&o1=1&o7=&o5=&o9=&o6=&o3=&o4=&h=000000000000]

> boot (n)
>
> ex: Let me put on my other ⁓.