package utils

import (
	"log"
	"os"
	"strconv"
)

// RequireEnv exits the program if environment vairables not set
func RequireEnv(key string) {
	value, found := os.LookupEnv(key)
	if !found {
		log.Fatalf("Missing required env %s", key)
	}

	if value == "" {
		log.Fatalf("Missing required env %s", key)
	}

	log.Printf("%s=%s", key, value)
}

// EnvOrDefault returns and environment variable or a default value
func EnvOrDefault(key string, fallback string) string {
	value, found := os.LookupEnv(key)
	if !found {
		log.Printf("%s not set, defaulting to %s", key, fallback)
		if err := os.Setenv(key, fallback); err != nil {
			log.Printf("Error setting %s, %v", key, err)
		}
		return fallback
	}
	log.Printf("%s=%s", key, value)
	return os.Getenv(key)
}

// UIntEnv returns the unsigned int environment variable or 0 if not found
func UIntEnv(key string) uint {
	value, found := os.LookupEnv(key)
	if !found {
		return 0
	}
	v, _ := strconv.ParseUint(value, 10, 8)
	return uint(v)
}

// IntEnv returns the int environment variable or 0 if not found
func IntEnv(key string) int {
	value, found := os.LookupEnv(key)
	if !found {
		return 0
	}
	v, _ := strconv.ParseInt(value, 10, 8)
	return int(v)
}

// BoolEnv return the boolean environment variables or false if not found
func BoolEnv(key string) bool {
	value, found := os.LookupEnv(key)
	if !found {
		return false
	}
	v, _ := strconv.ParseBool(value)
	return v
}
