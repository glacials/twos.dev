package cmd

import (
	"bytes"
	"fmt"
	"html/template"
	"io/fs"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v2"
)

func buildMarkdown(sourceDir, destinationDir string) error {
	templateHTML, err := ioutil.ReadFile(filepath.Join(sourceDir, "templates", "_essay.html"))
	if err != nil {
		return fmt.Errorf("can't read essay template: %w", err)
	}

	essay, err := template.New("essay").Parse(string(templateHTML))
	if err != nil {
		return fmt.Errorf("can't create essay template: %w", err)
	}

	if err := filepath.WalkDir(sourceDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("can't walk path `%s` for Markdown templates: %w", path, err)
		}

		if d.IsDir() {
			return nil
		}

		matched, err := filepath.Match("*.md", strings.ToLower(d.Name()))
		if err != nil {
			return fmt.Errorf("can't match against filename `%s` for HTML templates: %w", d.Name(), err)
		}
		if !matched {
			return nil
		}

		if d.Name()[0] == '_' {
			return nil
		}

		if strings.Contains(path, "node_modules") {
			return nil
		}

		matter, err := getFrontmatter(path)
		if err != nil {
			return fmt.Errorf("can't get frontmatter from Markdown file: %w", err)
		}

		relativeHTMLPath, err := filepath.Rel(sourceDir, filepath.Join(filepath.Dir(path), matter.Filename))
		if err != nil {
			return fmt.Errorf("can't get relative path of `%s`: %w", matter.Filename, err)
		}

		if err := os.MkdirAll(filepath.Join(destinationDir, filepath.Dir(relativeHTMLPath)), 0755); err != nil {
			return fmt.Errorf("can't make src-equivalent directory in dist: %w", err)
		}

		destinationFilePath := filepath.Join(destinationDir, filepath.Dir(relativeHTMLPath), matter.Filename)
		htmlFile, err := os.Create(destinationFilePath)
		if err != nil {
			return fmt.Errorf("can't render Markdown to `%s` (`%s` + `%s` + `%s`) from templates for `%s`: %w", destinationFilePath, destinationDir, filepath.Dir(relativeHTMLPath), matter.Filename, path, err)
		}

		renderedHTML := bytes.NewBuffer([]byte{})
		renderCmd := exec.Command("src/js/build.js", "body", path)
		renderCmd.Stdout = renderedHTML
		renderCmd.Stderr = os.Stderr

		if err := renderCmd.Run(); err != nil {
			return fmt.Errorf("can't run `src/js/build.js body '%s'`: %w", path, err)
		}

		var created string
		if matter.CreatedAt != "" {
			createdAt, err := time.Parse("2006-01", matter.CreatedAt)
			if err != nil {
				return fmt.Errorf("can't parse created date `%s` from markdown in `%s`: %w", matter.CreatedAt, path, err)
			}

			created = createdAt.Format("2006 January")
		} else if matter.Date != "" {
			date, err := time.Parse("2006-01", matter.Date)
			if err != nil {
				return fmt.Errorf("can't parse created date `%s` from markdown in `%s`: %w", matter.CreatedAt, path, err)
			}

			created = date.Format("2006 January")
		}

		var updated string
		if matter.UpdatedAt != "" {
			updatedAt, err := time.Parse("2006-01", matter.UpdatedAt)
			if err != nil {
				return fmt.Errorf("can't parse updated date `%s`: %w", matter.UpdatedAt, err)
			}

			updated = updatedAt.Format("2006 January")
		}

		v := htmlFileVars{
			Body:      template.HTML(renderedHTML.String()),
			Title:     d.Name()[0 : len(d.Name())-len(".md")],
			SourceURL: fmt.Sprintf("https://github.com/glacials/twos.dev/blob/main/%s", path),

			CreatedAt: created,
			UpdatedAt: updated,
		}

		if err := essay.Execute(htmlFile, v); err != nil {
			return fmt.Errorf("can't execute essay template: %w", err)
		}

		return nil
	}); err != nil {
		return fmt.Errorf("can't perform HTML loop: %w", err)
	}

	return nil
}

func getFrontmatter(markdownFilePath string) (essayFrontmatter, error) {
	yml := bytes.NewBuffer([]byte{})
	filenameCmd := exec.Command("src/js/build.js", "frontmatter", markdownFilePath)
	filenameCmd.Stdout = yml
	filenameCmd.Stderr = os.Stderr

	if err := filenameCmd.Run(); err != nil {
		return essayFrontmatter{}, fmt.Errorf("can't run `src/js/build.js frontmatter '%s'`: %w", markdownFilePath, err)
	}

	var out essayFrontmatter
	if err := yaml.Unmarshal(yml.Bytes(), &out); err != nil {
		return essayFrontmatter{}, fmt.Errorf("can't unmarshal frontmatter YAML in `%s`: %w", markdownFilePath, err)
	}

	if out.Date != "" && out.CreatedAt != "" {
		return essayFrontmatter{}, fmt.Errorf(
			"frontmatter for file `%s` specified date and created, but date is an alias for created; only one should be specified",
			markdownFilePath,
		)
	}

	if out.Date != "" {
		out.CreatedAt = out.Date
	} else if out.CreatedAt != "" {
		out.Date = out.CreatedAt
	}

	if out.Filename == "" {
		return essayFrontmatter{}, fmt.Errorf("no `filename` attribute in frontmatter for `%s`", markdownFilePath)
	}

	return out, nil
}
