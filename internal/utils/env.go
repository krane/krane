package utils

import (
	"log"
	"os"
	"strconv"
)

// RequireEnv : hard check if an environment variable is set
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

// EnvOrDefault : get an environment variable, returns a default value if env not found
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

// GetUIntEnv : get environment variable converted into an unsigned int
func GetUIntEnv(key string) uint {
	value, found := os.LookupEnv(key)
	if !found {
		return 0
	}
	v, _ := strconv.ParseUint(value, 10, 8)
	return uint(v)
}

// GetIntEnv : get environment variable converted into an int
func GetIntEnv(key string) int {
	value, found := os.LookupEnv(key)
	if !found {
		return 0
	}
	v, _ := strconv.ParseInt(value, 10, 8)
	return int(v)
}

// GetBoolEnv : get environment variable converted into a boolean, returns false if env not found
func GetBoolEnv(key string) bool {
	value, found := os.LookupEnv(key)
	if !found {
		return false
	}
	v, _ := strconv.ParseBool(value)
	return v
}
