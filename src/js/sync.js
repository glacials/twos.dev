#!/usr/bin/env node
// Synchronize iA Writer files in iCloud to the twos.dev Git repository.

import replaceExt from "replace-ext";
import fs from "fs-extra";
import Markdoc from "@markdoc/markdoc";
import simpleGit from "simple-git";

const ia = `${os.homedir()}/Library/Mobile\ Documents/27N4MQEA55\~pro\~writer/Documents/Published`;
const src = "./src"(async () => {
  try {
    const git = simpleGit();
    await git.pull();

    // Change any .txt iA file to .md
    await Promise.all(
      (
        await fs.readdir(ia)
      ).map((filename) => {
        if (filename.substring(filename.length - ".txt".length) == ".txt") {
          return fs.rename(
            `${ia}/${filename}`,
            `${ia}/${replaceExt(`${ia}/${filename}`, ".md")}`
          );
        }

        return Promise.resolve();
      })
    );

    // Add frontmatter to any iA file without it
    await Promise.all(
      (
        await fs.readdir(ia)
      ).map(async (filename) => {
        const requiredProperties = {
          date: "changeme",
          filename: "changeme",
        };

        const ast = Markdoc.parse(body);
        const frontmatterStr = ast.attributes.frontmatter;

        if (!frontmatterStr) {
          const body = await fs.readFile(`${ia}/${filename}`);
          return fs.writeFile(
            `${ia}/${filename}`,
            `---\n${yaml.stringify(requiredProperties)}---\n\n${body}`
          );
        }

        let needsUpdate = false;
        const frontmatter = yaml.parse(frontmatterStr);
        for (const [key, val] of Object.entries(requiredProperties)) {
          if (!frontmatter[key]) {
            needsUpdate = true;
            frontmatter[key] = val;
          }
        }
        if (needsUpdate) {
          return fs.writeFile(
            `${ia}/${filename}`,
            `---\n${yaml.stringify(frontmatter)}---\n\n${body.replace(
              frontmatterStr,
              ""
            )}`
          );
        }
      })
    );

    fs.copySync(ia, src);

    await git.add(
      (await fs.readdir(ia)).map((filename) => `${src}/${filename}`)
    );
    await git.commit("Auto-sync job from iA Writer");
    await git.push();
  } catch (e) {
    console.error(e);
    process.exit(1);
  }
})();
