package docker

import (
	"encoding/base64"
	"encoding/json"
)

type RegistryCredentials struct {
	URL      string `json:"url"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// Base64RegistryCredentials returns base64 container registry credentials
func Base64RegistryCredentials(username string, password string) string {
	bytes, _ := json.Marshal(RegistryCredentials{
		Username: username,
		Password: password,
	})
	return base64.StdEncoding.EncodeToString(bytes)
}
