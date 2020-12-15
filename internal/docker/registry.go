package docker

import (
	"encoding/base64"
	"encoding/json"
	"os"

	"github.com/biensupernice/krane/internal/constants"
)

// Base64DockerRegistryCredentials : returns base64 docker registry credentials
func Base64DockerRegistryCredentials() string {
	bytes, _ := json.Marshal(struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{
		Username: os.Getenv(constants.EnvDockerBasicAuthUsername),
		Password: os.Getenv(constants.EnvDockerBasicAuthPassword),
	})
	return base64.StdEncoding.EncodeToString(bytes)
}
