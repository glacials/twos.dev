package cmd

import (
	"fmt"
	"os"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"twos.dev/winter"
)

func newConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Interact with Winter configuration",
		Long: wrap(`
			Interact with Winter configuration.

			Configuration in ./winter.yml takes first precedence.
			Otherwise, configuration is stored according to the XDG spec.
			This is typically

			    ~/.config/winter/winter.yml

			on Linux and

			    ~/Library/Application Support/winter/winter.yml

			on macOS. For more information on possible locations, see the XDG specification:
			https://specifications.freedesktop.org/basedir-spec/basedir-spec-latest.html
		`),
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := winter.InteractiveConfig(); err != nil {
				return err
			}
			return nil
		},
	}
	cmd.AddCommand(newConfigGetCmd())
	cmd.AddCommand(newConfigClearCmd())
	return cmd
}

func newConfigGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get [key]",
		Short: "Get or list Winter config",
		Long: wrap(`
			Get the value of the Winter configuration variable named KEY,
			or all configuration if KEY is omitted.

			See winter config --help for information on configuration storage locations.
		`),
		Args: cobra.RangeArgs(0, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := winter.NewConfig()
			if err != nil {
				return err
			}
			if len(args) == 0 {
				bytes, err := yaml.Marshal(c)
				if err != nil {
					return err
				}
				if _, err := os.Stdout.Write(bytes); err != nil {
					return err
				}
				return nil
			}
			if len(args) > 1 {
				return fmt.Errorf("must take 0–1 arguments")
			}
			for _, arg := range args {
				if err := mapstructure.Decode(arg, &c); err != nil {
					return err
				}
			}
			return c.Save()
		},
	}
}

func newConfigClearCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "clear",
		Short: "Erase all config",
		Long: wrap(`
			Erase all Winter configuration. Cannot be undone.
			See winter config --help for information on configuration storage locations.
		`),
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			p, err := winter.ConfigPath()
			if err != nil {
				return err
			}
			return os.Remove(p)
		},
	}
}
