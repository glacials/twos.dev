# twos.dev

`twos.dev` (pronounced "tooz dot dev" or "two esses dot dev" if you like) is my
portfolio, whose name is a joke based on how I have to describe my last name to
people. This website is basically an extended version of my resume.

I originally wrote this in Vue single-file components, with a lot of CSS and JavaScript and it looked great! Until I tried to work on it on a bad internet connection and realized I couldn't even get far enough through a `yarn install` to work on it. A simple, single, static web page and I was thwarted by JavaScript bloat. Out of frustration I rewrote it as an HTML file with no dependencies and I'm very happy I did that. You can still see the old version at [twos.dev/fancy][fancy], but I don't keep the information there up-to-date.

[fancy]: https://twos.dev/fancy

## Development

### Dependencies

- [`yarn`][yarn]

[yarn]: https://github.com/yarnpkg/yarn

### Build

```sh
yarn install
```

### Run

```sh
npm rebuild node-sass
yarn serve
```
