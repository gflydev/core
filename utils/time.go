package utils

import (
	"time"
)

// FormatTime formats a time.Time value according to the specified layout.
//
// Parameters:
//   - t time.Time: The time value to format.
//   - layout string: The layout string to use for formatting (e.g., "2006-01-02").
//
// Returns:
//   - string: The formatted time string.
func FormatTime(t time.Time, layout string) string {
	return t.Format(layout)
}

// ParseTime parses a time string according to the specified layout.
//
// Parameters:
//   - value string: The time string to parse.
//   - layout string: The layout string to use for parsing (e.g., "2006-01-02").
//
// Returns:
//   - time.Time: The parsed time value.
//   - error: An error if the time string couldn't be parsed.
func ParseTime(value, layout string) (time.Time, error) {
	return time.Parse(layout, value)
}

// Now returns the current local time.
//
// Returns:
//   - time.Time: The current local time.
func Now() time.Time {
	return time.Now()
}

// UTC returns the current UTC time.
//
// Returns:
//   - time.Time: The current UTC time.
func UTC() time.Time {
	return time.Now().UTC()
}

// IsWeekend checks if the given time falls on a weekend (Saturday or Sunday).
//
// Parameters:
//   - t time.Time: The time to check.
//
// Returns:
//   - bool: True if the time falls on a weekend, false otherwise.
func IsWeekend(t time.Time) bool {
	day := t.Weekday()
	return day == time.Saturday || day == time.Sunday
}

// IsWeekday checks if the given time falls on a weekday (Monday through Friday).
//
// Parameters:
//   - t time.Time: The time to check.
//
// Returns:
//   - bool: True if the time falls on a weekday, false otherwise.
func IsWeekday(t time.Time) bool {
	return !IsWeekend(t)
}

// StartOfDay returns the start of the day (00:00:00) for the given time.
//
// Parameters:
//   - t time.Time: The time to get the start of the day for.
//
// Returns:
//   - time.Time: The start of the day for the given time.
func StartOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

// EndOfDay returns the end of the day (23:59:59.999999999) for the given time.
//
// Parameters:
//   - t time.Time: The time to get the end of the day for.
//
// Returns:
//   - time.Time: The end of the day for the given time.
func EndOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 999999999, t.Location())
}

// StartOfWeek returns the start of the week (Monday 00:00:00) for the given time.
//
// Parameters:
//   - t time.Time: The time to get the start of the week for.
//
// Returns:
//   - time.Time: The start of the week for the given time.
func StartOfWeek(t time.Time) time.Time {
	weekday := int(t.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	return StartOfDay(t.AddDate(0, 0, -weekday+1))
}

// EndOfWeek returns the end of the week (Sunday 23:59:59.999999999) for the given time.
//
// Parameters:
//   - t time.Time: The time to get the end of the week for.
//
// Returns:
//   - time.Time: The end of the week for the given time.
func EndOfWeek(t time.Time) time.Time {
	return EndOfDay(StartOfWeek(t).AddDate(0, 0, 6))
}

// StartOfMonth returns the start of the month (1st day 00:00:00) for the given time.
//
// Parameters:
//   - t time.Time: The time to get the start of the month for.
//
// Returns:
//   - time.Time: The start of the month for the given time.
func StartOfMonth(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
}

// EndOfMonth returns the end of the month (last day 23:59:59.999999999) for the given time.
//
// Parameters:
//   - t time.Time: The time to get the end of the month for.
//
// Returns:
//   - time.Time: The end of the month for the given time.
func EndOfMonth(t time.Time) time.Time {
	return EndOfDay(StartOfMonth(t).AddDate(0, 1, -1))
}

// DaysInMonth returns the number of days in the month for the given time.
//
// Parameters:
//   - t time.Time: The time to get the number of days in the month for.
//
// Returns:
//   - int: The number of days in the month.
func DaysInMonth(t time.Time) int {
	return EndOfMonth(t).Day()
}

// AddDays adds the specified number of days to the given time.
//
// Parameters:
//   - t time.Time: The time to add days to.
//   - days int: The number of days to add.
//
// Returns:
//   - time.Time: The resulting time after adding the days.
func AddDays(t time.Time, days int) time.Time {
	return t.AddDate(0, 0, days)
}

// SubtractDays subtracts the specified number of days from the given time.
//
// Parameters:
//   - t time.Time: The time to subtract days from.
//   - days int: The number of days to subtract.
//
// Returns:
//   - time.Time: The resulting time after subtracting the days.
func SubtractDays(t time.Time, days int) time.Time {
	return t.AddDate(0, 0, -days)
}

// DaysBetween calculates the number of days between two times.
//
// Parameters:
//   - t1 time.Time: The first time.
//   - t2 time.Time: The second time.
//
// Returns:
//   - int: The number of days between the two times.
func DaysBetween(t1, t2 time.Time) int {
	t1 = StartOfDay(t1)
	t2 = StartOfDay(t2)
	return int(t2.Sub(t1).Hours() / 24)
}

// IsSameDay checks if two times fall on the same day.
//
// Parameters:
//   - t1 time.Time: The first time.
//   - t2 time.Time: The second time.
//
// Returns:
//   - bool: True if the two times fall on the same day, false otherwise.
func IsSameDay(t1, t2 time.Time) bool {
	return t1.Year() == t2.Year() && t1.Month() == t2.Month() && t1.Day() == t2.Day()
}

// IsSameMonth checks if two times fall in the same month.
//
// Parameters:
//   - t1 time.Time: The first time.
//   - t2 time.Time: The second time.
//
// Returns:
//   - bool: True if the two times fall in the same month, false otherwise.
func IsSameMonth(t1, t2 time.Time) bool {
	return t1.Year() == t2.Year() && t1.Month() == t2.Month()
}

// IsSameYear checks if two times fall in the same year.
//
// Parameters:
//   - t1 time.Time: The first time.
//   - t2 time.Time: The second time.
//
// Returns:
//   - bool: True if the two times fall in the same year, false otherwise.
func IsSameYear(t1, t2 time.Time) bool {
	return t1.Year() == t2.Year()
}
