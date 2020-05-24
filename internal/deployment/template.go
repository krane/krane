package deployment

import (
	"encoding/json"
	"errors"

	"github.com/biensupernice/krane/internal/logger"
	"github.com/biensupernice/krane/internal/store"
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
	ContainerPort string `json:"container_port"`
	HostPort      string `json:"host_port"`
}

// SaveTemplate : save deployment template to datastore
func SaveTemplate(t *Template) (err error) {
	SetTemplateDefaults(t)

	if t.Name == "" {
		return errors.New("Unable to save template, missing field `name`")
	}

	// Converts template to bytes
	tBytes, err := json.Marshal(t)
	if err != nil {
		logger.Debugf("Unable to convert the deployment template into bytes - %s", err.Error())
		return err
	}

	// Save deployment to the datastore
	store.Put(store.TemplatesBucket, t.Name, tBytes)
	logger.Debugf("Template %s saved succesfuly", t.Name)

	return
}

// FindTemplate : from datastore by template name
func FindTemplate(name string) *Template {
	// Returns bytes
	tBytes := store.Get(store.TemplatesBucket, name)

	// Unmarshal bytes into template
	var t Template
	json.Unmarshal(tBytes, &t)

	return &t
}

// FindAllTemplates : from datastore
func FindAllTemplates() []Template {
	deploymentData := store.GetAll(store.TemplatesBucket)

	var deployments []Template
	for v := 0; v < len(deploymentData); v++ {
		var t Template
		err := json.Unmarshal(*deploymentData[v], &t)
		if err != nil {
			logger.Debugf("Unable to parse deployment [skipping] - %s", err.Error())
			continue
		}
		deployments = append(deployments, t)
	}

	return deployments
}

// ParseTemplate : validates a template
func ParseTemplate(t []byte) *Template {
	var tmpl Template
	err := json.Unmarshal(t, &tmpl)
	if err != nil {
		logger.Debugf("Unable to parse deployment template- %s", err.Error())
	}

	// Compare with a zero value composite literal because all fields are comparable
	// And check for empty name and image
	if tmpl == (Template{}) || tmpl.Name == "" || tmpl.Config.Image == "" {
		logger.Debug("Deployment template is missing values")
		return nil
	}

	logger.Debugf("%v", tmpl)

	return &tmpl
}

// SetTemplateDefaults :
func SetTemplateDefaults(t *Template) {
	const (
		DefaultRegistry = "docker.io"
	)

	if t.Config.Registry == "" {
		t.Config.Registry = DefaultRegistry
	}
}
