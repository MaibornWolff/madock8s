package utils

import (
	"os"
	"strconv"
)

func INCLUDE_LABELED_ONLY() bool {
	raw := os.Getenv("INCLUDE_LABELED_ONLY")
	value, err := strconv.ParseBool(raw)
	if err != nil {
		return false
	}
	return value
}
