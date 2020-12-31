package response

import (
	"encoding/json"
	"net/http"
)

// HTTPOk writes http response code 200
func HTTPOk(w http.ResponseWriter, data interface{}) {
	payload, _ := json.Marshal(data)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(payload)
	return
}

// HTTPNoContent writes http response code 204
func HTTPNoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
	return
}

// HTTPCreated writes http response code 201
func HTTPCreated(w http.ResponseWriter) {
	w.WriteHeader(http.StatusCreated)
	return
}

// HTTPAccepted writes http response code 202
func HTTPAccepted(w http.ResponseWriter) {
	w.WriteHeader(http.StatusAccepted)
	return
}

// HTTPBad writes http response code 400
func HTTPBad(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	_, _ = w.Write([]byte(err.Error()))
	return
}

// HTTPNotFound writes http response code 404
func HTTPNotFound(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	_, _ = w.Write([]byte(err.Error()))
	return
}
