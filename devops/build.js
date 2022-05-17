const fsSync = require('fs')
const fs = require('fs/promises')
const os = require('os')
const path = require('path')
const yaml = require('yaml')

const Markdoc = require('@markdoc/markdoc')

const sourceDir = `${os.homedir()}/Library/Mobile\ Documents/27N4MQEA55\~pro\~writer/Documents/Published`;
const destinationDir = './dist';

(async () => {
  try {
    if (fsSync.existsSync(destinationDir)) {
      await fs.rm(destinationDir, {recursive: true})
    }
    await fs.mkdir(destinationDir)

    // Add .md to any file without it
    await Promise.all((await fs.readdir(sourceDir)).map(filename => {
      if (filename.substring(filename.length - ".md".length) == ".md") {
        return Promise.resolve()
      }

      if (filename.substring(filename.length - ".txt".length) == ".txt") {
        return fs.rename(
          `${sourceDir}/${filename}`,
          `${sourceDir}/${filename.substring(0, filename.length - ".txt".length)}.md`,
        )
      }

      return fs.rename(`${sourceDir}/${filename}`, `${sourceDir}/${filename}.md`)
    }))

    // Re-glob after renames 
    // TODO: Just remember the renames
    const files = (await fs.readdir(sourceDir)).map(filename => {
      return fs.readFile(`${sourceDir}/${filename}`, {encoding: 'utf8'}).then(body => {
        const ast = Markdoc.parse(body)
        const content = Markdoc.transform(ast)
        const html = Markdoc.renderers.html(content)
        if (!ast.attributes.frontmatter) {
          throw `Frontmatter not found for "${filename}"`
        }
        return {
          shortname: yaml.parse(ast.attributes.frontmatter).shortname,
          filename: filename.substring(0, filename.length - ".md".length),
          title: filename.substring(0, filename.length - ".md".length),
          body: html,
        }
      })
    })

    await Promise.all(
      files.map(async file => fs.writeFile(
        `${destinationDir}/${(await file).shortname}.html`,
        (await file).body,
      ))
    )
  } catch (e) {
    console.error(e)
    process.exit(1)
  }
})()
