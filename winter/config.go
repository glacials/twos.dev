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
	// Author is the information for the website author. This is used in metadata such as
	// that of the RSS feed.
	Author feeds.Author `yaml:",omitempty"`
	// Debug is a flag that enables debug mode.
	Debug bool `yaml:",omitempty"`
	// Description is the Description of the website. This is used as metadata
	// for the RSS feed.
	Description string `yaml:",omitempty"`
	// Dist is the path to the distribution directory. If blank,
	// defaults to "./dist" relative to the working directory.
	Dist string `yaml:",omitempty"`
	// Hostname is the host portion of the website URL (e.g. "one.twos.dev" for a URL of
	// "https://one.twos.dev/index.html"). This is used in various backlinks, like those
	// in the RSS feed.
	Hostname string `yaml:",omitempty"`
	// Name is the name of the website. This is used in various places
	// in and out of templates.
	Name string `yaml:",omitempty"`
	// Since is the year the website was established, whether through
	// Winter or otherwise. This is used as metadata for the RSS feed.
	//
	// TODO: Use this for copyright in page footer
	Since int `yaml:",omitempty"`
	// Src is an additional list of directories to search
	// for source files beyond ./src.
	Src []string `yaml:",omitempty"`
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
