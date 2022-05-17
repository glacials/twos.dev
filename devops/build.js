#!/usr/bin/env node

const fsSync = require('fs')
const fs = require('fs/promises')
const os = require('os')
const path = require('path')
const yaml = require('yaml')

const Markdoc = require('@markdoc/markdoc')

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
    const files = (await fs.readdir(writingDir)).map(filename => {
      return fs.readFile(`${writingDir}/${filename}`, {encoding: 'utf8'}).then(body => {
        const ast = Markdoc.parse(body)
        const content = Markdoc.transform(ast)
        const html = Markdoc.renderers.html(content)
        if (!ast.attributes.frontmatter) {
          promises.push(fs.writeFile(
            `${writingDir}/${filename}`,
            `---\nerror: Must specify shortname\n---\n\n${body}`,
          ))
          throw `Frontmatter not found for "${filename}"`
        }
        const frontmatter = yaml.parse(ast.attributes.frontmatter)
        if (!frontmatter.shortname) {
          frontmatter.error = "Must specify shortname"
          promises.push(fs.writeFile(
            `${writingDir}/${filename}`,
            body.replace(ast.attributes.frontmatter + '\n', `${yaml.stringify(frontmatter)}`),
          ))
        }
        return {
          shortname: frontmatter.shortname,
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
