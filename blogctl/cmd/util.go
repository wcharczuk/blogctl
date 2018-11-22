package cmd

import (
	"io/ioutil"

	"github.com/blend/go-sdk/exception"
	"github.com/blend/go-sdk/template"
	"github.com/blend/go-sdk/yaml"

	"github.com/wcharczuk/photoblog/pkg/config"
)

// ReadConfig reads a config at a given path as yaml.
func ReadConfig(configPath string) (config config.Config, err error) {
	contents, readErr := ioutil.ReadFile(configPath)
	if readErr != nil {
		err = exception.New(readErr)
		return
	}
	err = exception.New(yaml.Unmarshal(contents, &config))
	return
}

// Vars is a helper alias.
type Vars = map[string]interface{}

// WriteYAML writes an object as yaml to disk.
func WriteYAML(path string, obj interface{}) error {
	contents, err := yaml.Marshal(obj)
	if err != nil {
		return err
	}
	return exception.New(ioutil.WriteFile(path, contents, 0666))
}

// WriteFile writes a file to disk.
func WriteFile(path string, contents []byte) error {
	return exception.New(ioutil.WriteFile(path, contents, 0666))
}

// RenderTemplate renders a template to a string.
func RenderTemplate(tpl string, vars map[string]interface{}) (string, error) {
	return template.New().WithBody(tpl).WithVars(vars).ProcessString()
}
