package fileutil

import (
	"os"

	"github.com/blend/go-sdk/exception"
)

// RemoveMany removes an array of files.
func RemoveMany(filePaths ...string) error {
	var err error
	for _, path := range filePaths {
		err = os.Remove(path)
		if err != nil {
			return exception.New(err).WithMessage(path)
		}
	}
	return err
}
