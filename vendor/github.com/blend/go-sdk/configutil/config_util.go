package configutil

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/blend/go-sdk/env"
	"github.com/blend/go-sdk/exception"
	"github.com/blend/go-sdk/yaml"
)

const (
	// EnvVarConfigPath is the env var for configs.
	EnvVarConfigPath = "CONFIG_PATH"
	// ExtensionJSON is a file extension.
	ExtensionJSON = ".json"
	// ExtensionYAML is a file extension.
	ExtensionYAML = ".yaml"
	// ExtensionYML is a file extension.
	ExtensionYML = ".yml"
)

var (
	// DefaultPaths are default path locations.
	// They are tested and read in order, so the later
	// paths will override data found in the earlier ones.
	DefaultPaths = []string{
		"/var/secrets/config.yml",
		"/var/secrets/config.yaml",
		"/var/secrets/config.json",
		"./_config/config.yml",
		"./_config/config.yaml",
		"./_config/config.json",
		"./config.yml",
		"./config.yaml",
		"./config.json",
	}
)

// Read reads a config from optional path(s).
// Paths will be tested from a standard set of defaults (ex. config.yml)
// and optionally a csv named in the `CONFIG_PATH` environment variable.
func Read(ref Any, paths ...string) error {
	_, err := TryReadFromPaths(ref, PathsWithDefaults(paths...)...)
	return err
}

// TryReadFromPaths tries to read the config from a list of given paths, reading from the first file that exists.
func TryReadFromPaths(ref Any, paths ...string) (path string, err error) {
	if len(paths) == 0 {
		err = exception.New(ErrConfigPathUnset)
		return
	}

	// for each of the paths
	// if the path doesn't exist, continue, read the path that is found.
	var f *os.File
	for _, path = range paths {
		if path == "" {
			continue
		}

		f, err = os.Open(path)
		if IsNotExist(err) {
			continue
		}
		if err != nil {
			err = exception.New(err)
			return
		}
		defer f.Close()

		err = ReadFromReader(ref, f, filepath.Ext(path))
		return
	}
	err = exception.New(os.ErrNotExist).WithMessagef("no provided paths exist: %#v", paths)
	return
}

// ReadFromReader reads a config from a given reader.
func ReadFromReader(ref Any, r io.Reader, ext string) error {
	if err := Deserialize(ext, r, ref); err != nil {
		return err
	}
	return env.Env().ReadInto(ref)
}

// PathsWithDefaults returns the default paths and additional optional paths.
func PathsWithDefaults(paths ...string) []string {
	return Paths(append(DefaultPaths, paths...)...)
}

// Paths returns config paths.
// The results are the provided defaults and the `CONFIG_PATH`
// environment variable as a csv if it's set.
func Paths(defaults ...string) (output []string) {
	if env.Env().Has(EnvVarConfigPath) {
		output = append(output, env.Env().CSV(EnvVarConfigPath)...)
	}
	output = append(output, defaults...)
	return
}

// Deserialize deserializes a config.
func Deserialize(ext string, r io.Reader, ref Any) error {
	// make sure the extension starts with a "."
	if !strings.HasPrefix(ext, ".") {
		ext = "." + ext
	}

	// based off the extension, use the appropriate deserializer
	switch strings.ToLower(ext) {
	case ExtensionJSON:
		return exception.New(json.NewDecoder(r).Decode(ref))
	case ExtensionYAML, ExtensionYML:
		return exception.New(yaml.NewDecoder(r).Decode(ref))
	default: // return an error if we're passed a weird extension
		return exception.New(ErrInvalidConfigExtension).WithMessagef("extension: %s", ext)
	}
}
