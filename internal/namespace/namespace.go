package namespace

import (
	"github.com/biensupernice/krane/internal/collection"
	"github.com/biensupernice/krane/internal/kranecfg"
	"github.com/biensupernice/krane/internal/store"
)

func Exist(namespace string) bool {
	deployments, err := store.Instance().GetAll(collection.Deployments)
	if err != nil {
		return false
	}

	found := false
	for _, deployment := range deployments {
		var d kranecfg.KraneConfig
		err := store.Deserialize(deployment, &d)
		if err != nil {
			return false
		}

		if namespace == d.Name {
			found = true
		}
	}

	if !found {
		return false
	}

	return true
}
