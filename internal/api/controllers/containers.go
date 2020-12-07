package controllers

import (
	"errors"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/biensupernice/krane/internal/api/response"
	"github.com/biensupernice/krane/internal/deployment/container"
	"github.com/biensupernice/krane/internal/deployment/namespace"
)

// GetContainers : gets all containers for a deployment
func GetContainers(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	name := params["name"]

	if name == "" {
		response.HTTPBad(w, errors.New("deployment name required"))
		return
	}

	if !namespace.Exist(name) {
		response.HTTPBad(w, errors.New("deployment does not exist"))
		return
	}

	containers, err := container.GetContainersByDeployment(name)
	if err != nil {
		response.HTTPBad(w, err)
		return
	}

	response.HTTPOk(w, containers)
	return
}
