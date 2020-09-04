package routes

import (
	"net/http"
	"strconv"
	"time"

	"github.com/biensupernice/krane/internal/activity"
	"github.com/biensupernice/krane/internal/api/status"
)

// GetRecentActivity : get recent activity
func GetRecentActivity(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	daysAgo := query.Get("days_ago")

	//  Defaults to 1 day ago
	if daysAgo == "" {
		daysAgo = "1"
	}

	daysAgoNum, _ := strconv.Atoi(daysAgo)
	start, end := calculateTimeRange(daysAgoNum)

	recentActivity, err := activity.GetInRange(start, end)
	if err != nil {
		status.HTTPBad(w, err)
		return
	}

	status.HTTPOk(w, recentActivity)
	return
}

func calculateTimeRange(daysAgo int) (string, string) {
	start := time.Now().Add(time.Duration(-24*daysAgo) * time.Hour).Format(time.RFC3339)
	end := time.Now().Local().Format(time.RFC3339)

	return start, end
}
