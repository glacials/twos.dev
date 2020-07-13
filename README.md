# twos.dev
`twos.dev` (pronounced "tooz dot dev" or "two esses dot dev" if you like) is my
portfolio, whose name is a joke based on how I have to describe my last name to
people. This website is basically an extended version of my resume.

The site is written in Vue single-file components, for no other reason than
that's what I was learning last time I rewrote it.

## Development

### Dependencies
- [`yarn`][yarn]
- [`yb`][yb]

[yb]: https://github.com/yourbase/yb
[yarn]: https://github.com/yarnpkg/yarn

### Build and test
```
yb build
```

### Run
```
yb exec
```

### Deploy
Deploys to GitHub Pages happen automatically on push to `master` via
[YourBase][yourbase], but with a [GitHub `repo` token][token] in `GITHUB_TOKEN`
you can initiate one manually with
```sh
yb build deploy
```

[yourbase]: https://yourbase.io
[token]: https://github.com/settings/tokens
