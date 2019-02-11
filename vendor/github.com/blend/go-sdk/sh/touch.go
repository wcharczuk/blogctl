package sh

import "os"

// Touch creates a file at a given path.
func Touch(path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	return nil
}
