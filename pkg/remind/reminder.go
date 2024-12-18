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

func CalculateReminderTime(unit string, value int) (time.Duration, string) {
	switch unit {
	case "s":
		return time.Duration(value) * time.Second, timeBeautifier(value, unit)
	case "h":
		return time.Duration(value) * time.Hour, timeBeautifier(value, unit)
	case "d":
		return time.Duration(value) * 24 * time.Hour, timeBeautifier(value, unit)
	case "w":
		return time.Duration(value) * 7 * 24 * time.Hour, timeBeautifier(value, unit)
	case "m":
		return time.Duration(value) * 30 * 24 * time.Hour, timeBeautifier(value, unit)
	default:
		return 0, ""
	}
}

func timeBeautifier(duration int, format string) string {
	switch format {
	case "s":
		return getRightForm(duration, "секунды", "секунд", "секунду")
	case "h":
		return getRightForm(duration, "часа", "часов", "час")
	case "d":
		return getRightForm(duration, "дня", "дней", "день")
	case "w":
		return getRightForm(duration, "недели", "недель", "неделю")
	case "m":
		return getRightForm(duration, "месяца", "месяцев", "месяц")
	default:
		return "Неверный формат"
	}
}

func getRightForm(n int, form1, form2, form3 string) string {
	if n == 1 {
		return form3
	}
	if n >= 2 && n <= 5 {
		return form1
	}
	n = n % 100
	if n >= 11 && n <= 19 {
		return form2
	}
	n = n % 10
	if n == 1 {
		return form3
	}
	return form2
}
