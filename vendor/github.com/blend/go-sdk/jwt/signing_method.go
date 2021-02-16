/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package jwt

import "crypto"

// Common signing method names.
const (
	SigningMethodNameHMAC256 = "HS256"
	SigningMethodNameHMAC384 = "HS384"
	SingingMethodNameHMAC512 = "HS512"

	SigningMethodNameES256 = "ES256"
	SigningMethodNameES384 = "ES384"
	SigningMethodNameES512 = "ES512"

	SigningMethodNameRS256 = "RS256"
	SigningMethodNameRS384 = "RS384"
	SigningMethodNameRS512 = "RS512"
)

// SigningMethod is a type that implements methods required to sign tokens.
type SigningMethod interface {
	Verify(signingString, signature string, key interface{}) error // Returns nil if signature is valid
	Sign(signingString string, key interface{}) (string, error)    // Returns encoded signature or error
	Alg() string                                                   // returns the alg identifier for this method (example: 'HS256')
}

// Static references for specific signing methods.
var (
	SigningMethodHMAC256 = &SigningMethodHMAC{SigningMethodNameHMAC256, crypto.SHA256}
	SigningMethodHMAC384 = &SigningMethodHMAC{SigningMethodNameHMAC384, crypto.SHA384}
	SigningMethodHMAC512 = &SigningMethodHMAC{SingingMethodNameHMAC512, crypto.SHA512}

	SigningMethodES256 = &SigningMethodECDSA{SigningMethodNameES256, crypto.SHA256, 32, 256}
	SigningMethodES384 = &SigningMethodECDSA{SigningMethodNameES384, crypto.SHA384, 48, 384}
	SigningMethodES512 = &SigningMethodECDSA{SigningMethodNameES512, crypto.SHA512, 66, 521}

	SigningMethodRS256 = &SigningMethodRSA{SigningMethodNameRS256, crypto.SHA256}
	SigningMethodRS384 = &SigningMethodRSA{SigningMethodNameRS384, crypto.SHA384}
	SigningMethodRS512 = &SigningMethodRSA{SigningMethodNameRS512, crypto.SHA512}
)

// GetSigningMethod returns a signing method with a given name.
func GetSigningMethod(name string) SigningMethod {
	switch name {
	case SigningMethodNameHMAC256:
		return SigningMethodHMAC256
	case SigningMethodNameHMAC384:
		return SigningMethodHMAC384
	case SingingMethodNameHMAC512:
		return SigningMethodHMAC512
	case SigningMethodNameES256:
		return SigningMethodES256
	case SigningMethodNameES384:
		return SigningMethodES384
	case SigningMethodNameES512:
		return SigningMethodES512
	case SigningMethodNameRS256:
		return SigningMethodRS256
	case SigningMethodNameRS384:
		return SigningMethodRS384
	case SigningMethodNameRS512:
		return SigningMethodRS512
	default:
		return nil
	}
}
