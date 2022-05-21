#!/usr/bin/env node
import fsSync from 'fs-extra'
import fs from 'fs/promises'
import yaml from 'yaml'

import Markdoc from '@markdoc/markdoc'

const src = './src';
const dist = './dist';

(async () => {
  // Anything in here will be waited on before the process exits.
  const promises = []

  try {
    if (fsSync.existsSync(dist)) {
      await fs.rm(dist, {recursive: true})
    }
    await fs.mkdir(dist)

    const header = fs.readFile(`${src}/templates/_header.html`, {encoding: 'utf8'});
    const footer = fs.readFile(`${src}/templates/_footer.html`, {encoding: 'utf8'});

    (await fs.readdir(src)).filter(filename => {
      return filename.substring(filename.length - ".md".length) == ".md"
    }).map(filename => fs.readFile(`${src}/${filename}`, {encoding: 'utf8'}).then(async body => {
      const stat = fs.stat(`${src}/${filename}`)
      const ast = Markdoc.parse(body)
      const content = Markdoc.transform(ast)
      let html = Markdoc.renderers.html(content)
      const frontmatter = yaml.parse(ast.attributes.frontmatter) || {}

      const date = new Date()
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
      promises.push(fs.writeFile(
        `${dist}/${frontmatter.filename}`,
        (await header).replaceAll(
          "{{title}}",
          filename.substring(0, filename.length - ".md".length),
        ) + html + (await footer).replace(
          `https://github.com/glacials/twos.dev`,
          `https://github.com/glacials/twos.dev/blob/main/src/${filename}`,
        )
      ))
    }))

    await Promise.all(promises)
    fsSync.copySync('./public', './dist')
  } catch (e) {
    console.error(e)
    await Promise.all(promises)
    process.exit(1)
  }
})()
