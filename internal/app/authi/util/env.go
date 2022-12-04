package util

import (
	"fmt"
	"os"
	"strconv"
)

func GetEnvWithFallback(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func GetEnv(key string) (string, error) {
	if value, ok := os.LookupEnv(key); ok {
		return value, nil
	}
	return "", fmt.Errorf("environment variable %s not found", key)
}

func GetEnvIntWithFallback(key string, fallback int) (int, error) {
	if value, ok := os.LookupEnv(key); ok {
		return strconv.Atoi(value)
	}
	return fallback, nil
}
