package sh

import "os"

// MustToFile opens or creates a file and panics on error.
func MustToFile(path string) *os.File {
	file, err := ToFile(path)
	if err != nil {
		panic(err)
	}
	return file
}

// ToFile opens or creates a file.
func ToFile(path string) (*os.File, error) {
	if _, err := os.Stat(path); err == nil {
		return os.Open(path)
	}
	return os.Create(path)
}
