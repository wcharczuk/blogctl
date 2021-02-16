/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package jwt

// Keyfunc should return the key used in verification based on the raw token passed to it.
type Keyfunc func(*Token) (interface{}, error)

// KeyfuncStatic returns a static key func.
func KeyfuncStatic(key []byte) Keyfunc {
	return func(_ *Token) (interface{}, error) {
		return key, nil
	}
}
