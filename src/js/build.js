#!/usr/bin/env node
import { argv } from 'node:process'
import fs from 'fs-extra'
import yaml from 'yaml'

import Markdoc from '@markdoc/markdoc'

(async () => {
  // Two modes of calling:
  // 
  //     src/js/build.js filename FILEPATH
  //     src/js/build.js body FILEPATH
  //
  // The first will print to stdout filename the Markdown file given to stdin wants,
  // based on frontmatter. The second will print to stdout the HTML rendering of the
  // Markdown file given to stdin.
  try {
    const src = argv[3]
    const stat = fs.stat(src)
    const body = fs.readFileSync(src, {encoding: "utf8"})
    const ast = Markdoc.parse(body)
    const content = Markdoc.transform(ast)
    let html = Markdoc.renderers.html(content)
    const frontmatter = (ast.attributes.frontmatter ? yaml.parse(ast.attributes.frontmatter) : {}) || {}

    if (argv[2] == "filename") {
      console.log(frontmatter.filename)
      return
    }

    if (frontmatter.date) {
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
    }
    console.log(html)
  } catch (e) {
    console.error(e)
    process.exit(1)
  }
})()
