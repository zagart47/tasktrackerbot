package remind

import (
	"errors"
	"time"
)

func ValidateReminderDuration(unit string, value int) error {
	validUnits := map[string]bool{
		"s": true,
		"h": true,
		"d": true,
		"w": true,
		"m": true,
	}
	if !validUnits[unit] {
		return errors.New("invalid remind unit")
	}
	if value <= 0 {
		return errors.New("remind value must be positive")
	}
	return nil
}

func CalculateReminderTime(unit string, value int) time.Time {
	now := time.Now()
	switch unit {
	case "s":
		return now.Add(time.Duration(value) * time.Second)
	case "h":
		return now.Add(time.Duration(value) * time.Hour)
	case "d":
		return now.AddDate(0, 0, value)
	case "w":
		return now.AddDate(0, 0, value*7)
	case "m":
		return now.AddDate(0, value, 0)
	default:
		return now
	}
}
