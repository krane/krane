package status

import (
	"encoding/json"
	"net/http"
)

// HTTPOk : response with status code 200
func HTTPOk(w http.ResponseWriter, data interface{}) {
	payload, _ := json.Marshal(data)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(payload)
	return
}

// HTTPCreated : response with status code 201
func HTTPCreated(w http.ResponseWriter) {
	w.WriteHeader(http.StatusCreated)
	return
}

// HTTPAccepted : response with status code 202
func HTTPAccepted(w http.ResponseWriter) {
	w.WriteHeader(http.StatusAccepted)
	return
}

// HTTPBad : response with code 400
func HTTPBad(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	_, _ = w.Write([]byte(err.Error()))
	return
}
