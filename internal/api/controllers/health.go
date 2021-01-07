package controllers

import (
	"net/http"
	"os"

	"github.com/krane/krane/internal/api/response"
	"github.com/krane/krane/internal/docker"
	"github.com/krane/krane/internal/utils"
)

// HealthCheck returns the health and status of the running Krane instance
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	host, _ := os.Hostname()
	response.HTTPOk(w, struct {
		Docker    bool   `json:"docker"`
		Host      string `json:"host"`
		Timestamp string `json:"timestamp"`
	}{
		Docker:    docker.Ping(),
		Host:      host,
		Timestamp: utils.UTCDateString(),
	})
}
