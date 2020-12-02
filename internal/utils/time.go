package utils

import (
	"crypto/rand"
	"fmt"
	"io"
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

// MakeIdentifier : create a short unique identifier
func MakeIdentifier() string {
	b := make([]byte, 12)
	_, err := io.ReadFull(rand.Reader, b)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%x", b)
}
