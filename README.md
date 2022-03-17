# twos.dev

This is the source for my personal website. I post thoughts, hobbies, and other random
things here.

"twos" (pronounce "two esses" or "tooz") is a play on how I have to describe my last
name to people.

I originally wrote this in Vue single-file components, with plenty of CSS and JavaScript
and it looked really fancy! Then I tried to work on it on a bad internet connection and
couldn't even get through a `yarn install` -- thwarted by JS bloat for what ultimately
was a simple static website.

Out of frustration I rewrote it in raw HTML with no dependencies and I'm very happy I
did that. It's not as fancy but it's darned easy to work on, which is a great thing to
optimize for a personal website.

## First run

### Dependencies

No dependencies.

### Starting a dev "server"

```sh
# macOS
open src/index.html

# Linux
xdg-open src/index.html
```
