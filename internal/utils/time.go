package utils

import (
	"fmt"
	"time"
)

var (
	// OneYear : unix time for 1 year
	OneYear = time.Now().Add(time.Minute * 525600).Unix()
)

// UTCDateString : current date time in RFC3339 format
func UTCDateString() string {
	t := time.Now().Local()
	return t.Format(time.RFC3339)
}

// UnixToDate : format unix date into MM/DD/YYYY
func UnixToDate(u int64) string {
	t := time.Unix(u, 0)
	year, month, day := t.Date()
	return fmt.Sprintf("%02d/%d/%d", int(month), day, year)
}
