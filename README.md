# twos.dev
`twos.dev` (pronounced "tooz dot dev" or "two esses dot dev" if you like) is my portfolio, whose name is a joke based on
how I have to describe my last name to people. This website is basically an extended version of my resume.

The site is written in Vue single-file components, for no other reason than that's what I was learning last time I
rewrote it.

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
Deploys to GitHub Pages happen automatically on push to `master` via [YourBase][1], but you can initiate one manually
with
```
yb build deploy
```
For that to work you will have to set up Git to use a [GitHub token][2] with `repo` permissions for authentication by
running something like this first:
```
git remote set-url origin https://glacials:INSERT_TOKEN_HERE@github.com/glacials/twos.dev.git
```

[1]: https://yourbase.io
[2]: https://github.com/settings/tokens
