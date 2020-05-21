package deployment

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/biensupernice/krane/internal/data"
)

// Template :
type Template struct {
	Name   string `json:"name" binding:"required"`
	Config Config `json:"config" binding:"required"`
}

// Config :
type Config struct {
	Registry      string `json:"registry"`
	Image         string `json:"image" binding:"required"`
	Tag           string `json:"tag"`
	ContainerPort string `json:"container_port"`
	HostPort      string `json:"host_port"`
}

// SaveDeployment : to datastore
func SaveDeployment(t *Template) (err error) {
	SetTemplateDefaults(t)

	if t.Name == "" {
		return errors.New("Unable to save template, missing field `name`")
	}

	// Converts template to bytes
	tBytes, err := json.Marshal(t)
	if err != nil {
		return err
	}

	data.Put(data.DeploymentsBucket, t.Name, tBytes)

	return
}

// FindDeployment : from datastore by deployment name
func FindDeployment(name string) *Template {
	// Returns bytes
	tBytes := data.Get(data.DeploymentsBucket, name)

	// Unmarshal bytes into template
	var t Template
	json.Unmarshal(tBytes, &t)

	return &t
}

// FindAllDeployments : from datastore
func FindAllDeployments() []Template {
	deploymentData := data.GetAll(data.DeploymentsBucket)

	var deployments []Template
	for v := 0; v < len(deploymentData); v++ {
		var t Template
		err := json.Unmarshal(*deploymentData[v], &t)
		if err != nil {
			log.Printf("Unable to parse deployment [skipping] - %s", err.Error())
			continue
		}
		deployments = append(deployments, t)
	}

	return deployments
}

// ParseTemplate : validates a template
func ParseTemplate(t []byte) *Template {
	var tmpl Template
	json.Unmarshal(t, &tmpl)

	if tmpl.Name == "" ||
		tmpl.Config.Image == "" {
		return nil
	}

	return &tmpl
}

// SetTemplateDefaults :
func SetTemplateDefaults(t *Template) {
	const (
		DefaultRegistry = "docker.io"
		DefaultTag      = "latest"
	)

	if t.Config.Registry == "" {
		t.Config.Registry = DefaultRegistry
	}

	if t.Config.Tag == "" {
		t.Config.Tag = DefaultTag
	}
}
