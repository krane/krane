package controllers

import (
	"net/http"

	"github.com/krane/krane/internal/api/response"
)

// RootPath returns a plain-text response and 200 OK
func RootPath(w http.ResponseWriter, _ *http.Request) {
	response.HTTPOk(w, "Krane")
}
