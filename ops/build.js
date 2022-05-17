#!/usr/bin/env node
import fsSync from 'fs'
import fs from 'fs/promises'
import os from 'os'
import path from 'path'
import yaml from 'yaml'

import Markdoc from '@markdoc/markdoc'

const writingDir = `${os.homedir()}/Library/Mobile\ Documents/27N4MQEA55\~pro\~writer/Documents/Published`;
const buildDir = './dist';

(async () => {
  // Anything in here will be waited on before the process exits.
  const promises = []

  try {
    if (fsSync.existsSync(buildDir)) {
      await fs.rm(buildDir, {recursive: true})
    }
    await fs.mkdir(buildDir)

    // Add .md to any file without it
    await Promise.all((await fs.readdir(writingDir)).map(filename => {
      if (filename.substring(filename.length - ".md".length) == ".md") {
        return Promise.resolve()
      }

      if (filename.substring(filename.length - ".txt".length) == ".txt") {
        return fs.rename(
          `${writingDir}/${filename}`,
          `${writingDir}/${filename.substring(0, filename.length - ".txt".length)}.md`,
        )
      }

      return fs.rename(`${writingDir}/${filename}`, `${writingDir}/${filename}.md`)
    }))

    const header = fs.readFile(`./src/_header.html`)
    const footer = fs.readFile(`./src/_footer.html`)
    promises.push(fs.copyFile('./src/index.html', `${buildDir}/index.html`))
    promises.push(fs.copyFile('./src/style.css', `${buildDir}/style.css`))

    // Re-glob after renames 
    // TODO: Just remember the renames
    const files = (await fs.readdir(writingDir)).map(async filename => {
      return fs.readFile(`${writingDir}/${filename}`, {encoding: 'utf8'}).then(async body => {
        const stat = fs.stat(`${writingDir}/${filename}`)
        const ast = Markdoc.parse(body)
        const content = Markdoc.transform(ast)
        let html = Markdoc.renderers.html(content)
        if (!ast.attributes.frontmatter) {
          ast.attributes.frontmatter = {some: "frontmatter"}
        }
        const frontmatter = yaml.parse(ast.attributes.frontmatter) || {}
        let error = false
        if (!frontmatter.shortname) {
          frontmatter.shortname = "TODO"
          error = true
        }
        const date = new Date()
        if (!frontmatter.date) {
          frontmatter.date = "TODO"
          error = true
        } else {
          const parts = frontmatter.date.split('-')
          date.setYear(parts[0])
          if (parts[1]) {
            date.setMonth(parts[1] - 1)
          }
          if (parts[2]) {
            date.setDate(parts[2])
          }
          let dateStr = `${date.toLocaleString('default', {month: 'long'})} ${date.getFullYear()}`
          if ((await stat).mtime.getMonth() != date.getMonth() || (await stat).mtime.getFullYear() != date.getFullYear()) {
            dateStr += `; last updated ${(await stat).mtime.toLocaleString('default', {month: 'long'})} ${(await stat).mtime.getFullYear()}`
          }
          html = html.replace('</h1>', `</h1><p>${dateStr}</p>`)
        }
        if (error) {
          promises.push(fs.writeFile(
            `${writingDir}/${filename}`,
            `---\n${yaml.stringify(frontmatter)}---\n\n${body}`,
          ))
        }
        return {
          shortname: frontmatter.shortname || 'error',
          date,
          filename: filename.substring(0, filename.length - ".md".length),
          title: filename.substring(0, filename.length - ".md".length),
          body: html,
        }
      })
    })

    await Promise.all(
      files.map(async file => fs.writeFile(
        `${buildDir}/${(await file).shortname}.html`,
        (await header) + (await file).body + (await footer),
      ))
    )
    await Promise.all(promises)
  } catch (e) {
    console.error(e)
    await Promise.all(promises)
    process.exit(1)
  }
})()
