package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/biensupernice/krane/internal/api/status"
	"github.com/biensupernice/krane/internal/kranecfg"
)

// SaveDeployment :
func SaveDeployment(w http.ResponseWriter, r *http.Request) {
	var cfg kranecfg.KraneConfig
	err := json.NewDecoder(r.Body).Decode(&cfg)
	if err != nil {
		status.HTTPBad(w, err)
		return
	}

	err = cfg.Save()
	if err != nil {
		status.HTTPBad(w, err)
		return
	}

	status.HTTPOk(w, cfg)
	return
}

// DeleteDeployment :
func DeleteDeployment(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	name := params["name"]

	err := kranecfg.Delete(name)
	if err != nil {
		status.HTTPBad(w, err)
		return
	}

	status.HTTPOk(w, nil)
	return
}

// GetDeployment : get a deployment
func GetDeployment(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	name := params["name"]

	// Find deployment
	cfg, err := kranecfg.Get(name)
	if err != nil {
		status.HTTPBad(w, err)
		return
	}

	status.HTTPOk(w, cfg)
	return
}

// GetDeployments : get all deployments
func GetDeployments(w http.ResponseWriter, r *http.Request) {
	deployments, err := kranecfg.GetAll()
	if err != nil {
		status.HTTPBad(w, err)
		return
	}

	status.HTTPOk(w, deployments)
	return
}
