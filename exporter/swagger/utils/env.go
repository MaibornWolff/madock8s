package utils

import (
	"os"
)

func GetSwaggerJsonURLFromEnv() string {
	if value, ok := os.LookupEnv("SWAGGER_JSON_URL"); ok {
		return string(value)
	} else {
		return value
	}
}

func GetSwaggerBaseURLFromEnv() string {
	if value, ok := os.LookupEnv("BASE_URL"); ok {
		return string(value)
	} else {
		return value
	}
}

func GetSwaggerJSONFromEnv() string {
	if value, ok := os.LookupEnv("SWAGGER_JSON"); ok {
		return string(value)
	} else {
		return value
	}
}

func GetSwaggerPortFromEnv() string {
	if value, ok := os.LookupEnv("SWAGGER_PORT"); ok {
		return string(value)
	} else {
		return value
	}
}
