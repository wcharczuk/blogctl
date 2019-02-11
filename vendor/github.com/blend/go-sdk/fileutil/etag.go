package fileutil

import (
	"crypto/md5"
	"encoding/hex"

	"github.com/blend/go-sdk/exception"
)

// ETag creates an etag for a given blob.
func ETag(contents []byte) (string, error) {
	hash := md5.New()
	_, err := hash.Write(contents)
	if err != nil {
		return "", exception.New(err)
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}
