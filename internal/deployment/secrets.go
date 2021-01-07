package deployment

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/krane/krane/internal/constants"
	"github.com/krane/krane/internal/store"
)

type Secret struct {
	Deployment string `json:"deployment"`
	Key        string `json:"key"`
	Value      string `json:"value"`
	Alias      string `json:"alias"`
}

// AddSecret adds a secret to a deployment. Secrets are injected to the container during the container 'run' step.
// When a secret is created, an alias is returned and can be used to reference the secret in the `deployment.json`
// ie. SECRET_TOKEN=@secret-token (@secret-token was returned and how you reference the value for SECRET_TOKEN)
func AddSecret(deployment, key, value string) (*Secret, error) {
	if !isValidSecretKey(key) {
		return &Secret{}, fmt.Errorf("invalid secret name %s", key)
	}

	secret := &Secret{
		Deployment: deployment,
		Key:        key,
		Value:      value,
		Alias:      formatSecretAlias(key),
	}

	collection := getSecretsCollectionName(deployment)
	bytes, _ := secret.SerializeSecret()
	err := store.Client().Put(collection, secret.Key, bytes)
	if err != nil {
		return nil, err
	}

	return secret, nil
}

// DeleteSecret deletes a deployment secret
func DeleteSecret(deployment, key string) error {
	collection := getSecretsCollectionName(deployment)
	return store.Client().Remove(collection, key)
}

// CreateSecretsCollection creates secrets collection for a deployment
func CreateSecretsCollection(deployment string) error {
	collection := getSecretsCollectionName(deployment)
	return store.Client().CreateCollection(collection)
}

// DeleteCollection deletes secrets collection for a deployment
func DeleteSecretsCollection(deployment string) error {
	collection := getSecretsCollectionName(deployment)
	return store.Client().DeleteCollection(collection)
}

// GetAll returns all secrets for a deployment
func GetAllSecrets(deployment string) ([]*Secret, error) {
	collection := getSecretsCollectionName(deployment)
	bytes, err := store.Client().GetAll(collection)
	if err != nil {
		return make([]*Secret, 0), err
	}

	secrets := make([]*Secret, 0)
	for _, secret := range bytes {
		var s Secret
		err := json.Unmarshal(secret, &s)
		if err != nil {
			return make([]*Secret, 0), err
		}
		secrets = append(secrets, &s)
	}

	return secrets, nil
}

// GetAllSecretsRedacted returns all deployment secrets with <redacted> a their value
func GetAllSecretsRedacted(deployment string) []Secret {
	plainSecrets, _ := GetAllSecrets(deployment)
	redactedSecrets := make([]Secret, 0)
	for _, secret := range plainSecrets {
		secret.Redact()
		redactedSecrets = append(redactedSecrets, *secret)
	}
	return redactedSecrets
}

// GetSecret returns a deployment secret if it exists
func GetSecret(deployment, key string) (*Secret, error) {
	collection := getSecretsCollectionName(deployment)
	bytes, err := store.Client().Get(collection, key)
	if err != nil {
		return nil, err
	}

	if bytes == nil {
		return nil, fmt.Errorf("secret with key %s not found for deployment %s", key, deployment)
	}

	var s *Secret
	_ = json.Unmarshal(bytes, &s)

	return s, nil
}

// Redact masks the value for a secret
func (s *Secret) Redact() { s.Value = "<redacted>" }

func formatSecretAlias(key string) string {
	asLowerCase := strings.ToUpper(key)
	asUnderScore := strings.ReplaceAll(asLowerCase, "-", "_")
	return fmt.Sprintf("@%s", asUnderScore)
}

func isValidSecretKey(secret string) bool {
	if len(secret) <= 1 || len(secret) > 50 {
		return false
	}

	startsWithLetter := "[a-zA-Z0-9]"
	allowedCharacters := "[a-zA-Z0-9_-]"
	endWithLowerCaseAlphanumeric := "[a-zA-Z0-9]"

	matchers := fmt.Sprintf(`^%s%s*%s$`, // ^[a-zA-z0-9][a-zA-z0-9_-]*[a-zA-Z0-9]$
		startsWithLetter,
		allowedCharacters,
		endWithLowerCaseAlphanumeric)

	match := regexp.MustCompile(matchers)
	return match.MatchString(secret)
}

func (s Secret) SerializeSecret() ([]byte, error) { return json.Marshal(s) }

func getSecretsCollectionName(deployment string) string {
	return strings.ToLower(fmt.Sprintf("%s-%s", deployment, constants.SecretsCollectionName))
}
