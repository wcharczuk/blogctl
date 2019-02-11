package sh

import "io"

// Tee returns a new tee writer for a given set of writers.
func Tee(writers ...io.Writer) io.Writer {
	return TeeWriter(writers)
}

// Assert that tee writer implements io.Writer.
var (
	_ io.Writer = (*TeeWriter)(nil)
)

// TeeWriter returns a io.Writer for a given array of writers.
type TeeWriter []io.Writer

// Write writes the contents to each of the streams.
func (tw TeeWriter) Write(contents []byte) (int, error) {
	var count int
	var err error
	for _, writer := range tw {
		if count, err = writer.Write(contents); err != nil {
			return count, err
		}
	}
	return count, err
}
