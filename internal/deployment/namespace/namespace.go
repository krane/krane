package namespace

import (
	"github.com/biensupernice/krane/internal/constants"
	"github.com/biensupernice/krane/internal/deployment/config"
	"github.com/biensupernice/krane/internal/store"
)

// Exist : check if a deployment exists
func Exist(namespace string) bool {
	deployments, err := store.Client().GetAll(constants.DeploymentsCollectionName)
	if err != nil {
		return false
	}

	for _, deployment := range deployments {
		var d config.DeploymentConfig
		if err := store.Deserialize(deployment, &d); err != nil {
			return false
		}

		if namespace == d.Name {
			return true
		}
	}

	return false
}
