/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package logger

import (
	"log"
)

// StdlibShim returns a stdlib logger that writes to a given logger instance.
func StdlibShim(handler Triggerable, opts ...ShimWriterOption) *log.Logger {
	shim := NewShimWriter(handler)
	for _, opt := range opts {
		opt(&shim)
	}
	return log.New(shim, "", 0)
}
