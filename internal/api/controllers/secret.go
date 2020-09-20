package controllers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/biensupernice/krane/internal/api/status"
	"github.com/biensupernice/krane/internal/secrets"
)

func GetSecrets(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	deploymentName := params["name"]

	if deploymentName == "" {
		status.HTTPBad(w, errors.New("deployment name required"))
		return
	}

	secrets, err := secrets.GetAll(deploymentName)
	if err != nil {
		status.HTTPBad(w, err)
		return
	}

	for _, secret := range secrets {
		secret.Redact()
	}

	status.HTTPOk(w, secrets)
	return
}

func CreateSecret(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	namespace := params["name"]

	if namespace == "" {
		status.HTTPBad(w, errors.New("namespace required"))
		return
	}

	type secretBody struct {
		Key   string `json:"key" binding:"required"`
		Value string `json:"value" binding:"required"`
	}

	var body secretBody
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		status.HTTPBad(w, err)
		return
	}

	alias, err := secrets.Add(body.Key, body.Value, namespace)
	if err != nil {
		status.HTTPBad(w, err)
		return
	}

	status.HTTPOk(w, alias)
	return
}
