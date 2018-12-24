package web

import (
	"fmt"

	"github.com/blend/go-sdk/exception"
	"github.com/blend/go-sdk/jwt"
)

const (
	// ErrSessionIDEmpty is thrown if a session id is empty.
	ErrSessionIDEmpty exception.Class = "auth session id is empty"
	// ErrSecureSessionIDEmpty is an error that is thrown if a given secure session id is invalid.
	ErrSecureSessionIDEmpty exception.Class = "auth secure session id is empty"
	// ErrUnsetViewTemplate is an error that is thrown if a given secure session id is invalid.
	ErrUnsetViewTemplate exception.Class = "view result template is unset"

	// ErrParameterMissing is an error on request validation.
	ErrParameterMissing exception.Class = "parameter is missing"
)

func newParameterMissingError(paramName string) error {
	return fmt.Errorf("`%s` parameter is missing", paramName)
}

// IsErrSessionInvalid returns if an error is a session invalid error.
func IsErrSessionInvalid(err error) bool {
	if err == nil {
		return false
	}
	if exception.Is(err, ErrSessionIDEmpty) ||
		exception.Is(err, ErrSecureSessionIDEmpty) ||
		exception.Is(err, jwt.ErrValidation) {
		return true
	}
	return false
}
