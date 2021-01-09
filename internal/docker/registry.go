package docker

import (
	"encoding/base64"
	"encoding/json"
	"os"

	"github.com/krane/krane/internal/constants"
)

type RegistryCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Base64RegistryCredentials returns base64 container registry credentials
func Base64RegistryCredentials() string {
	bytes, _ := json.Marshal(RegistryCredentials{
		Username: os.Getenv(constants.EnvDockerBasicAuthUsername),
		Password: os.Getenv(constants.EnvDockerBasicAuthPassword),
	})
	return base64.StdEncoding.EncodeToString(bytes)
}
