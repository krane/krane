package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/krane/krane/internal/api/response"
	"github.com/krane/krane/internal/deployment"
)

// GetSecrets returns all secrets for a deployment
func GetSecrets(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	deploymentName := params["deployment"]

	if deploymentName == "" {
		response.HTTPBad(w, errors.New("deployment name required"))
		return
	}

	redactedSecrets := deployment.GetAllSecretsRedacted(deploymentName)
	response.HTTPOk(w, redactedSecrets)
	return
}

// CreateOrUpdateSecret saves a deployment secret
func CreateOrUpdateSecret(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	deploymentName := params["deployment"]

	if deploymentName == "" {
		response.HTTPBad(w, errors.New("deployment name required"))
		return
	}

	if !deployment.Exist(deploymentName) {
		response.HTTPBad(w, fmt.Errorf("unable to find deployment %s", deploymentName))
		return
	}

	type SecretRequest struct {
		Key   string `json:"key" binding:"required"`
		Value string `json:"value" binding:"required"`
	}

	var body SecretRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		response.HTTPBad(w, err)
		return
	}

	newSecret, err := deployment.AddSecret(deploymentName, body.Key, body.Value)
	if err != nil {
		response.HTTPBad(w, err)
		return
	}

	newSecret.Redact()

	response.HTTPOk(w, newSecret)
	return
}

// DeleteSecret removes a deployment secret
func DeleteSecret(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	deploymentName := params["deployment"]
	key := params["key"]

	if deploymentName == "" {
		response.HTTPBad(w, errors.New("deployment name required"))
		return
	}

	if key == "" {
		response.HTTPBad(w, errors.New("secret key required"))
		return
	}

	if err := deployment.DeleteSecret(deploymentName, key); err != nil {
		response.HTTPBad(w, err)
		return
	}

	response.HTTPNoContent(w)
	return
}
