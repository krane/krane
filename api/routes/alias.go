package routes

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/biensupernice/krane/api/utils"
	"github.com/biensupernice/krane/internal/service/deployment"
)

// UpdateDeploymentAlias : update an alias
func UpdateDeploymentAlias(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	deploymentName := params["name"]

	type Body struct {
		Alias string `json:"alias" binding:"required"`
	}

	// Decode request
	var body Body
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		utils.HTTPBad(w, err)
		return
	}

	// Find deployment
	_, err = deployment.GetDeployment(deploymentName)
	if err != nil {
		utils.HTTPBad(w, err)
		return
	}

	// jobs.UpdateDeploymentAlias(d, body.Alias)

	utils.HTTPAccepted(w)
	return
}

// DeleteDeploymentAlias : delete an alias
func DeleteDeploymentAlias(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	deploymentName := params["name"]

	// Find deployment
	_, err := deployment.GetDeployment(deploymentName)
	if err != nil {
		utils.HTTPBad(w, err)
		return
	}

	// jobs.DeleteDeploymentAlias(d)

	utils.HTTPAccepted(w)
	return
}
