package secrets

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/biensupernice/krane/internal/constants"
	"github.com/biensupernice/krane/internal/deployment/namespace"
	"github.com/biensupernice/krane/internal/store"
)

type Secret struct {
	Namespace string `json:"namespace"`
	Key       string `json:"key"`
	Value     string `json:"value"`
	Alias     string `json:"alias"`
}

// Add : a secret to a deployment. SecretsCollectionName are injected to the container during the container run step.
// When a secret is created, an alias is returned and can be used to reference the secret in the `krane.json`
// ie. SECRET_TOKEN=@secret-token (@secret-token was returned and how you reference the value for SECRET_TOKEN)
func Add(deploymentName, key, value string) (*Secret, error) {
	if !isValidSecretKey(key) {
		return &Secret{}, fmt.Errorf("invalid secret name %s", key)
	}

	if !namespace.Exist(deploymentName) {
		return nil, fmt.Errorf("unable to find namespace %s", deploymentName)
	}

	s := &Secret{
		Namespace: deploymentName,
		Key:       key,
		Value:     value,
		Alias:     formatSecretAlias(key),
	}

	bytes, _ := s.serialize()
	collection := getNamespaceCollectionName(deploymentName)
	err := store.Instance().Put(collection, s.Key, bytes)
	if err != nil {
		return nil, err
	}

	return s, nil
}

func Delete(namespace, key string) error {
	collection := getNamespaceCollectionName(namespace)
	return store.Instance().Remove(collection, key)
}

func CreateCollection(namespace string) error {
	collection := getNamespaceCollectionName(namespace)
	return store.Instance().CreateCollection(collection)
}

func DeleteCollection(namespace string) error {
	collection := getNamespaceCollectionName(namespace)
	return store.Instance().DeleteCollection(collection)
}

func GetAll(namespace string) ([]*Secret, error) {
	collection := getNamespaceCollectionName(namespace)
	bytes, err := store.Instance().GetAll(collection)
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

func GetAllRedacted(namespace string) []Secret {
	plainSecrets, _ := GetAll(namespace)
	redactedSecrets := make([]Secret, 0)
	for _, secret := range plainSecrets {
		secret.Redact()
		redactedSecrets = append(redactedSecrets, *secret)
	}
	return redactedSecrets
}

func Get(namespace, key string) (*Secret, error) {
	collection := getNamespaceCollectionName(namespace)
	bytes, err := store.Instance().Get(collection, key)

	if err != nil {
		return nil, err
	}

	if bytes == nil {
		return nil, fmt.Errorf("secret with key %s not found", key)
	}

	var s *Secret
	_ = json.Unmarshal(bytes, &s)

	return s, nil
}

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

func (s Secret) serialize() ([]byte, error) { return json.Marshal(s) }

func getNamespaceCollectionName(namespace string) string {
	return strings.ToLower(fmt.Sprintf("%s-%s", namespace, constants.SecretsCollectionName))
}
