package winter

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/gorilla/feeds"
	"gopkg.in/yaml.v3"
)

const (
	configRelDir   = "winter"
	configFileName = "winter.yml"
)

// Config is a configuration for the Winter build.
type Config struct {
	// Author is the information for the website author.
	// This is used in metadata such as that of the RSS feed.
	Author feeds.Author `yaml:"author,omitempty"`
	// Debug is a flag that enables debug mode.
	Debug bool `yaml:"debug,omitempty"`
	// Development contains options specific to development.
	// They have no impact when building for production.
	Development struct {
		// URL is the base URL you will connect to while developing your website or Winter.
		// If blank, it defaults to "http://localhost:8100".
		URL string `yaml:"url,omitempty"`
	} `yaml:"development,omitempty"`
	// Description is the Description of the website.
	// This is used as metadata for the RSS feed.
	Description string `yaml:"description,omitempty"`
	// Dist is the location the site will be built into,
	// relative to the working directory.
	// After a build, this directory is suitable for deployment to the web as a set of static files.
	//
	// In other words, the path of any file in dist,
	// relative to dist,
	// is equivalent to the path component of the URL for that file.
	//
	// If blank, defaults to ./dist.
	Dist       string `yaml:"dist,omitempty"`
	Production struct {
		// URL is the base URL you will connect to to view your deployed website
		// (e.g. twos.dev or one.twos.dev or twos.dev:6667).
		// This is used in various backlinks, like those in the RSS feed.
		//
		// Must not be blank.
		URL string `yaml:"url,omitempty"`
	} `yaml:"production,omitempty"`
	// Known helps the generated site follow the "Cool URIs don't change" rule
	// by remembering certain facts about what the site looks like,
	// and checking newly-generated sites against that memory.
	Known struct {
		// URIs holds the path to the known URIs file,
		// which Winter will generate, update, and maintain.
		//
		// You should commit this file.
		//
		// If unset, defaults to src/uris.txt.
		URIs string `yaml:"urls,omitempty"`
	} `yaml:"known,omitempty"`
	// Name is the name of the website.
	// This is used in various places in and out of templates.
	Name string `yaml:"name,omitempty"`
	// Since is the year the website was established, whether through Winter or otherwise.
	// This is used as metadata for the RSS feed.
	//
	// TODO: Use this for copyright in page footer
	Since int `yaml:"since,omitempty"`
	// Src is an additional list of directories to search for source files beyond ./src.
	Src []string `yaml:"srca,omitempty"`
}

func NewConfig() (Config, error) {
	var c Config
	p, err := ConfigPath()
	if err != nil {
		return Config{}, err
	}
	f, err := os.ReadFile(p)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			fmt.Println("No config file found. Run winter init to create one interactively.")
		}
	}
	if err := yaml.Unmarshal(f, &c); err != nil {
		return Config{}, err
	}
	if c.Development.URL == "" {
		c.Development.URL = "http://localhost:8100"
	}
	if c.Production.URL == "" {
		return Config{}, fmt.Errorf("production.url must be specified in winter.yml")
	}
	if c.Known.URIs == "" {
		c.Known.URIs = "src/uris.txt"
	}
	for i := range c.Src {
		c.Src[i] = os.ExpandEnv(strings.ReplaceAll(c.Src[i], "~", "$HOME"))
	}
	return c, nil
}

func InteractiveConfig() error {
	p, err := ConfigPath()
	if err != nil {
		return err
	}
	w, err := os.OpenFile(p, os.O_CREATE, 0o644)
	if err != nil {
		return err
	}
	var c Config
	c.Author.Name = ask("Author name:")
	c.Author.Email = ask("Author email:")
	b, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	if _, err := w.Write(b); err != nil {
		return err
	}
	return nil
}

func (c Config) Save() error {
	bytes, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	p, err := ConfigPath()
	if err != nil {
		return err
	}
	return os.WriteFile(p, bytes, fs.FileMode(os.O_WRONLY))
}

func ask(question string) (answer string) {
	fmt.Println(question)
	fmt.Scanf(answer)
	return
}

func ConfigPath() (string, error) {
	if _, err := os.Stat(configFileName); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return configFileName, nil
		}
	} else {
		return configFileName, nil
	}
	userCfg, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(userCfg, configRelDir, configFileName), nil
}
