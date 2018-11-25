package aws

import (
	"strings"
	"time"
)

// RefStr returns a string reference.
func RefStr(str string) *string {
	if str == "" {
		return nil
	}
	return &str
}

// RefStrs returns a list of string references.
func RefStrs(inputs ...string) []*string {
	strs := make([]*string, len(inputs))
	for index := range inputs {
		strs[index] = &inputs[index]
	}
	return strs
}

// DerefStr safely dereferences a string.
func DerefStr(str *string) string {
	if str == nil {
		return ""
	}
	return *str
}

// RefTime returns a time.Time reference.
func RefTime(t time.Time) *time.Time {
	return &t
}

// DerefTime deferences a time.Time.
func DerefTime(t *time.Time) time.Time {
	if t == nil {
		return time.Time{}
	}
	return *t
}

// StripQuotes strips leading or trailing quotes.
func StripQuotes(v string) string {
	v = strings.TrimSpace(v)
	v = strings.TrimPrefix(v, "\"")
	v = strings.TrimSuffix(v, "\"")
	return v
}
