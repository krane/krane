package namespace

import (
	"github.com/biensupernice/krane/internal/constants"
	"github.com/biensupernice/krane/internal/deployment/config"
	"github.com/biensupernice/krane/internal/store"
)

func Exist(namespace string) bool {
	deployments, err := store.Instance().GetAll(constants.DeploymentsCollectionName)
	if err != nil {
		return false
	}

	for _, deployment := range deployments {
		var d config.Kconfig
		if err := store.Deserialize(deployment, &d); err != nil {
			return false
		}

		if namespace == d.Name {
			return true
		}
	}

	return false
}
