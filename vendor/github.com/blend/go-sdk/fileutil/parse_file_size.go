package fileutil

import (
	"strconv"
	"strings"
)

// ParseFileSize parses a file size
func ParseFileSize(fileSizeValue string) int64 {
	if len(fileSizeValue) == 0 {
		return 0
	}

	if len(fileSizeValue) < 2 {
		val, err := strconv.Atoi(fileSizeValue)
		if err != nil {
			return 0
		}
		return int64(val)
	}

	units := strings.ToLower(fileSizeValue[len(fileSizeValue)-2:])
	value, err := strconv.ParseInt(fileSizeValue[:len(fileSizeValue)-2], 10, 64)
	if err != nil {
		return 0
	}
	switch units {
	case "tb":
		return value * Terrabyte
	case "gb":
		return value * Gigabyte
	case "mb":
		return value * Megabyte
	case "kb":
		return value * Kilobyte
	}
	fullValue, err := strconv.ParseInt(fileSizeValue, 10, 64)
	if err != nil {
		return 0
	}
	return fullValue
}
