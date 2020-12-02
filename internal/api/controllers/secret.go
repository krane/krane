package controllers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/biensupernice/krane/internal/api/status"
	"github.com/biensupernice/krane/internal/secrets"
)

// GetSecrets : get deployment secrets
func GetSecrets(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	deploymentName := params["name"]

	if deploymentName == "" {
		status.HTTPBad(w, errors.New("deployment name required"))
		return
	}

	redactedSecrets := secrets.GetAllRedacted(deploymentName)
	status.HTTPOk(w, redactedSecrets)
	return
}

// CreateSecret : create a deployment secret
func CreateSecret(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	deploymentName := params["name"]

	if deploymentName == "" {
		status.HTTPBad(w, errors.New("deployment name required"))
		return
	}

	type secretBody struct {
		Key   string `json:"key" binding:"required"`
		Value string `json:"value" binding:"required"`
	}

	var body secretBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		status.HTTPBad(w, err)
		return
	}

	s, err := secrets.Add(deploymentName, body.Key, body.Value)
	if err != nil {
		status.HTTPBad(w, err)
		return
	}

	s.Redact()

	status.HTTPOk(w, s)
	return
}

// DeleteSecret : delete a deployment secret
func DeleteSecret(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	deploymentName := params["name"]
	key := params["key"]

	if deploymentName == "" {
		status.HTTPBad(w, errors.New("deployment name required"))
		return
	}

	if key == "" {
		status.HTTPBad(w, errors.New("key required"))
		return
	}

	if err := secrets.Delete(deploymentName, key); err != nil {
		status.HTTPBad(w, err)
		return
	}

	status.HTTPOk(w, nil)
	return
}
