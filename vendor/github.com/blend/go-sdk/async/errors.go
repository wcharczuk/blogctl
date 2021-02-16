/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package async

import "github.com/blend/go-sdk/ex"

// Errors
var (
	ErrCannotStart ex.Class = "cannot start; already started"
	ErrCannotStop  ex.Class = "cannot stop; already stopped"
)
