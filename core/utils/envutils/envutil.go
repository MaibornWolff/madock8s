package envutils

import (
	"os"
	"strings"
)

func GetDeletionStrategyFromEnv() string {
	if value, ok := os.LookupEnv("DELETION_STRATEGY"); ok {
		return strings.ToUpper(value)
	} else {
		return "UPDATE"
	}
}
