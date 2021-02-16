/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package graceful

import (
	"os"
	"os/signal"
)

// Notify returns a channel that listens for a given set of os signals.
func Notify(signals ...os.Signal) chan os.Signal {
	terminateSignal := make(chan os.Signal, 1)
	signal.Notify(terminateSignal, signals...)
	return terminateSignal
}
