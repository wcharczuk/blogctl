package jwt

import (
	"encoding/base64"
	"strings"
)

// Parse validates and returns a token from a given string.
// KeyFunc will receive the parsed token and should return the key for validating.
func Parse(tokenString string, keyFunc Keyfunc) (*Token, error) {
	return new(Parser).Parse(tokenString, keyFunc)
}

// ParseWithClaims parses a token with a given set of claims.
func ParseWithClaims(tokenString string, claims Claims, keyFunc Keyfunc) (*Token, error) {
	return new(Parser).ParseWithClaims(tokenString, claims, keyFunc)
}

// EncodeSegment en codes JWT specific base64url encoding with suffix '=' padding stripped.
func EncodeSegment(seg []byte) string {
	return strings.TrimRight(base64.URLEncoding.EncodeToString(seg), "=")
}

// DecodeSegment decodes a JWT specific base64url encoding with suffix '=' padding added back if not present.
func DecodeSegment(seg string) ([]byte, error) {
	if l := len(seg) % 4; l > 0 {
		seg += strings.Repeat("=", 4-l)
	}
	return base64.URLEncoding.DecodeString(seg)
}
