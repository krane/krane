package spec

import (
	"encoding/json"
	"fmt"

	"github.com/biensupernice/krane/internal/logger"
	"github.com/biensupernice/krane/internal/store"
)

// Spec : a spec that defines a deployment
type Spec struct {
	Name   string `json:"name" binding:"required"`
	Config Config `json:"config" binding:"required"`
}

// Config : for a deployment spec
type Config struct {
	Registry      string `json:"registry"`
	Image         string `json:"image" binding:"required"`
	ContainerPort string `json:"container_port"`
	HostPort      string `json:"host_port"`
}

// Create : a spec
func (s Spec) Create() (err error) {
	spec := Find(s.Name)

	// Check if deployment already exist
	if spec != (Spec{}) {
		errMsg := fmt.Sprintf("Deployment %s already exist", spec.Name)
		logger.Debugf(errMsg)
		return
	}

	// Save to datastore
	return s.Save()
}

// Delete : spec from datastore
func (s Spec) Delete() (err error) { return store.Remove(store.SpecsBucket, s.Name) }

// Find : spec from datastore
func Find(deploymentName string) Spec {
	// Returns bytes
	tBytes := store.Get(store.SpecsBucket, deploymentName)

	// Unmarshal bytes into Spec
	var s Spec
	json.Unmarshal(tBytes, &s)

	return s
}

// FindAll : specs from datastore
func FindAll() []Spec {
	specsBytes := store.GetAll(store.SpecsBucket)

	specs := make([]Spec, 0)
	for v := 0; v < len(specsBytes); v++ {
		var s Spec
		err := json.Unmarshal(*specsBytes[v], &s)
		if err != nil {
			logger.Debugf("Unable to parse deployment [skipping] - %s", err.Error())
			continue
		}
		specs = append(specs, s)
	}

	return specs
}

// Save : spec to datastore
func (s Spec) Save() (err error) {
	sBytes, err := json.Marshal(s)
	if err != nil {
		logger.Debugf("Unable to convert the deployment spec into bytes - %s", err.Error())
		return
	}

	// Save spec to the datastore
	return store.Put(store.SpecsBucket, s.Name, sBytes)
}

// SetDefaults : for a spec
func (s *Spec) SetDefaults() {
	const (
		DefaultRegistry = "docker.io"
	)

	if s.Config.Registry == "" {
		s.Config.Registry = DefaultRegistry
	}
}
