package setlist

import (
	"time"
)

// generateTimestamp returns the current UTC timestamp formatted for logs.
func generateTimestamp() string {
	now := time.Now().UTC()

	return now.Format("2006-01-02T15:04:05 MST")
}
