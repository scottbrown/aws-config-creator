package setlist

import (
	"time"
)

func generateTimestamp() string {
	now := time.Now().UTC()

	return now.Format("2006-01-02T15:04:05 MST")
}
