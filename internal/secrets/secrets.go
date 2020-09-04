package secrets

const SECRETS_COLLECTION = "secrets"

type Secret struct {
	Key            string
	Value          string
	DeploymentName string
	Alias          string
}

// Add : a secret to a deployment. Secrets are injected to the container during the container build step.
// When a secret is created, an alias is returned and can be used to reference the secret in the `krane.json`
// ie. SECRET_TOKEN=@secret-token (@secret-token was returned and how you reference the value for SECRET_TOKEN)
// func Add(key, value, deploymentName string, store storage.Storage) {
// 	alias := makeAlias(key)
//
// 	b, err := store.Get(SECRETS_COLLECTION, alias)
//
// 	s := &Secret{key, value, deploymentName, alias}
// 	bytes, _ := json.Marshal(s)
// 	storage.Put(SECRETS_COLLECTION, s.Alias, bytes)
// }
//
// func makeAlias(key string) string {
// 	asLowerCase := strings.ToLower(key)
// 	// asDashed := strings.ReplaceAll(asLowerCase, "_", "-")
//
// 	return asLowerCase
// }
