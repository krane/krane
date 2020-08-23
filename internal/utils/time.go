package utils

import "time"

// UTCDateString : Get the current date time in RFC3339 format
func UTCDateString() string {
	t := time.Now().Local()
	return t.Format(time.RFC3339)
}
