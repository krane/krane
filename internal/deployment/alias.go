package deployment

import (
	"github.com/biensupernice/krane/internal/storage"
)

var (
	AliasCollectionName = "alias"
)

func (d *Deployment) UpdateAlias(props map[string]string) error {
	alias := props["alias"]
	return storage.Put(AliasCollectionName, d.Spec.Name, []byte(alias))
}

func (d *Deployment) DeleteAlias(props map[string]string) error {
	return storage.Remove(AliasCollectionName, d.Spec.Name)
}

func GetDeploymentAlias(deploymentName string) (string, error) {
	bytes, err := storage.Get(AliasCollectionName, deploymentName)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}
