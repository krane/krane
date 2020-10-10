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
func Add(key, value, nspace string) (string, error) {
	if !isValidSecretKey(key) {
		return "", fmt.Errorf("invalid secret name %s", key)
	}

	if !namespace.Exist(nspace) {
		return "", fmt.Errorf("unable to find namespace %s", nspace)
	}

	s := &Secret{
		Namespace: nspace,
		Key:       key,
		Value:     value,
		Alias:     formatSecretAlias(key),
	}
	bytes, _ := s.serialize()
	collection := getNamespaceCollectionName(nspace)
	store.Instance().Put(collection, s.Alias, bytes)

	return s.Alias, nil
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

func Get(namespace, alias string) (*Secret, error) {
	collection := getNamespaceCollectionName(namespace)
	bytes, err := store.Instance().Get(collection, alias)

	var s *Secret
	json.Unmarshal(bytes, &s)

	return s, err
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
