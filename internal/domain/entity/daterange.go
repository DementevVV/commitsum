package entity

import (
	"fmt"
	"time"
)

// DateRange represents a date range selection.
type DateRange struct {
	StartDate string
	EndDate   string
	Label     string
}

// DateRangePreset represents a predefined date range option.
type DateRangePreset struct {
	Key   string
	Label string
}

// DateRangePresets contains available date range presets.
var DateRangePresets = []DateRangePreset{
	{Key: "today", Label: "Today"},
	{Key: "yesterday", Label: "Yesterday"},
	{Key: "week", Label: "Last 7 days"},
	{Key: "month", Label: "Last 30 days"},
	{Key: "custom", Label: "Custom date"},
}

// GetDateRange returns start and end dates for a preset.
func GetDateRange(preset string) DateRange {
	now := time.Now()
	today := now.Format("2006-01-02")

	switch preset {
	case "today":
		return DateRange{
			StartDate: today,
			EndDate:   today,
			Label:     "Today",
		}
	case "yesterday":
		yesterday := now.AddDate(0, 0, -1).Format("2006-01-02")
		return DateRange{
			StartDate: yesterday,
			EndDate:   yesterday,
			Label:     "Yesterday",
		}
	case "week":
		weekAgo := now.AddDate(0, 0, -7).Format("2006-01-02")
		return DateRange{
			StartDate: weekAgo,
			EndDate:   today,
			Label:     "Last 7 days",
		}
	case "month":
		monthAgo := now.AddDate(0, 0, -30).Format("2006-01-02")
		return DateRange{
			StartDate: monthAgo,
			EndDate:   today,
			Label:     "Last 30 days",
		}
	default:
		return DateRange{
			StartDate: today,
			EndDate:   today,
			Label:     "Custom",
		}
	}
}

// FormatDateDisplay formats date for display.
func FormatDateDisplay(startDate, endDate string) string {
	if startDate == endDate {
		return startDate
	}
	return fmt.Sprintf("%s → %s", startDate, endDate)
}
