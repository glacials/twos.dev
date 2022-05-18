#!/usr/bin/env node
import fsSync from 'fs-extra'
import fs from 'fs/promises'
import os from 'os'
import path from 'path'
import yaml from 'yaml'

import Markdoc from '@markdoc/markdoc'
import replaceExt from 'replace-ext'

const writingDir = './src';
const buildDir = './dist';

(async () => {
  // Anything in here will be waited on before the process exits.
  const promises = []

  try {
    if (fsSync.existsSync(buildDir)) {
      await fs.rm(buildDir, {recursive: true})
    }
    await fs.mkdir(buildDir)

    // Change any .txt or no-extension file to .md
    await Promise.all((await fs.readdir(writingDir)).map(filename => {
      if (filename.substring(filename.length - ".txt".length) == ".txt" || filename.indexOf('.') < 0) {
        return fs.rename(
          `${writingDir}/${filename}`,
          `${writingDir}/${replaceExt(`${writingDir}/${filename}`, '.md')}`,
        )
      }

      return Promise.resolve()
    }))

    const header = fs.readFile(`./src/_header.html`)
    const footer = fs.readFile(`./src/_footer.html`)

    // Re-glob after renames 
    // TODO: Just remember the renames
    const files = (await fs.readdir(writingDir)).map(async filename => {
      return fs.readFile(`${writingDir}/${filename}`, {encoding: 'utf8'}).then(async body => {
        const stat = fs.stat(`${writingDir}/${filename}`)
        const ast = Markdoc.parse(body)
        const content = Markdoc.transform(ast)
        let html = Markdoc.renderers.html(content)
        const originalFrontmatter = ast.attributes.frontmatter
        if (!ast.attributes.frontmatter) {
          ast.attributes.frontmatter = 'some: frontmatter\n'
        }
        const frontmatter = yaml.parse(ast.attributes.frontmatter) || {}
        let error = false
        if (!frontmatter.filename) {
          frontmatter.filename = "TODO"
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
            `---\n${yaml.stringify(frontmatter)}---\n\n` + body.replace(
              `---\n${originalFrontmatter}\n---\n\n`,
              '',
            )
          ))
        }
        return {
          desiredFilename: frontmatter.filename || 'TODO.html',
          date,
          filename: filename.substring(0, filename.length - ".md".length),
          title: filename.substring(0, filename.length - ".md".length),
          body: html,
        }
      })
    })

    await Promise.all(
      files.map(async file => fs.writeFile(
        `${buildDir}/${(await file).desiredFilename}`,
        (await header) + (await file).body + (await footer),
      ))
    )
    await Promise.all(promises)
    fsSync.copySync('./public', './dist')
  } catch (e) {
    console.error(e)
    await Promise.all(promises)
    process.exit(1)
  }
})()
