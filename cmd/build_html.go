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
)

type htmlFileVars struct {
	Body      template.HTML
	Title     string
	SourceURL string

	CreatedAt string
	UpdatedAt string
}

func buildHTML(sourceDir, destinationDir string) error {
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
			return fmt.Errorf("can't walk path `%s` for HTML templates: %w", path, err)
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

		desiredHTMLFilename, err := desiredFilename(path)
		if err != nil {
			return fmt.Errorf("can't generate desired filename from Markdown file: %w", err)
		}

		relativeHTMLPath, err := filepath.Rel(sourceDir, filepath.Join(filepath.Dir(path), desiredHTMLFilename))
		if err != nil {
			return fmt.Errorf("can't get relative path of `%s`: %w", desiredHTMLFilename, err)
		}

		if err := os.MkdirAll(filepath.Join(destinationDir, filepath.Dir(relativeHTMLPath)), 0755); err != nil {
			return fmt.Errorf("can't make src-equivalent directory in dist: %w", err)
		}

		destinationFilePath := filepath.Join(destinationDir, filepath.Dir(relativeHTMLPath), desiredHTMLFilename)
		htmlFile, err := os.Create(destinationFilePath)
		if err != nil {
			return fmt.Errorf("can't create HTML file from templates for `%s`: %w", path, err)
		}

		renderedHTML := bytes.NewBuffer([]byte{})
		renderCmd := exec.Command("src/js/build.js", "body", path)
		renderCmd.Stdout = renderedHTML
		renderCmd.Stderr = os.Stderr

		if err := renderCmd.Run(); err != nil {
			return fmt.Errorf("can't run `src/js/build.js body '%s'`: %w", path, err)
		}

		v := htmlFileVars{
			Body:      template.HTML(renderedHTML.String()),
			Title:     d.Name()[0 : len(d.Name())-len(".md")],
			SourceURL: fmt.Sprintf("https://github.com/glacials/twos.dev/blob/main/src/%s", relativeHTMLPath),

			CreatedAt: "TODO",
			UpdatedAt: "TODO",
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

func desiredFilename(markdownFilePath string) (string, error) {
	filename := bytes.NewBuffer([]byte{})
	filenameCmd := exec.Command("src/js/build.js", "filename", markdownFilePath)
	filenameCmd.Stdout = filename
	filenameCmd.Stderr = os.Stderr

	if err := filenameCmd.Run(); err != nil {
		return "", fmt.Errorf("can't run `src/js/build.js filename '%s'`: %w", markdownFilePath, err)
	}

	return strings.TrimSpace(filename.String()), nil
}
