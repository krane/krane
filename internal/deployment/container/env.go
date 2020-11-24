package container

import (
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/biensupernice/krane/internal/deployment/kconfig"
	"github.com/biensupernice/krane/internal/secrets"
)

// from Kconfig to Docker environment variables string formatted list
func fromKconfigDockerEnvList(cfg kconfig.Kconfig) []string {
	envs := make([]string, 0)

	// kconfig environment variables
	for k, v := range cfg.Env {
		envs = append(envs, fmt.Sprintf("%s=%s", k, v))
	}

	// resolve secret by alias
	for key, alias := range cfg.Secrets {
		secret, err := secrets.Get(cfg.Name, alias)
		if err != nil || secret == nil {
			logrus.Debugf("Unable to get resolve secret for %s with alias @%s", cfg.Name, alias)
			continue
		}
		envs = append(envs, fmt.Sprintf("%s=%s", key, secret.Value))
	}

	return envs
}

// from Docker environment variables string list to environment variable map
func fromDockerToEnvMap(envs []string) map[string]string {
	envMap := make(map[string]string, 0)
	for _, env := range envs {
		keyValue := strings.Split(env, "=")

		key := keyValue[0]
		value := keyValue[1]

		envMap[key] = value
	}

	return envMap
}
