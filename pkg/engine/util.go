package engine

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"image"
	// these are needed to read an image metadata
	_ "image/jpeg"
	_ "image/png"

	"github.com/blend/go-sdk/yaml"
	"github.com/wcharczuk/photoblog/pkg/exif"
	"github.com/wcharczuk/photoblog/pkg/model"
)

// GetFileInfos returns all the file infos within a given directory by path.
func GetFileInfos(path string) (files []os.FileInfo, err error) {
	err = filepath.Walk(path, func(currentPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return filepath.SkipDir
		}
		files = append(files, info)
		return nil
	})
	return
}

// ReadYAML reads a yaml file into a given object reference.
func ReadYAML(path string, obj interface{}) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return yaml.NewDecoder(f).Decode(obj)
}

// HasExtension returns if a given filename has any of a given set of extensions.
func HasExtension(filename string, extensions ...string) bool {
	for _, ext := range extensions {
		if strings.HasSuffix(filename, ext) {
			return true
		}
	}
	return false
}

// ReadImage reads image metadata.
func ReadImage(path string, ref *model.Image) error {
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	image, _, err := image.DecodeConfig(bytes.NewBuffer(contents))
	if err != nil {
		return err
	}

	exif, err := exif.Decode(bytes.NewBuffer(contents))
	if err != nil {
		return err
	}

	ref = &model.Image{
		Width:  image.Width,
		Height: image.Height,
		Exif:   model.Exif(exif.Values()),
	}
	return nil
}
