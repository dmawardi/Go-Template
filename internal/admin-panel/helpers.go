package adminpanel

import (
	"time"
)

// Helper function to format time.Time fields for forms
func formatTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format(time.RFC3339) // or another suitable format
}
