package utils

import "time"

func UTCDateString() string {
	t := time.Now().Local()
	return t.Format(time.RFC3339)
}
