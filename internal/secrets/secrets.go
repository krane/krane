package secrets

import (
	"encoding/json"
	"strings"

	"github.com/biensupernice/krane/internal/collection"
	"github.com/biensupernice/krane/internal/store"
)

type Secret struct {
	Namespace string
	Key       string
	Value     string
	Alias     string
}

// Add : a secret to a deployment. Secrets are injected to the container during the container run step.
// When a secret is created, an alias is returned and can be used to reference the secret in the `krane.json`
// ie. SECRET_TOKEN=@secret-token (@secret-token was returned and how you reference the value for SECRET_TOKEN)
func Add(key, value, namespace string, store store.Store) string {
	alias := generateSecretAlias(key)

	s := &Secret{namespace, key, value, alias}
	bytes, _ := s.serialize()

	store.Put(collection.SECRETS, s.Alias, bytes)

	return alias
}

func generateSecretAlias(key string) string {
	asLowerCase := strings.ToLower(key)
	// asDashed := strings.ReplaceAll(asLowerCase, "_", "-")

	return asLowerCase
}

func (s Secret) serialize() ([]byte, error) { return json.Marshal(s) }
