#!/usr/bin/env node
// Synchronize iA Writer files in iCloud to the twos.dev Git repository.

import replaceExt from 'replace-ext'
import fs from 'fs-extra'
import simpleGit from 'simple-git'

const ia = `${os.homedir()}/Library/Mobile\ Documents/27N4MQEA55\~pro\~writer/Documents/Published`;
const src = './src'

(async () => {
  try {
    const git = simpleGit()
    await git.pull()

    // Change any .txt or no-extension iA file to .md
    await Promise.all((await fs.readdir(ia)).map(filename => {
      if (filename.substring(filename.length - ".txt".length) == ".txt" || filename.indexOf('.') < 0) {
        return fs.rename(
          `${ia}/${filename}`,
          `${ia}/${replaceExt(`${ia}/${filename}`, '.md')}`,
        )
      }

      return Promise.resolve()
    }))

    fs.copySync(ia, src)

    await git.add((await fs.readdir(ia)).map(filename => `${src}/${filename}`))
    await git.commit("Auto-sync job from iA Writer")
    await git.push()
  } catch (e) {
    console.error(e)
    process.exit(1)
  }
})()
