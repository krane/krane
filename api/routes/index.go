package routes

import (
	"net/http"
	"os"

	api_utils "github.com/biensupernice/krane/api/utils"
	"github.com/biensupernice/krane/internal/utils"
)

func IndexRoute(w http.ResponseWriter, r *http.Request) {
	host, _ := os.Hostname()
	api_utils.HTTPOk(w, struct {
		Host      string `json:"host"`
		Timestamp string `json:"timestamp"`
		Path      string `json:"path"`
		Env       string `json:"env"`
		Port      string `json:"port"`
		Domain    string `json:"domain"`
	}{
		Host:      host,
		Timestamp: utils.UTCDateString(),
		Path:      "/",
		Env:       os.Getenv("ENV"),
		Port:      os.Getenv("PORT"),
		Domain:    os.Getenv("DOMAIN"),
	})
}
