package jwt

import "github.com/blend/go-sdk/exception"

// Error constants.
var (
	// ErrValidation will be the top most class in most cases.
	ErrValidation exception.Class = "validation error"

	ErrValidationAudienceUnset exception.Class = "token claims audience unset"
	ErrValidationExpired       exception.Class = "token expired"
	ErrValidationIssued        exception.Class = "token issued in future"
	ErrValidationNotBefore     exception.Class = "token not before"

	ErrValidationSignature exception.Class = "signature is invalid"

	ErrKeyfuncUnset         exception.Class = "keyfunc is unset"
	ErrInvalidKey           exception.Class = "key is invalid"
	ErrInvalidKeyType       exception.Class = "key is of invalid type"
	ErrInvalidSigningMethod exception.Class = "invalid signing method"
	ErrHashUnavailable      exception.Class = "the requested hash function is unavailable"

	ErrHMACSignatureInvalid exception.Class = "hmac signature is invalid"

	ErrECDSAVerification exception.Class = "crypto/ecdsa: verification error"

	ErrKeyMustBePEMEncoded exception.Class = "invalid key: key must be pem encoded pkcs1 or pkcs8 private key"
	ErrNotRSAPrivateKey    exception.Class = "key is not a valid rsa private key"
	ErrNotRSAPublicKey     exception.Class = "key is not a valid rsa public key"
)

// IsValidation returns if the error is a validation error
// instead of a more structural error with the key infrastructure.
func IsValidation(err error) bool {
	return exception.Is(err, ErrValidation)
}
