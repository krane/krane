package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
)

func (cfg Config) Serialize() ([]byte, error) { return json.Marshal(cfg) }

func (cfg Config) validate() error {
	isValidName := cfg.validateName()
	if !isValidName {
		return errors.New("invalid name")
	}

	return nil
}

func (cfg *Config) applyDefaults() {
	if cfg.Registry == "" {
		cfg.Registry = "docker.io"
	}

	if cfg.Alias == nil {
		cfg.Alias = make([]string, 0)
	}

	if cfg.Secrets == nil {
		cfg.Secrets = make(map[string]string, 0)
	}

	if cfg.Env == nil {
		cfg.Env = make(map[string]string, 0)
	}

	if cfg.Volumes == nil {
		cfg.Volumes = make(map[string]string, 0)
	}

	if cfg.Ports == nil {
		cfg.Ports = make(map[string]string, 0)
	}

	if cfg.Tag == "" {
		cfg.Tag = "latest"
	}

	return
}

func (cfg Config) validateName() bool {
	startsWithLetter := "[a-z]"
	allowedCharacters := "[a-z0-9_-]"
	endWithLowerCaseAlphanumeric := "[0-9a-z]"
	characterLimit := "{1,}"

	matchers := fmt.Sprintf(`^%s%s*%s%s$`, // ^[a - z][a - z0 - 9_ -]*[0-9a-z]$
		startsWithLetter,
		allowedCharacters,
		endWithLowerCaseAlphanumeric,
		characterLimit)

	match := regexp.MustCompile(matchers)
	return match.MatchString(cfg.Name)
}
