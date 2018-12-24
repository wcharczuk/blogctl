package fileutil

import (
	"io"
	"os"

	"github.com/blend/go-sdk/exception"
)

// ReadChunks reads a file in `chunkSize` pieces, dispatched to the handler.
func ReadChunks(filePath string, chunkSize int, handler func([]byte) error) error {
	f, err := os.Open(filePath)
	if err != nil {
		return exception.New(err)
	}
	defer f.Close()

	chunk := make([]byte, chunkSize)
	for {
		readBytes, err := f.Read(chunk)
		if err == io.EOF {
			break
		}
		readData := chunk[:readBytes]
		err = handler(readData)
		if err != nil {
			return exception.New(err)
		}
	}
	return nil
}
