/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package configutil

// StringPtr returns a StringSource for a given string pointer.
func StringPtr(value *string) StringSource {
	if value == nil || *value == "" {
		return String("")
	}
	return String(*value)
}
