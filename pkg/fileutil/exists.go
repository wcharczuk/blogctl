package fileutil

import "os"

// Exists returns if a given file exists.
func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
