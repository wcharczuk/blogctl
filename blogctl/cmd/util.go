package cmd

import (
	"io/ioutil"

	"github.com/blend/go-sdk/exception"
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
