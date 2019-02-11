package sh

import (
	"io"
	"sync"

	"github.com/blend/go-sdk/exception"
)

// Errors
const (
	ErrMaxBytesWriterCapacityLimit exception.Class = "write failed; maximum capacity reached or would be exceede"
)

// LimitBytes returns a new max bytes writer.
func LimitBytes(maxBytes int, inner io.Writer) *MaxBytesWriter {
	return &MaxBytesWriter{
		max:   maxBytes,
		inner: inner,
	}
}

// MaxBytesWriter returns a maximum bytes writer.
type MaxBytesWriter struct {
	sync.Mutex

	max   int
	count int
	inner io.Writer
}

// Max returns the maximum bytes allowed.
func (mbw *MaxBytesWriter) Max() int {
	return int(mbw.max)
}

// Count returns the current written bytes..
func (mbw *MaxBytesWriter) Count() int {
	return int(mbw.count)
}

func (mbw *MaxBytesWriter) Write(contents []byte) (int, error) {
	mbw.Lock()
	defer mbw.Unlock()

	if mbw.count >= mbw.max {
		return 0, exception.New(ErrMaxBytesWriterCapacityLimit)
	}
	if len(contents)+mbw.count >= mbw.max {
		return 0, exception.New(ErrMaxBytesWriterCapacityLimit)
	}

	written, err := mbw.inner.Write(contents)
	mbw.count += written
	if err != nil {
		return written, exception.New(err)
	}
	return written, nil
}
