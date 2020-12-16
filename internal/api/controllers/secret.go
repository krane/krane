package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/biensupernice/krane/internal/api/response"
	"github.com/biensupernice/krane/internal/deployment/secrets"
	"github.com/biensupernice/krane/internal/deployment/service"
)

// GetSecrets : get deployment secrets
func GetSecrets(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	deploymentName := params["name"]

	if deploymentName == "" {
		response.HTTPBad(w, errors.New("deployment name required"))
		return
	}

	redactedSecrets := secrets.GetAllRedacted(deploymentName)
	response.HTTPOk(w, redactedSecrets)
	return
}

// CreateSecret : create a deployment secret
func CreateSecret(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	deploymentName := params["name"]

	if deploymentName == "" {
		response.HTTPBad(w, errors.New("deployment name required"))
		return
	}

	if !service.DeploymentExist(deploymentName) {
		response.HTTPBad(w, fmt.Errorf("unable to find deployment %s", deploymentName))
		return
	}

	type secretBody struct {
		Key   string `json:"key" binding:"required"`
		Value string `json:"value" binding:"required"`
	}

	var body secretBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		response.HTTPBad(w, err)
		return
	}

	s, err := secrets.Add(deploymentName, body.Key, body.Value)
	if err != nil {
		response.HTTPBad(w, err)
		return
	}

	s.Redact()

	response.HTTPOk(w, s)
	return
}

// DeleteSecret : delete a deployment secret
func DeleteSecret(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	deploymentName := params["name"]
	key := params["key"]

	if deploymentName == "" {
		response.HTTPBad(w, errors.New("deployment name required"))
		return
	}

	if key == "" {
		response.HTTPBad(w, errors.New("key required"))
		return
	}

	if err := secrets.Delete(deploymentName, key); err != nil {
		response.HTTPBad(w, err)
		return
	}

	response.HTTPOk(w, nil)
	return
}
