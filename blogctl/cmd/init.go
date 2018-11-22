package cmd

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/wcharczuk/photoblog/pkg/constants"

	"github.com/blend/go-sdk/template"

	"github.com/blend/go-sdk/exception"
	"github.com/blend/go-sdk/logger"
	"github.com/spf13/cobra"
)

// Init returns the init command.
func Init(configPath *string, log *logger.Logger) *cobra.Command {
	return &cobra.Command{
		Use:   "init [NAME]",
		Short: "Initialize a new blog",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			name := args[0]

			if err := MakeDir(name); err != nil {
				log.SyncFatalExit(err)
			}
			if err := MakeDir(filepath.Join(name, constants.ImagesPath)); err != nil {
				log.SyncFatalExit(err)
			}
			if err := MakeDir(filepath.Join(name, constants.LayoutPath)); err != nil {
				log.SyncFatalExit(err)
			}
			if err := MakeDir(filepath.Join(name, constants.PartialsPath)); err != nil {
				log.SyncFatalExit(err)
			}

			//create the config
			config, err := RenderTemplate(configTemplate, Vars{"name": name})
			if err != nil {
				log.SyncFatalExit(err)
			}
			if err := WriteFile(filepath.Join(name, constants.ConfigPath), []byte(config)); err != nil {
				log.SyncFatalExit(err)
			}
		},
	}
}

// MakeDir creates a new directory.
func MakeDir(path string) error {
	return exception.New(os.MkdirAll(path, 0666))
}

// Vars is a helper alias.
type Vars = map[string]interface{}

// WriteFile writes a file to disk.
func WriteFile(path string, contents []byte) error {
	return exception.New(ioutil.WriteFile(path, contents, 0666))
}

// RenderTemplate renders a template to a string.
func RenderTemplate(tpl string, vars map[string]interface{}) (string, error) {
	return template.New().WithBody(tpl).WithVars(vars).ProcessString()
}

const (
	configTemplate = `
name: {{ .Var "name" }}
images: images
output: dist
layout:
	index: layout/index.html
	home: layout/home.html
	single: layout/single.html
	partials:
	- layout/partials/image.html
	- layout/partials/list.html
	- layout/partials/about.html
	`

	indexHtml = `

	`
)
