package activity

import (
	"encoding/json"

	"github.com/docker/distribution/uuid"
	"github.com/sirupsen/logrus"

	"github.com/biensupernice/krane/internal/storage"
	"github.com/biensupernice/krane/internal/utils"
	"github.com/biensupernice/krane/pkg/bbq"
)

type Activity struct {
	ID        string  `json:"activity_id"`
	CreatedAt string  `json:"created_at"`
	Job       bbq.Job `json:"job"`
	// Session
}

var (
	ActivityCollectionName = "activity"
)

func Capture(a *Activity) {
	a.ID = uuid.Generate().String()
	a.CreatedAt = utils.UTCDateString()

	// Convert activity struct to bytes
	bytes, err := json.Marshal(a)
	if err != nil {
		logrus.Error(err.Error())
		return
	}

	// The key is an RFC3339 encoded time and is required to do date range scans across the db
	// this way we can perform looks up for activity within a time range
	err = storage.Put(ActivityCollectionName, a.CreatedAt, bytes)
	if err != nil {
		logrus.Error(err.Error())
		return
	}

	logrus.Debugf("[%s] -> Job activity captured", a.Job.ID)
}

func GetInRange(minDate, maxDate string) ([]Activity, error) {
	// Find activity in the range
	bytes, err := storage.GetInRange(ActivityCollectionName, minDate, maxDate)
	if err != nil {
		return make([]Activity, 0), err
	}

	recentActivity := make([]Activity, 0)
	for _, activity := range bytes {
		var a Activity
		err := json.Unmarshal(activity, &a)
		if err != nil {
			return recentActivity, err
		}

		recentActivity = append(recentActivity, a)
	}
	return recentActivity, nil
}
