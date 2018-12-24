package logger

import (
	"io"
	"sync"
)

// NewInterlockedWriter returns a new interlocked writer.
func NewInterlockedWriter(output io.Writer) io.Writer {
	return &InterlockedWriter{
		output:   output,
		syncRoot: sync.Mutex{},
	}
}

// InterlockedWriter is a writer that serializes access to the Write() method.
type InterlockedWriter struct {
	output   io.Writer
	syncRoot sync.Mutex
}

// Write writes the given bytes to the inner writer.
func (iw *InterlockedWriter) Write(buffer []byte) (count int, err error) {
	iw.syncRoot.Lock()

	count, err = iw.output.Write(buffer)
	if err != nil {
		iw.syncRoot.Unlock()
		return
	}
	iw.syncRoot.Unlock()
	return
}

// Close closes any outputs that are io.WriteCloser's.
func (iw *InterlockedWriter) Close() (err error) {
	iw.syncRoot.Lock()
	defer iw.syncRoot.Unlock()

	if typed, isTyped := iw.output.(io.WriteCloser); isTyped {
		err = typed.Close()
		if err != nil {
			return
		}
	}
	return
}
