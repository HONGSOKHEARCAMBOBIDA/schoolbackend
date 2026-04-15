package helper

import "time"

// FormatDate trims timezone and time part from date string.
// Example: "1990-05-01T00:00:00+07:00" → "1990-05-01"
func FormatDate(dateStr string) string {
	if dateStr == "" {
		return ""
	}

	// Try to parse full ISO format
	t, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		// Fallback for SQL datetime style
		t, err = time.Parse("2006-01-02 15:04:05", dateStr)
		if err != nil {
			return dateStr // fallback to original string
		}
	}

	return t.Format("2006-01-02")
}

func FormatTime(timeStr string) string {
	if timeStr == "" {
		return ""
	}

	t, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		t, err = time.Parse("2006-01-02 15:04:05", timeStr)
		if err != nil {
			return timeStr
		}
	}

	return t.Format("15:04:05")
}
