package adminpanel

import (
	"strconv"
	"time"
)

// Helper function to format time.Time fields for forms
func formatTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format(time.RFC3339) // or another suitable format
}

// Helper function to convert list of strings to list of ints
func convertStringSliceToIntSlice(stringSlice []string) ([]int, error) {
	intSlice := make([]int, 0, len(stringSlice)) // Create a slice of ints with the same length

	for _, str := range stringSlice {
		num, err := strconv.Atoi(str) // Convert string to int
		if err != nil {
			return nil, err
		}
		intSlice = append(intSlice, num) // Append the converted int to the slice
	}
	return intSlice, nil
}
