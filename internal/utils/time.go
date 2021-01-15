package utils

import (
	"fmt"
	"time"
)

var (
	// OneYear is the unix time for 1 year
	OneYear = time.Now().Add(time.Minute * 525600).Unix()

	OneMinMs  = "60000"
	TwoMinMs  = "120000"
	FiveMinMs = "300000"
	TenMinMs  = "600000"
)

// UTCDateString returns the current date time in RFC3339 format
func UTCDateString() string {
	t := time.Now().Local()
	return t.Format(time.RFC3339)
}

// UnixToDate converts unix date into MM/DD/YYYY
func UnixToDate(u int64) string {
	t := time.Unix(u, 0)
	year, month, day := t.Date()
	return fmt.Sprintf("%02d/%d/%d", int(month), day, year)
}

// CalculateTimeRange returns the start & end date from a number of days ago
func CalculateTimeRange(daysAgo int) (string, string) {
	start := time.Now().AddDate(0, 0, -daysAgo).Format(time.RFC3339)
	end := time.Now().Local().Format(time.RFC3339)
	return start, end
}
