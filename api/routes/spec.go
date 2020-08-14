package routes

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/biensupernice/krane/api/utils"
	"github.com/biensupernice/krane/internal/service/spec"
	"github.com/biensupernice/krane/internal/storage"
)

// CreateSpec : create a spec
func CreateSpec(w http.ResponseWriter, r *http.Request) {
	var s spec.Spec
	err := json.NewDecoder(r.Body).Decode(&s)
	if err != nil {
		utils.HTTPBad(w, err)
		return
	}

	err = s.CreateSpec()
	if err != nil {
		utils.HTTPBad(w, err)
		return
	}

	createdSpec, _ := spec.GetOne(s.Name)
	utils.HTTPOk(w, createdSpec)
	return
}

// UpdateSpec : update a spec
func UpdateSpec(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	name := params["name"]

	var s spec.Spec
	err := json.NewDecoder(r.Body).Decode(&s)
	if err != nil {
		utils.HTTPBad(w, err)
		return
	}

	// Verify the Spec exist
	data, err := storage.Get(spec.Collection, name)
	if err != nil {
		utils.HTTPBad(w, err)
		return
	}

	if data == nil {
		utils.HTTPBad(w, errors.New("Spec not found"))
		return
	}

	// Update spec
	err = s.UpdateSpec(name)
	if err != nil {
		utils.HTTPBad(w, err)
		return
	}

	utils.HTTPCreated(w)
	return
}
