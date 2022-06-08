# `src/warm`

This directory is a sync target from a directory in my personal [iA
Writer](https://ia.net/writer) library. This allows me to write with low
friction from any platform, without worrying about context switching into a
proper development environment. The build pipeline automatically renders
Markdown in this directory to HTML at build time.

If a document becomes too complex to be represented by Markdown, or if I'm
pretty sure I'll never touch it again, I render it to HTML one final time and
move it [`src/cold`](../cold).
