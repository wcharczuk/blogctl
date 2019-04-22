package engine

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"image"
	// these are needed to read an image metadata
	_ "image/jpeg"
	_ "image/png"

	"github.com/blend/go-sdk/ex"
	"github.com/blend/go-sdk/yaml"

	"github.com/wcharczuk/blogctl/pkg/exif"
	"github.com/wcharczuk/blogctl/pkg/model"
	"github.com/wcharczuk/blogctl/pkg/stringutil"
)

// ListDirectory returns all the file infos within a given directory by path.
func ListDirectory(path string) (files []os.FileInfo, err error) {
	err = ex.New(filepath.Walk(path, func(currentPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if currentPath == path {
			return nil
		}
		if info.IsDir() {
			return filepath.SkipDir
		}
		files = append(files, info)
		return nil
	}))
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
func ReadImage(path string) (model.Image, error) {
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return model.Image{}, err
	}

	image, _, err := image.DecodeConfig(bytes.NewBuffer(contents))
	if err != nil {
		return model.Image{}, err
	}

	rawExifData, err := exif.Decode(bytes.NewBuffer(contents))
	if err != nil {
		return model.Image{}, err
	}

	exifData, err := GetExifData(rawExifData)
	if err != nil {
		return model.Image{}, err
	}

	return model.Image{
		Width:  image.Width,
		Height: image.Height,
		Exif:   exifData,
	}, nil
}

// GetExifData generates the parsed exif data for the post.
func GetExifData(exifData *exif.Exif) (data model.Exif, err error) {
	// fnumber
	if tag, tagErr := exifData.Get(exif.FNumber); tagErr == nil {
		nominator, denominator, ratErr := tag.Rat2(0)
		if ratErr != nil {
			err = ex.New(ratErr)
			return
		}

		if denominator != 0 {
			data.FNumber = fmt.Sprintf("F%g", float64(nominator)/float64(denominator))
		}
	}

	if tag, tagErr := exifData.Get(exif.FocalLength); tagErr == nil {
		nominator, denominator, ratErr := tag.Rat2(0)
		if ratErr != nil {
			err = ex.New(ratErr)
			return
		}
		if denominator != 0 {
			data.FocalLength = fmt.Sprintf("%gmm", float64(nominator)/float64(denominator))
		}
	}

	if tag, tagErr := exifData.Get(exif.ExposureTime); tagErr == nil {
		data.ExposureTime = stringutil.StripQuotes(tag.String()) + " sec"
	}
	if tag, tagErr := exifData.Get(exif.ISOSpeedRatings); tagErr == nil {
		data.ISOSpeedRatings = stringutil.StripQuotes(tag.String())
	}
	if tag, tagErr := exifData.Get(exif.Make); tagErr == nil {
		data.CameraMake = stringutil.StripQuotes(tag.String())
	}
	if tag, tagErr := exifData.Get(exif.Model); tagErr == nil {
		data.CameraModel = stringutil.StripQuotes(tag.String())
	}
	if tag, tagErr := exifData.Get(exif.LensMake); tagErr == nil {
		data.LensMake = tag.String()
	}
	if tag, tagErr := exifData.Get(exif.LensModel); tagErr == nil {
		data.LensModel = tag.String()
	}
	return
}

// GenerateExifData gets exif data from a file on disk.
func GenerateExifData(imagePath string) (*exif.Exif, error) {
	contents, err := ioutil.ReadFile(imagePath)
	if err != nil {
		return nil, ex.New(err)
	}

	rawExifData, err := exif.Decode(bytes.NewBuffer(contents))
	if err != nil {
		return nil, ex.New(err)
	}

	return rawExifData, nil
}

// ExtractCaptureDate extracts the capture date from an image file by path.
func ExtractCaptureDate(imagePath string) (captureDate time.Time, err error) {
	var exifData *exif.Exif
	exifData, err = GenerateExifData(imagePath)
	if err != nil {
		return
	}
	captureDate, err = exifData.DateTime()
	return
}

// MakeDir creates a new directory.
func MakeDir(path string) error {
	return ex.New(os.MkdirAll(path, 0755))
}

// WriteFile writes a file with default perms.
func WriteFile(path string, contents []byte) error {
	return ex.New(ioutil.WriteFile(path, contents, 0666))
}

// WriteYAML writes an object as yaml to disk.
func WriteYAML(path string, obj interface{}) error {
	contents, err := yaml.Marshal(obj)
	if err != nil {
		return err
	}
	return ex.New(ioutil.WriteFile(path, contents, 0666))
}

// Exists returns if a given file exists.
func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
