package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/biensupernice/krane/internal/api/status"
	"github.com/biensupernice/krane/internal/deployment/config"
	"github.com/biensupernice/krane/internal/deployment/service"
)

func SaveDeployment(w http.ResponseWriter, r *http.Request) {
	var cfg config.Config

	if err := json.NewDecoder(r.Body).Decode(&cfg); err != nil {
		status.HTTPBad(w, err)
		return
	}

	if err := cfg.Save(); err != nil {
		status.HTTPBad(w, err)
		return
	}

	if err := service.StartDeployment(cfg); err != nil {
		status.HTTPBad(w, err)
		return
	}

	status.HTTPOk(w, cfg)
	return
}

func DeleteDeployment(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	name := params["name"]

	if err := config.Delete(name); err != nil {
		status.HTTPBad(w, err)
		return
	}

	cfg, err := config.Get(name)
	if err != nil {
		status.HTTPBad(w, err)
		return
	}

	_ = service.DeleteDeployment(cfg)

	status.HTTPOk(w, nil)
	return
}

func GetDeployment(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	name := params["name"]

	cfg, err := config.Get(name)
	if err != nil {
		status.HTTPBad(w, err)
		return
	}

	status.HTTPOk(w, cfg)
	return
}

func GetDeployments(w http.ResponseWriter, r *http.Request) {
	deployments, err := config.GetAll()
	if err != nil {
		status.HTTPBad(w, err)
		return
	}

	status.HTTPOk(w, deployments)
	return
}
