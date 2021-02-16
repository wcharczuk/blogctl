/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package stringutil

import "strings"

// CSV produces a csv from a given set of values.
func CSV(values []string) string {
	return strings.Join(values, ",")
}
