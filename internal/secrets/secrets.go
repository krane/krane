package secrets

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/biensupernice/krane/internal/collection"
	"github.com/biensupernice/krane/internal/namespace"
	"github.com/biensupernice/krane/internal/store"
)

type Secret struct {
	Namespace string `json:"namespace"`
	Key       string `json:"key"`
	Value     string `json:"value"`
	Alias     string `json:"alias"`
}

// Add : a secret to a deployment. Secrets are injected to the container during the container run step.
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
		Alias:     createAlias(key),
	}
	bytes, _ := s.serialize()
	collection := getNamespaceCollectionName(nspace)
	store.Instance().Put(collection, s.Alias, bytes)

	return s.Alias, nil
}

func Get(namespace string) ([]*Secret, error) {
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

func (s *Secret) Redact() { s.Value = "<redacted>" }

func createAlias(key string) string {
	asLowerCase := strings.ToUpper(key)
	asUnderScore := strings.ReplaceAll(asLowerCase, "-", "_")
	return fmt.Sprintf("@%s", asUnderScore)
}

func isValidSecretKey(secret string) bool {
	startsWithLetter := "[a-zA-z]"
	allowedCharacters := "[a-zA-z0-9_-]"
	endWithLowerCaseAlphanumeric := "[0-9a-zA-z	]"
	characterLimit := "{1,}"

	matchers := fmt.Sprintf(`^%s%s*%s%s$`, // ^[a - z][a - z0 - 9_ -]*[0-9a-z]$
		startsWithLetter,
		allowedCharacters,
		endWithLowerCaseAlphanumeric,
		characterLimit)

	match := regexp.MustCompile(matchers)
	return match.MatchString(secret)
}

func (s Secret) serialize() ([]byte, error) { return json.Marshal(s) }

func getNamespaceCollectionName(namespace string) string {
	return strings.ToLower(fmt.Sprintf("%s-%s", namespace, collection.SECRETS))
}
