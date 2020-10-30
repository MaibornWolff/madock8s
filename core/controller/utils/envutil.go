package utils

import "os"

func GetEnvTargetNamespace() string {
	n := os.Getenv("TARGET_NAMESPACE")
	if n == "" {
		return "default"
	} else {
		return n
	}
}
