/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package sanitize

// ValueSanitizer is a function that sanitizes values.
//
// It's designed to accept a variadic list of values that represent
// the Value of maps like `http.Header` and `url.Values`
type ValueSanitizer func(key string, values ...string) []string

// DefaultValueSanitizer is the default value sanitizer.
func DefaultValueSanitizer(_ string, _ ...string) []string {
	return nil
}
