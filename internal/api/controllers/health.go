package controllers

import (
	"net/http"
	"os"

	"github.com/biensupernice/krane/internal/api/response"
	time "github.com/biensupernice/krane/internal/utils"
)

// HealthCheck : get Krane server response
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	host, _ := os.Hostname()
	response.HTTPOk(w, struct {
		Status    string `json:"response"`
		Host      string `json:"host"`
		Timestamp string `json:"timestamp"`
	}{
		Status:    "Running",
		Host:      host,
		Timestamp: time.UTCDateString(),
	})
}
