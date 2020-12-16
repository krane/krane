package service

import (
	"github.com/biensupernice/krane/internal/constants"
	"github.com/biensupernice/krane/internal/deployment/config"
	"github.com/biensupernice/krane/internal/store"
)

// DeploymentExist : check if a deployment exists
func DeploymentExist(deploymentName string) bool {
	deployments, err := store.Client().GetAll(constants.DeploymentsCollectionName)
	if err != nil {
		return false
	}

	for _, deployment := range deployments {
		var d config.DeploymentConfig
		if err := store.Deserialize(deployment, &d); err != nil {
			return false
		}

		if deploymentName == d.Name {
			return true
		}
	}

	return false
}
