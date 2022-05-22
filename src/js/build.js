#!/usr/bin/env node
import { argv } from 'node:process'
import fs from 'fs-extra'
import yaml from 'yaml'

import Markdoc from '@markdoc/markdoc'

(async () => {
  // Two modes of calling:
  // 
  //     src/js/build.js frontmatter FILEPATH
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

    if (argv[2] == "frontmatter") {
      if (ast.attributes.frontmatter) {
        console.log(ast.attributes.frontmatter)
      }
      return
    }

    if (argv[2] == "body") {
      const content = Markdoc.transform(ast)
      const html = Markdoc.renderers.html(content)
      if (html) {
        console.log(html)
      }
      return
    }
  } catch (e) {
    console.error(e)
    process.exit(1)
  }
})()
