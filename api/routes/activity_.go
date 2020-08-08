package routes

import (
	"net/http"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/biensupernice/krane/api/utils"
	"github.com/biensupernice/krane/internal/service/activity"
)

func GetRecentActivity(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	daysAgo := query.Get("days_ago")

	//  Default to 1 day ago
	if daysAgo == "" {
		daysAgo = "1"
	}

	logrus.Info("Getting activities")
	daysAgoNum, _ := strconv.Atoi(daysAgo)
	start := time.Now().Add(time.Duration(-24*daysAgoNum) * time.Hour).Format(time.RFC3339)
	end := time.Now().Local().Format(time.RFC3339)

	recentActivity, err := activity.GetInRange(start, end)
	if err != nil {
		utils.HTTPBad(w, err)
		return
	}

	utils.HTTPOk(w, recentActivity)
	return
}
