package utils

import (
	"log"
	"os"
	"strconv"
)

func RequireEnv(key string) {
	_, found := os.LookupEnv(key)
	if !found {
		log.Fatalf("Missing required env %s", key)
	}
}

func EnvOrDefault(key string, fallback string) string {
	value, found := os.LookupEnv(key)
	if !found {
		log.Printf("%s not set, defaulting to %s", key, fallback)
		os.Setenv(key, fallback)
		return fallback
	}
	log.Printf("%s already set with value %s", key, value)
	return os.Getenv(key)
}

func GetUIntEnv(key string) uint {
	value, found := os.LookupEnv(key)
	if !found {
		return 0
	}
	v, _ := strconv.ParseUint(value, 10, 8)
	return uint(v)
}
