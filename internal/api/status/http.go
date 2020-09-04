package status

import (
	"encoding/json"
	"net/http"
)

// HTTPResponse :
type HTTPResponse struct {
	Success bool        `json:"success"`
	Code    uint        `json:"code"`
	Data    interface{} `json:"data"`
}

// HTTPOk : response with status code 200
func HTTPOk(w http.ResponseWriter, data interface{}) {
	payload, _ := json.Marshal(&HTTPResponse{
		Success: true,
		Code:    http.StatusOK,
		Data:    data,
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(payload)
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
	payload, _ := json.Marshal(&HTTPResponse{
		Success: false,
		Code:    http.StatusBadRequest,
		Data:    err.Error(),
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	w.Write(payload)
	return
}
