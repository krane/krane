package routes

import (
	"net/http"
	"os"

	"github.com/biensupernice/krane/api/utils"
	time "github.com/biensupernice/krane/internal/utils"
)

// GetServerStatus : get server status
func GetServerStatus(w http.ResponseWriter, r *http.Request) {
	host, _ := os.Hostname()
	utils.HTTPOk(w, struct {
		Status    string `json:"status"`
		Host      string `json:"host"`
		Timestamp string `json:"timestamp"`
	}{
		Status:    "Running",
		Host:      host,
		Timestamp: time.UTCDateString(),
	})
}
