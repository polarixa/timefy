package timefy

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

// getWeekStartDay returns the week start day for this Timex instance.
// If Rule is nil, it returns the global default using GetWeekStartDay().
func (t *Timex) getWeekStartDay() time.Weekday {
	if t.Rule == nil {
		return GetWeekStartDay()
	}
	return t.weekStartDay
}

// getTimeFormats returns the time formats for this Timex instance.
// If Rule is nil, it returns the global TimeFormats.
func (t *Timex) getTimeFormats() []string {
	if t.Rule == nil {
		return TimeFormats
	}
	return t.timeFormats
}

// getTimeLocation returns the time location for this Timex instance.
// If Rule is nil or timeLocation is nil, it returns nil (which will use local).
func (t *Timex) getTimeLocation() *time.Location {
	if t.Rule == nil {
		return nil
	}
	return t.timeLocation
}

// With wraps the provided time value `v` into a Timex object, applying the default configuration.
//
// The function first checks if the global `DefaultConfig` is set. If it is not, it initializes a new
// `Config` with default values for `WeekStartDay` and `TimeFormats`. The configuration and time
// are then used to create a `Timex` object.
//
// Parameters:
//
//   - `v`: A time.Time value representing the time to be wrapped inside the Timex struct.
//
// Returns:
//
//   - A pointer to a `Timex` struct, which includes the provided time and either the default or
//     a custom configuration.
//
// Example:
//
//	t := time.Now()
//	timex := With(t) // This wraps the current time into a Timex object with the default configuration.
func With(v time.Time) *Timex {
	rule := NewRule()
	return &Timex{Time: v, Rule: rule}
}

// New creates a new Timex object for the provided time value `v`.
//
// The function calls the `With()` function, which wraps the given time in a `Timex` struct and applies
// the default configuration (or a custom one if set). Essentially, this function is a shorthand for
// creating a `Timex` object.
//
// Parameters:
//
//   - `v`: A time.Time value representing the time to be wrapped inside the Timex struct.
//
// Returns:
//   - A pointer to a `Timex` struct, which includes the provided time and the applicable configuration.
//
// Example:
//
//	t := time.Now()
//	timex := New(t) // This creates a new Timex object for the current time.
func New(v time.Time) *Timex {
	return With(v)
}

// With wraps the provided time value `v` into a Timex object using the current configuration `c`.
//
// This method is called on a `Config` object to create a new `Timex` instance, where the provided
// time and the configuration are bundled together. Unlike the global `With()` function, this method
// uses the specific configuration tied to the `Config` instance `c`.
//
// Parameters:
//
//   - `v`: A time.Time value representing the time to be wrapped inside the Timex struct.
//
// Returns:
//   - A pointer to a `Timex` struct, which includes the provided time and the specific `Config` instance.
//
// Example:
//
//	customConfig := &Config{WeekStartDay: time.Monday}
//	t := time.Now()
//	timex := customConfig.With(t) // This creates a Timex object with the custom configuration.
func (c *Rule) With(v time.Time) *Timex {
	return &Timex{Time: v, Rule: c}
}

// Parse attempts to parse a given set of string representations of time using the current `Config`.
//
// The function checks whether the `TimeLocation` field in the configuration `c` is set. If it is not set,
// the function defaults to parsing the strings based on the current time (using the system's local time zone).
// If `TimeLocation` is set, it uses the specified location to interpret the time. It delegates the parsing
// task to the `Parse()` method of a new `Timex` object created using `With()`, which wraps the current time
// with the existing configuration.
//
// Parameters:
//
//   - `s`: A variadic list of strings representing dates or times to be parsed.
//
// Returns:
//   - A `time.Time` value if one of the provided strings is successfully parsed.
//   - An error if parsing fails.
//
// Example:
//
//	config := &Config{TimeLocation: time.UTC}
//	parsedTime, err := config.Parse("2023-10-24T12:00:00Z") // Parses the string using UTC time location.
//
// If no `TimeLocation` is set, the function uses the local time zone:
//
//	config := &Config{}
//	parsedTime, err := config.Parse("2023-10-24T12:00:00") // Parses using the local time zone.
func (c *Rule) Parse(s ...string) (time.Time, error) {
	if c.timeLocation == nil {
		return c.With(time.Now()).Parse(s...)
	} else {
		return c.With(time.Now().In(c.timeLocation)).Parse(s...)
	}
}

// MustParse attempts to parse a given set of string representations of time using the current `Config`,
// and panics if parsing fails.
//
// This function checks if the `TimeLocation` field in the configuration `c` is set. If it is not set,
// it uses the system's local time zone for parsing. If `TimeLocation` is set, it utilizes the specified
// location to interpret the time. The actual parsing task is delegated to the `MustParse()` method
// of a new `Timex` object created using `With()`, which wraps the current time with the existing
// configuration.
//
// Parameters:
//
//   - `s`: A variadic list of strings representing dates or times to be parsed.
//
// Returns:
//   - A `time.Time` value representing the parsed time. This function will panic if parsing fails.
//
// Example:
//
//	config := &Config{TimeLocation: time.UTC}
//	parsedTime := config.MustParse("2023-10-24T12:00:00Z") // Parses the string using UTC time location.
//
//	 If no `TimeLocation` is set, the function uses the local time zone:
//
//	config := &Config{}
//	parsedTime := config.MustParse("2023-10-24T12:00:00") // Parses using the local time zone, panicking on failure.
func (c *Rule) MustParse(s ...string) time.Time {
	if c.timeLocation == nil {
		return c.With(time.Now()).MustParse(s...)
	} else {
		return c.With(time.Now().In(c.timeLocation)).MustParse(s...)
	}
}

// BeginningOfMinute returns a new time.Time value representing the start of the minute for the
// given Timex instance.
//
// The function uses the `Truncate()` method on the underlying time.Time of the Timex struct,
// truncating the time to the nearest minute. This effectively resets the seconds and
// nanoseconds to zero, providing the time corresponding to the beginning of the current minute.
//
// Returns:
//   - A `time.Time` value representing the start of the minute for the current Timex instance.
//
// Example:
//
//	t := Timex{Time: time.Now()}
//	startOfMinute := t.BeginningOfMinute() // This will return the current time truncated to the start of the minute.
func (t *Timex) BeginningOfMinute() time.Time {
	_, _, sec := t.Time.Clock()
	// Subtract the wall-clock seconds and nanoseconds; Round(0) strips the monotonic reading.
	return t.Time.Add(-(time.Duration(sec)*time.Second + time.Duration(t.Time.Nanosecond()))).Round(0)
}

// BeginningOfHour returns a new time.Time value representing the start of the hour for the
// given Timex instance.
//
// The function retrieves the year, month, and day from the underlying time.Time of the Timex struct
// using the `Date()` method. It then constructs a new time.Time value with the same year, month, and
// day but sets the hour to the current hour and resets the minutes, seconds, and nanoseconds to zero.
// This effectively provides the time corresponding to the beginning of the current hour.
//
// Returns:
//   - A `time.Time` value representing the start of the hour for the current Timex instance.
//
// Example:
//
//	t := Timex{Time: time.Now()}
//	startOfHour := t.BeginningOfHour() // This will return the current time truncated to the start of the hour.
func (t *Timex) BeginningOfHour() time.Time {
	_, min, sec := t.Time.Clock()
	// Work on the instant rather than rebuilding it with time.Date, which would always pick the
	// first occurrence of an ambiguous hour (DST fall-back). Round(0) strips the monotonic reading.
	off := time.Duration(min)*time.Minute + time.Duration(sec)*time.Second + time.Duration(t.Time.Nanosecond())
	return t.Time.Add(-off).Round(0)
}

// BeginningOfDay returns a new time.Time value representing the start of the day for the
// given Timex instance.
//
// The function retrieves the year, month, and day from the underlying time.Time of the Timex struct
// using the `Date()` method. It then constructs a new time.Time value with the same year, month, and
// day but sets the hour, minute, second, and nanosecond to zero. This effectively resets the time to
// midnight (00:00:00) of the same date, providing the start of the day.
//
// Returns:
//   - A `time.Time` value representing the start of the day for the current Timex instance.
//
// Example:
//
//	t := Timex{Time: time.Now()}
//	startOfDay := t.BeginningOfDay() // This will return the current date at 00:00:00 (midnight).
func (t *Timex) BeginningOfDay() time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, t.Time.Location())
}

// BeginningOfWeek returns a new time.Time value representing the start of the week for the
//
//	given Timex instance, based on the configured `WeekStartDay`.
//
// The function first calculates the beginning of the current day using `BeginningOfDay()`. Then, it determines
// the current weekday and calculates how many days to subtract in order to reach the start of the week, based on
// the configured start day (`WeekStartDay`). If `WeekStartDay` is set to a day other than Sunday, the function
// adjusts the number of days to move back accordingly.
//
// Returns:
//
//   - A `time.Time` value representing the start of the week for the current Timex instance.
//
// Example:
//
//	t := Timex{Time: time.Now(), Config: &Config{WeekStartDay: time.Monday}}
//	startOfWeek := t.BeginningOfWeek() // This will return the date and time at the start of the current week.
//
// Note:
//
//	If `WeekStartDay` is not set (defaults to Sunday), the function will return the preceding Sunday.
func (t *Timex) BeginningOfWeek() time.Time {
	day := t.BeginningOfDay()
	weekday := int(day.Weekday())
	startDay := t.getWeekStartDay()
	if startDay != time.Sunday {
		weekStartDayInt := int(startDay)
		if weekday < weekStartDayInt {
			weekday = weekday + 7 - weekStartDayInt
		} else {
			weekday = weekday - weekStartDayInt
		}
	}
	return day.AddDate(0, 0, -weekday)
}

// BeginningOfMonth returns a new time.Time value representing the start of the month for the
// given Timex instance.
//
// The function retrieves the year and month from the underlying time.Time of the Timex struct using the `Date()`
// method. It then constructs a new time.Time value with the same year and month, but sets the day to the first of
// the month and resets the hour, minute, second, and nanosecond to zero, providing the start of the month.
//
// Returns:
//   - A `time.Time` value representing the start of the month for the current Timex instance.
//
// Example:
//
//	t := Timex{Time: time.Now()}
//	startOfMonth := t.BeginningOfMonth() // This will return the date and time at the start of the current month.
func (t *Timex) BeginningOfMonth() time.Time {
	y, m, _ := t.Date()
	return time.Date(y, m, 1, 0, 0, 0, 0, t.Location())
}

// BeginningOfQuarter returns a new time.Time value representing the start of the quarter
// for the given Timex instance.
//
// The function first calculates the start of the current month using `BeginningOfMonth()`. It then determines
// how far the current month is into the quarter by calculating the offset (based on a 3-month quarter system).
// The offset is subtracted from the current month to return a time value corresponding to the start of the quarter
// (e.g., if the current month is April, it will return the start of April; if it's June, it will return the start of April).
//
// Returns:
//   - A `time.Time` value representing the start of the quarter for the current Timex instance.
//
// Example:
//
//	t := Timex{Time: time.Now()}
//	startOfQuarter := t.BeginningOfQuarter() // This will return the date and time at the start of the current quarter.
func (t *Timex) BeginningOfQuarter() time.Time {
	month := t.BeginningOfMonth()
	offset := (int(month.Month()) - 1) % 3
	return month.AddDate(0, -offset, 0)
}

// BeginningOfHalf returns a new time.Time value representing the start of the half-year
// for the given Timex instance.
//
// The function first calculates the start of the current month using `BeginningOfMonth()`. It then determines
// how far the current month is into the half-year by calculating the offset (based on a 6-month system).
// The offset is subtracted from the current month to return a time value corresponding to the start of the half-year
// (e.g., if the current month is July, it will return the start of July; if it's October, it will return the start of July).
//
// Returns:
//   - A `time.Time` value representing the start of the half-year for the current Timex instance.
//
// Example:
//
//	t := Timex{Time: time.Now()}
//	startOfHalf := t.BeginningOfHalf() // This will return the date and time at the start of the current half-year.
func (t *Timex) BeginningOfHalf() time.Time {
	month := t.BeginningOfMonth()
	offset := (int(month.Month()) - 1) % 6
	return month.AddDate(0, -offset, 0)
}

// BeginningOfYear returns a new time.Time value representing the start of the year
// for the given Timex instance.
//
// The function retrieves the year from the underlying time.Time of the Timex struct using the `Date()` method.
// It then constructs a new time.Time value with the same year but sets the month to January, and the day, hour,
// minute, second, and nanosecond to zero. This effectively resets the time to midnight (00:00:00) on January 1st
// of the current year.
//
// Returns:
//   - A `time.Time` value representing the start of the year for the current Timex instance.
//
// Example:
//
//	t := Timex{Time: time.Now()}
//	startOfYear := t.BeginningOfYear() // This will return January 1st of the current year at 00:00:00.
func (t *Timex) BeginningOfYear() time.Time {
	y, _, _ := t.Date()
	return time.Date(y, time.January, 1, 0, 0, 0, 0, t.Location())
}

// EndOfMinute returns a new time.Time value representing the end of the current minute
// for the given Timex instance.
//
// The function first calculates the start of the minute using `BeginningOfMinute()`, then adds one minute
// and subtracts a nanosecond. This adjustment effectively gives the time at the last nanosecond of the current minute.
//
// Returns:
//   - A `time.Time` value representing the end of the current minute for the Timex instance.
//
// Example:
//
//	t := Timex{Time: time.Now()}
//	endOfMinute := t.EndOfMinute() // This will return the time at the end of the current minute (59.999999999 nanoseconds).
func (t *Timex) EndOfMinute() time.Time {
	return t.BeginningOfMinute().Add(time.Minute - time.Nanosecond)
}

// EndOfHour returns a new time.Time value representing the end of the current hour
// for the given Timex instance.
//
// The function first calculates the start of the hour using `BeginningOfHour()`, then adds one hour
// and subtracts a nanosecond to reach the last nanosecond of the current hour.
//
// Returns:
//   - A `time.Time` value representing the end of the current hour for the Timex instance.
//
// Example:
//
//	t := Timex{Time: time.Now()}
//	endOfHour := t.EndOfHour() // This will return the time at the end of the current hour (59:59.999999999).
func (t *Timex) EndOfHour() time.Time {
	return t.BeginningOfHour().Add(time.Hour - time.Nanosecond)
}

// EndOfDay returns a new time.Time value representing the end of the current day
// for the given Timex instance.
//
// The function retrieves the year, month, and day from the underlying time.Time of the Timex struct
// using the `Date()` method. It then constructs a new time.Time value with the same year, month, and day
// but sets the time to 23:59:59.999999999, which is the last nanosecond of the day.
//
// Returns:
//   - A `time.Time` value representing the end of the current day for the Timex instance.
//
// Example:
//
//	t := Timex{Time: time.Now()}
//	endOfDay := t.EndOfDay() // This will return the date and time at the end of the current day (23:59:59.999999999).
func (t *Timex) EndOfDay() time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, 23, 59, 59, int(time.Second-time.Nanosecond), t.Location())
}

// EndOfWeek returns a new time.Time value representing the end of the current week
// for the given Timex instance, based on the configured `WeekStartDay`.
//
// The function first calculates the start of the current week using `BeginningOfWeek()`,
// then adds 7 days to reach the start of the next week, and subtracts one nanosecond
// to get the last nanosecond of the current week.
//
// Returns:
//   - A `time.Time` value representing the end of the current week for the Timex instance.
//
// Example:
//
//	t := Timex{Time: time.Now(), Config: &Config{WeekStartDay: time.Monday}}
//	endOfWeek := t.EndOfWeek() // This will return the date and time at the end of the current week.
//
// Note:
//
//	If `WeekStartDay` is not set (defaults to Sunday), the function will return the preceding Saturday at 23:59:59.999999999.
func (t *Timex) EndOfWeek() time.Time {
	return t.BeginningOfWeek().AddDate(0, 0, 7).Add(-time.Nanosecond)
}

// EndOfMonth returns a new time.Time value representing the end of the current month
// for the given Timex instance.
//
// The function first calculates the start of the month using `BeginningOfMonth()`,
// then adds one month to reach the beginning of the next month, and subtracts one nanosecond
// to obtain the last nanosecond of the current month.
//
// Returns:
//   - A `time.Time` value representing the end of the current month for the Timex instance.
//
// Example:
//
//	t := Timex{Time: time.Now()}
//	endOfMonth := t.EndOfMonth() // This will return the date and time at the end of the current month (e.g., 23:59:59.999999999 on the last day).
func (t *Timex) EndOfMonth() time.Time {
	return t.BeginningOfMonth().AddDate(0, 1, 0).Add(-time.Nanosecond)
}

// EndOfQuarter returns a new time.Time value representing the end of the current quarter
// for the given Timex instance.
//
// The function first calculates the start of the quarter using `BeginningOfQuarter()`,
// then adds three months to reach the beginning of the next quarter, and subtracts one nanosecond
// to obtain the last nanosecond of the current quarter.
//
// Returns:
//   - A `time.Time` value representing the end of the current quarter for the Timex instance.
//
// Example:
//
//	t := Timex{Time: time.Now()}
//	endOfQuarter := t.EndOfQuarter() // This will return the date and time at the end of the current quarter (e.g., 23:59:59.999999999 on the last day of the quarter).
func (t *Timex) EndOfQuarter() time.Time {
	return t.BeginningOfQuarter().AddDate(0, 3, 0).Add(-time.Nanosecond)
}

// EndOfHalf returns a new time.Time value representing the end of the current half-year
// for the given Timex instance.
//
// The function first calculates the start of the half-year using `BeginningOfHalf()`, then adds six months
// to reach the beginning of the next half-year, and subtracts one nanosecond to obtain the last nanosecond
// of the current half-year.
//
// Returns:
//   - A `time.Time` value representing the end of the current half-year for the Timex instance.
//
// Example:
//
//	t := Timex{Time: time.Now()}
//	endOfHalf := t.EndOfHalf() // This will return the date and time at the end of the current half-year (e.g., 23:59:59.999999999 on June 30th or December 31st).
func (t *Timex) EndOfHalf() time.Time {
	return t.BeginningOfHalf().AddDate(0, 6, 0).Add(-time.Nanosecond)
}

// EndOfYear returns a new time.Time value representing the end of the current year
// for the given Timex instance.
//
// The function first calculates the start of the year using `BeginningOfYear()`, then adds one year
// to reach the beginning of the next year, and subtracts one nanosecond to obtain the last nanosecond
// of the current year.
//
// Returns:
//   - A `time.Time` value representing the end of the current year for the Timex instance.
//
// Example:
//
//	t := Timex{Time: time.Now()}
//	endOfYear := t.EndOfYear() // This will return the date and time at the end of the current year (e.g., 23:59:59.999999999 on December 31st).
func (t *Timex) EndOfYear() time.Time {
	return t.BeginningOfYear().AddDate(1, 0, 0).Add(-time.Nanosecond)
}

// Monday returns a new time.Time value representing Monday of the Monday-to-Sunday week that
// contains the date given by the provided date string(s), or the current date if none are provided.
//
// If a date string is provided, it parses the date using the Timex `Parse` method. If no string is provided,
// the function defaults to the start of the current day. It then determines the weekday of the parsed date.
// If the weekday is Sunday (represented by 0), it adjusts to be the seventh day for calculation purposes.
// The function finally moves to Monday of that Monday-to-Sunday week (which may be after the reference day).
//
// Parameters:
//   - s ...string: Optional date string(s) to parse; if none are provided, it defaults to today.
//
// Returns:
//   - A `time.Time` value representing the date and time at the Monday of that week, keeping the time of day.
//
// Example:
//
//	t := Timex{Time: time.Now()}
//	lastMonday := t.Monday() // Returns the start of the most recent Monday.
//
// Note:
//
//	If parsing fails, the function will panic with an error.
func (t *Timex) Monday(s ...string) time.Time {
	var parseTime time.Time
	var err error
	if len(s) > 0 {
		parseTime, err = t.Parse(s...)
		if err != nil {
			panic(err)
		}
	} else {
		parseTime = t.BeginningOfDay()
	}
	weekday := int(parseTime.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	return parseTime.AddDate(0, 0, -weekday+1)
}

// Tuesday returns a new time.Time value representing Tuesday of the Monday-to-Sunday week that
// contains the date given by the provided date string(s), or the current date if none are provided.
//
// If a date string is provided, it parses the date using the Timex `Parse` method. If no string is provided,
// the function defaults to the start of the current day. It then determines the weekday of the parsed date.
// If the weekday is Sunday (represented by 0), it adjusts to be the seventh day for calculation purposes.
// The function finally moves to Tuesday of that Monday-to-Sunday week (which may be after the reference day).
//
// Parameters:
//   - s ...string: Optional date string(s) to parse; if none are provided, it defaults to today.
//
// Returns:
//   - A `time.Time` value representing the date and time at the Tuesday of that week, keeping the time of day.
//
// Example:
//
//	t := Timex{Time: time.Now()}
//	lastTuesday := t.Tuesday() // Returns the start of the most recent Tuesday.
//
// Note:
//
//	If parsing fails, the function will panic with an error.
func (t *Timex) Tuesday(s ...string) time.Time {
	var parseTime time.Time
	var err error
	if len(s) > 0 {
		parseTime, err = t.Parse(s...)
		if err != nil {
			panic(err)
		}
	} else {
		parseTime = t.BeginningOfDay()
	}
	weekday := int(parseTime.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	return parseTime.AddDate(0, 0, -weekday+2)
}

// Wednesday returns a new time.Time value representing Wednesday of the Monday-to-Sunday week that
// contains the date given by the provided date string(s), or the current date if none are provided.
//
// If a date string is provided, it parses the date using the Timex `Parse` method. If no string is provided,
// the function defaults to the start of the current day. It then determines the weekday of the parsed date.
// If the weekday is Sunday (represented by 0), it adjusts to be the seventh day for calculation purposes.
// The function finally moves to Wednesday of that Monday-to-Sunday week (which may be after the reference day).
//
// Parameters:
//   - s ...string: Optional date string(s) to parse; if none are provided, it defaults to today.
//
// Returns:
//   - A `time.Time` value representing the date and time at the Wednesday of that week, keeping the time of day.
//
// Example:
//
//	t := Timex{Time: time.Now()}
//	lastWednesday := t.Wednesday() // Returns the start of the most recent Wednesday.
//
// Note:
//
//	If parsing fails, the function will panic with an error.
func (t *Timex) Wednesday(s ...string) time.Time {
	var parseTime time.Time
	var err error
	if len(s) > 0 {
		parseTime, err = t.Parse(s...)
		if err != nil {
			panic(err)
		}
	} else {
		parseTime = t.BeginningOfDay()
	}
	weekday := int(parseTime.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	return parseTime.AddDate(0, 0, -weekday+3)
}

// Thursday returns a new time.Time value representing Thursday of the Monday-to-Sunday week that
// contains the date given by the provided date string(s), or the current date if none are provided.
//
// If a date string is provided, it parses the date using the Timex `Parse` method. If no string is provided,
// the function defaults to the start of the current day. It then determines the weekday of the parsed date.
// If the weekday is Sunday (represented by 0), it adjusts to be the seventh day for calculation purposes.
// The function finally moves to Thursday of that Monday-to-Sunday week (which may be after the reference day).
//
// Parameters:
//   - s ...string: Optional date string(s) to parse; if none are provided, it defaults to today.
//
// Returns:
//   - A `time.Time` value representing the date and time at the Thursday of that week, keeping the time of day.
//
// Example:
//
//	t := Timex{Time: time.Now()}
//	lastThursday := t.Thursday() // Returns the start of the most recent Thursday.
//
// Note:
//
//	If parsing fails, the function will panic with an error.
func (t *Timex) Thursday(s ...string) time.Time {
	var parseTime time.Time
	var err error
	if len(s) > 0 {
		parseTime, err = t.Parse(s...)
		if err != nil {
			panic(err)
		}
	} else {
		parseTime = t.BeginningOfDay()
	}
	weekday := int(parseTime.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	return parseTime.AddDate(0, 0, -weekday+4)
}

// Friday returns a new time.Time value representing Friday of the Monday-to-Sunday week that
// contains the date given by the provided date string(s), or the current date if none are provided.
//
// If a date string is provided, it parses the date using the Timex `Parse` method. If no string is provided,
// the function defaults to the start of the current day. It then determines the weekday of the parsed date.
// If the weekday is Sunday (represented by 0), it adjusts to be the seventh day for calculation purposes.
// The function finally moves to Friday of that Monday-to-Sunday week (which may be after the reference day).
//
// Parameters:
//   - s ...string: Optional date string(s) to parse; if none are provided, it defaults to today.
//
// Returns:
//   - A `time.Time` value representing the date and time at the Friday of that week, keeping the time of day.
//
// Example:
//
//	t := Timex{Time: time.Now()}
//	lastFriday := t.Friday() // Returns the start of the most recent Friday.
//
// Note:
//
//	If parsing fails, the function will panic with an error.
func (t *Timex) Friday(s ...string) time.Time {
	var parseTime time.Time
	var err error
	if len(s) > 0 {
		parseTime, err = t.Parse(s...)
		if err != nil {
			panic(err)
		}
	} else {
		parseTime = t.BeginningOfDay()
	}
	weekday := int(parseTime.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	return parseTime.AddDate(0, 0, -weekday+5)
}

// Saturday returns a new time.Time value representing Saturday of the Monday-to-Sunday week that
// contains the date given by the provided date string(s), or the current date if none are provided.
//
// If a date string is provided, it parses the date using the Timex `Parse` method. If no string is provided,
// the function defaults to the start of the current day. It then determines the weekday of the parsed date.
// If the weekday is Sunday (represented by 0), it adjusts to be the seventh day for calculation purposes.
// The function finally moves to Saturday of that Monday-to-Sunday week (which may be after the reference day).
//
// Parameters:
//   - s ...string: Optional date string(s) to parse; if none are provided, it defaults to today.
//
// Returns:
//   - A `time.Time` value representing the date and time at the Saturday of that week, keeping the time of day.
//
// Example:
//
//	t := Timex{Time: time.Now()}
//	lastSaturday := t.Saturday() // Returns the start of the most recent Saturday.
//
// Note:
//
//	If parsing fails, the function will panic with an error.
func (t *Timex) Saturday(s ...string) time.Time {
	var parseTime time.Time
	var err error
	if len(s) > 0 {
		parseTime, err = t.Parse(s...)
		if err != nil {
			panic(err)
		}
	} else {
		parseTime = t.BeginningOfDay()
	}
	weekday := int(parseTime.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	return parseTime.AddDate(0, 0, -weekday+6)
}

// Sunday returns a new time.Time value representing the most recent or upcoming Sunday
// based on the provided date string(s) or the current date if no date strings are provided.
//
// If a date string is provided, it parses the date using the Timex `Parse` method. If no string is provided,
// the function defaults to the start of the current day. It then determines the weekday of the parsed date.
// If the weekday is Sunday (represented by 0), it adjusts to be the seventh day for calculation purposes.
// The function finally adds the necessary days to reach the next Sunday.
//
// Parameters:
//   - s ...string: Optional date string(s) to parse; if none are provided, it defaults to today.
//
// Returns:
//   - A `time.Time` value representing the date and time at the start of the most recent or upcoming Sunday.
//
// Example:
//
//	t := Timex{Time: time.Now()}
//	lastSunday := t.Sunday() // Returns the date and time at the start of the most recent or upcoming Sunday.
//
// Note:
//
//	If parsing fails, the function will panic with an error.
func (t *Timex) Sunday(s ...string) time.Time {
	var parseTime time.Time
	var err error
	if len(s) > 0 {
		parseTime, err = t.Parse(s...)
		if err != nil {
			panic(err)
		}
	} else {
		parseTime = t.BeginningOfDay()
	}
	weekday := int(parseTime.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	return parseTime.AddDate(0, 0, (7 - weekday))
}

// EndOfSunday returns a new time.Time value representing the end of the most recent or upcoming Sunday
// for the given Timex instance.
//
// The function first calculates the start of the most recent or upcoming Sunday using `Sunday()`,
// then constructs a new Timex instance with that date and time. Finally, it calls `EndOfDay()` on this
// instance to retrieve the time at the last nanosecond of the Sunday.
//
// Returns:
//   - A `time.Time` value representing the end of the most recent or upcoming Sunday for the Timex instance.
//
// Example:
//
//	t := Timex{Time: time.Now()}
//	endOfSunday := t.EndOfSunday() // Returns the date and time at the end of the most recent or upcoming Sunday (23:59:59.999999999).
func (t *Timex) EndOfSunday() time.Time {
	return New(t.Sunday()).EndOfDay()
}

// Quarter returns the current quarter of the year for the given Timex instance,
// where the year is divided into four quarters: Q1 (January-March), Q2 (April-June),
// Q3 (July-September), and Q4 (October-December).
//
// The function calculates the quarter by taking the current month (retrieved using
// the `Month()` method), subtracting one (to account for zero-based indexing),
// dividing by 3, and then adding 1 to get the quarter number. This results in
// a value between 1 and 4, representing the quarter of the year.
//
// Returns:
//   - A `uint` value representing the current quarter (1 to 4) for the Timex instance.
//
// Example:
//
//	t := Timex{Time: time.Now()}
//	quarter := t.Quarter() // Returns the current quarter number (1, 2, 3, or 4).
func (t *Timex) Quarter() uint {
	return (uint(t.Month())-1)/3 + 1
}

// weekdaySafe is a helper function that calculates the time for a specific weekday.
// It handles parsing and weekday normalization logic used by all *Safe weekday methods.
// The offset parameter specifies the target weekday where Monday=1, Tuesday=2, ..., Sunday=7.
func (t *Timex) weekdaySafe(dayName string, offset int, s ...string) (time.Time, error) {
	var parseTime time.Time
	var err error
	if len(s) > 0 {
		parseTime, err = t.Parse(s...)
		if err != nil {
			return time.Time{}, fmt.Errorf("timefy: failed to parse date for %s: %w", dayName, err)
		}
	} else {
		parseTime = t.BeginningOfDay()
	}
	weekday := int(parseTime.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	return parseTime.AddDate(0, 0, -weekday+offset), nil
}

// MondaySafe returns a new time.Time value representing Monday of the Monday-to-Sunday week that
// contains the date given by the provided date string(s), or the current date if none are provided.
// Unlike Monday(), this method returns an error instead of panicking on parse failures.
//
// Parameters:
//   - s ...string: Optional date string(s) to parse; if none are provided, it defaults to today.
//
// Returns:
//   - A `time.Time` value representing the date and time at the Monday of that week, keeping the time of day.
//   - An error if parsing fails.
func (t *Timex) MondaySafe(s ...string) (time.Time, error) {
	return t.weekdaySafe("Monday", 1, s...)
}

// TuesdaySafe returns a new time.Time value representing Tuesday of the Monday-to-Sunday week that
// contains the date given by the provided date string(s), or the current date if none are provided.
// Unlike Tuesday(), this method returns an error instead of panicking on parse failures.
//
// Parameters:
//   - s ...string: Optional date string(s) to parse; if none are provided, it defaults to today.
//
// Returns:
//   - A `time.Time` value representing the date and time at the Tuesday of that week, keeping the time of day.
//   - An error if parsing fails.
func (t *Timex) TuesdaySafe(s ...string) (time.Time, error) {
	return t.weekdaySafe("Tuesday", 2, s...)
}

// WednesdaySafe returns a new time.Time value representing Wednesday of the Monday-to-Sunday week that
// contains the date given by the provided date string(s), or the current date if none are provided.
// Unlike Wednesday(), this method returns an error instead of panicking on parse failures.
//
// Parameters:
//   - s ...string: Optional date string(s) to parse; if none are provided, it defaults to today.
//
// Returns:
//   - A `time.Time` value representing the date and time at the Wednesday of that week, keeping the time of day.
//   - An error if parsing fails.
func (t *Timex) WednesdaySafe(s ...string) (time.Time, error) {
	return t.weekdaySafe("Wednesday", 3, s...)
}

// ThursdaySafe returns a new time.Time value representing Thursday of the Monday-to-Sunday week that
// contains the date given by the provided date string(s), or the current date if none are provided.
// Unlike Thursday(), this method returns an error instead of panicking on parse failures.
//
// Parameters:
//   - s ...string: Optional date string(s) to parse; if none are provided, it defaults to today.
//
// Returns:
//   - A `time.Time` value representing the date and time at the Thursday of that week, keeping the time of day.
//   - An error if parsing fails.
func (t *Timex) ThursdaySafe(s ...string) (time.Time, error) {
	return t.weekdaySafe("Thursday", 4, s...)
}

// FridaySafe returns a new time.Time value representing Friday of the Monday-to-Sunday week that
// contains the date given by the provided date string(s), or the current date if none are provided.
// Unlike Friday(), this method returns an error instead of panicking on parse failures.
//
// Parameters:
//   - s ...string: Optional date string(s) to parse; if none are provided, it defaults to today.
//
// Returns:
//   - A `time.Time` value representing the date and time at the Friday of that week, keeping the time of day.
//   - An error if parsing fails.
func (t *Timex) FridaySafe(s ...string) (time.Time, error) {
	return t.weekdaySafe("Friday", 5, s...)
}

// SaturdaySafe returns a new time.Time value representing Saturday of the Monday-to-Sunday week that
// contains the date given by the provided date string(s), or the current date if none are provided.
// Unlike Saturday(), this method returns an error instead of panicking on parse failures.
//
// Parameters:
//   - s ...string: Optional date string(s) to parse; if none are provided, it defaults to today.
//
// Returns:
//   - A `time.Time` value representing the date and time at the Saturday of that week, keeping the time of day.
//   - An error if parsing fails.
func (t *Timex) SaturdaySafe(s ...string) (time.Time, error) {
	return t.weekdaySafe("Saturday", 6, s...)
}

// SundaySafe returns a new time.Time value representing the upcoming Sunday
// based on the provided date string(s) or the current date if no date strings are provided.
// Unlike Sunday(), this method returns an error instead of panicking on parse failures.
//
// Note: Like Sunday(), the result is the Sunday that ends the current Monday-to-Sunday week.
// If the reference day is already a Sunday, the same day is returned (not the following Sunday).
//
// Parameters:
//   - s ...string: Optional date string(s) to parse; if none are provided, it defaults to today.
//
// Returns:
//   - A `time.Time` value representing the date and time at the start of the upcoming Sunday.
//   - An error if parsing fails.
func (t *Timex) SundaySafe(s ...string) (time.Time, error) {
	var parseTime time.Time
	var err error
	if len(s) > 0 {
		parseTime, err = t.Parse(s...)
		if err != nil {
			return time.Time{}, fmt.Errorf("timefy: failed to parse date for Sunday: %w", err)
		}
	} else {
		parseTime = t.BeginningOfDay()
	}
	weekday := int(parseTime.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	return parseTime.AddDate(0, 0, (7 - weekday)), nil
}

// Parse interprets the provided date string(s) and converts them into a time.Time value.
// It attempts to parse each string according to the configured formats, adjusting for the current time
// and location as necessary.
//
// The function checks each string for time information and adjusts the date components accordingly.
// If a date component is missing (e.g., day, month, year) in the parsed string, it will replace
// it with the corresponding component from the current time. If the input string includes only time,
// it sets the day and month to those of the current date. The function handles various cases for
// the input strings and parses them into a valid time.Time value.
//
// Parameters:
//   - `s ...string`: One or more date strings to be parsed. The function will try to parse each string
//     in the order provided and will return the first successful parsed time.
//
// Returns:
//   - `value`: A `time.Time` value representing the parsed date and time.
//   - `err`: An error value indicating any issues that occurred during parsing; if parsing is successful,
//     this will be nil.
//
// Example:
//
//	t := Timex{Time: time.Now()}
//	parsedTime, err := t.Parse("2023-10-25 15:04", "15:04") // Attempts to parse the given date strings.
//	if err != nil {
//		// Handle error
//	}
//
// Note:
//   - The function modifies the parsed date based on the current time when certain components are missing.
//   - Strings that fail to parse are skipped. The result is built from the strings that parsed, in order;
//     if none parse (or no string is given) it returns the zero time.Time and the last error.
func (t *Timex) Parse(s ...string) (value time.Time, err error) {
	// A location that failed to load (see Rule.WithLocationRFC) must not silently fall back to local time.
	if t.Rule != nil && t.locErr != nil {
		return time.Time{}, t.locErr
	}
	if len(s) == 0 {
		return time.Time{}, errors.New("timefy: no date string provided")
	}

	// The rule's location, when set, decides how strings are interpreted, whichever way the Timex was built.
	base := t.Time
	if loc := t.getTimeLocation(); loc != nil {
		base = base.In(loc)
	}

	var (
		setCurrentTime  bool
		parseTime       []int
		currentLocation = base.Location()
		onlyTimeInStr   = true
		currentTime     []int
		result          time.Time
		parsed          bool
		lastErr         error
	)

	currentTime = FormatTimex(base)

	for _, str := range s {
		parsedValue, layout, perr := t.parseWithFormat(str, currentLocation)
		if perr != nil {
			// A failing string never discards a value parsed from another string.
			lastErr = perr
			continue
		}

		// Decide whether the string carries a time from the layout that matched it,
		// not from a regular expression over the input text.
		hasTimeInStr := layoutHasTime(layout)
		onlyTimeInStr = hasTimeInStr && onlyTimeInStr && TimeOnlyRegexp.MatchString(str)

		location := parsedValue.Location()
		parseTime = FormatTimex(parsedValue)

		for i, v := range parseTime {
			// Don't reset hour, minute, second if current time str including time
			if hasTimeInStr && i <= 3 {
				continue
			}

			// If value is zero, replace it with current time
			if v == 0 {
				if setCurrentTime {
					parseTime[i] = currentTime[i]
				}
			} else {
				setCurrentTime = true
			}

			// if current time only includes time, should change day, month to current time
			if onlyTimeInStr {
				if i == 4 || i == 5 {
					parseTime[i] = currentTime[i]
					continue
				}
			}
		}
		result = time.Date(parseTime[6], time.Month(parseTime[5]), parseTime[4], parseTime[3], parseTime[2], parseTime[1], parseTime[0], location)
		currentTime = FormatTimex(result)
		parsed = true
	}

	if !parsed {
		return time.Time{}, lastErr
	}
	return result, nil
}

// layoutHasTime reports whether a time layout contains an hour element (15, 03 or 3),
// meaning the strings it parses carry a time of day.
func layoutHasTime(layout string) bool {
	return strings.Contains(layout, "15") || strings.Contains(layout, "3")
}

// MustParse interprets the provided date string(s) and converts them into a time.Time value,
// but it panics if parsing fails. This is useful for scenarios where a valid time is essential,
// and failure to parse should result in a program error.
//
// The function calls the `Parse` method, which attempts to interpret each date string according
// to the configured formats, adjusting for the current time and location as necessary. If
// the parsing succeeds, the resulting time.Time value is returned; if it fails, a panic is triggered
// with the associated error message.
//
// Parameters:
//   - `s ...string`: One or more date strings to be parsed. The function will try to parse each string
//     in the order provided and will panic if none can be successfully parsed.
//
// Returns:
//   - `v`: A `time.Time` value representing the parsed date and time, guaranteed to be valid
//     as the function will panic on failure.
//
// Example:
//
//	t := Timex{Time: time.Now()}
//	parsedTime := t.MustParse("2023-10-25 15:04") // Attempts to parse the given date strings.
//	// If parsing fails, it will panic.
//
// Note:
//   - This function is intended for use cases where the caller must have a valid time returned,
//     and any parsing errors should be treated as unrecoverable conditions.
func (t *Timex) MustParse(s ...string) (v time.Time) {
	v, err := t.Parse(s...)
	if err != nil {
		panic(err)
	}
	return v
}

// Between checks if the current Timex instance's time falls within the specified time range,
// defined by the provided `begin` and `end` date strings. The function parses the strings into
// time.Time values and determines if the current time is after the `begin` time and before the `end` time.
//
// Parameters:
//   - `begin`: A string representing the start of the time range; it will be parsed into a time.Time value.
//   - `end`: A string representing the end of the time range; it will also be parsed into a time.Time value.
//
// Returns:
//   - A boolean value indicating whether the current time is within the specified range (exclusive).
//     It returns true if the current time is after `begin` and before `end`; otherwise, it returns false.
//
// Example:
//
//	t := Timex{Time: time.Now()}
//	isInRange := t.Between("2023-10-01", "2023-10-31") // Checks if the current time is in October 2023.
//
// Note:
//   - If the parsing of `begin` or `end` fails, the function will panic due to the usage of `MustParse`.
func (t *Timex) Between(begin, end string) bool {
	beginTime := t.MustParse(begin)
	endTime := t.MustParse(end)
	return t.After(beginTime) && t.Before(endTime)
}

// FormatRFC formats the time using the provided layout.
//
// Parameters:
//   - `layout`: A TimeFormatRFC value representing the layout to use for formatting.
//
// Returns:
//   - A string representing the formatted time.
//
// Example:
//
//	t := Timex{Time: time.Now()}
//	formatted := t.FormatRFC(TimeFormat20060102T150405) // Returns the formatted time in the format "2023-08-15 13:45:30".
func (t *Timex) FormatRFC(layout TimeFormatRFC) string {
	return t.Time.Format(string(layout))
}

// FormatRFCshort formats the time using the provided layout.
//
// Parameters:
//   - `layout`: A TimeRFC value representing the layout to use for formatting.
//
// Returns:
//   - A string representing the formatted time.
//
// Example:
//
//	t := Timex{Time: time.Now()}
//	formatted := t.FormatRFCshort(TimeRFC01T150405) // Returns the formatted time in the format "13:45:30".
func (t *Timex) FormatRFCshort(layout TimeRFC) string {
	return t.Time.Format(string(layout))
}

// DefaultFormatRFC formats the time using a default layout that includes the date and time with nanosecond precision.
//
// Returns:
//   - A string representing the formatted time in the default RFC format.
func (t *Timex) DefaultFormatRFC() string {
	return t.FormatRFC(TimeFormat20060102150405999999)
}

// DurationInDays returns the number of whole calendar days from the current time to `other`
// (`other` minus the current time), truncated toward zero. It is negative when `other` is earlier.
// Days are counted on the calendar of the current time's location, so a 23- or 25-hour day around
// a DST change still counts as one day.
//
// Example:
//
//	t := Timex{Time: time.Now()}
//	days := t.DurationInDays(time.Now().AddDate(0, 0, 1)) // Returns 1
func (t *Timex) DurationInDays(other time.Time) int {
	return wholeUnits(t.Time, other, func(from time.Time, n int) time.Time { return from.AddDate(0, 0, n) },
		func(d time.Duration) int { return int(d.Hours() / 24) })
}

// DurationInWeeks returns the number of whole weeks from the current time to `other`
// (`other` minus the current time), truncated toward zero. It is negative when `other` is earlier.
//
// Example:
//
//	t := Timex{Time: time.Now()}
//	weeks := t.DurationInWeeks(time.Now().AddDate(0, 0, 7)) // Returns 1
func (t *Timex) DurationInWeeks(other time.Time) int {
	return t.DurationInDays(other) / 7
}

// DurationInMonths returns the number of whole calendar months from the current time to `other`
// (`other` minus the current time), truncated toward zero. A month is complete only when the same
// day and time of day has been reached, so January 31 to February 1 is 0 months.
// It is negative when `other` is earlier.
//
// Example:
//
//	t := Timex{Time: time.Now()}
//	months := t.DurationInMonths(time.Now().AddDate(0, 1, 0)) // Returns 1
func (t *Timex) DurationInMonths(other time.Time) int {
	return wholeUnits(t.Time, other, func(from time.Time, n int) time.Time { return from.AddDate(0, n, 0) },
		func(d time.Duration) int { return int(d.Hours() / (24 * 28)) })
}

// DurationInYears returns the number of whole calendar years from the current time to `other`
// (`other` minus the current time), truncated toward zero. It is negative when `other` is earlier.
//
// Example:
//
//	t := Timex{Time: time.Now()}
//	years := t.DurationInYears(time.Now().AddDate(1, 0, 0)) // Returns 1
func (t *Timex) DurationInYears(other time.Time) int {
	return wholeUnits(t.Time, other, func(from time.Time, n int) time.Time { return from.AddDate(n, 0, 0) },
		func(d time.Duration) int { return int(d.Hours() / (24 * 365)) })
}

// wholeUnits returns the largest n (toward zero) such that add(from, n) has not passed `to`.
// `add` advances a time by n calendar units and `estimate` gives a starting guess from a duration;
// the result is then corrected exactly, so DST changes and month lengths cannot skew it.
func wholeUnits(from, to time.Time, add func(time.Time, int) time.Time, estimate func(time.Duration) int) int {
	// Work in the location of `from` so calendar arithmetic follows its wall clock.
	to = to.In(from.Location())
	if to.Before(from) {
		return -wholeUnits(to, from, add, estimate)
	}
	n := estimate(to.Sub(from))
	for n > 0 && add(from, n).After(to) {
		n--
	}
	for !add(from, n+1).After(to) {
		n++
	}
	return n
}

// IsWeekend checks if the current time is a weekend (Saturday or Sunday).
//
// Returns:
//   - A boolean value indicating whether the current time is a weekend.
//
// Example:
//
//	t := Timex{Time: time.Now()}
//	isWeekend := t.IsWeekend() // Returns true if the current time is a weekend
func (t *Timex) IsWeekend() bool {
	weekday := t.Weekday()
	return weekday == time.Saturday || weekday == time.Sunday
}

// IsWeekday checks if the current time is a weekday (Monday to Friday).
//
// Returns:
//   - A boolean value indicating whether the current time is a weekday.
//
// Example:
//
//	t := Timex{Time: time.Now()}
//	isWeekday := t.IsWeekday() // Returns true if the current time is a weekday
func (t *Timex) IsWeekday() bool {
	return !t.IsWeekend()
}

// IsBefore checks if the current time is before the specified time.
//
// Parameters:
//   - `other`: The other time to compare against.
//
// Returns:
//   - A boolean value indicating whether the current time is before the specified time.
//
// Example:
//
//	t := Timex{Time: time.Now()}
//	isBefore := t.IsBefore(time.Now().AddDate(0, 0, 1)) // Returns true if the current time is before the specified time
func (t *Timex) IsBefore(other time.Time) bool {
	return t.Time.Before(other)
}

// IsAfter checks if the current time is after the specified time.
//
// Parameters:
//   - `other`: The other time to compare against.
//
// Returns:
//   - A boolean value indicating whether the current time is after the specified time.
//
// Example:
//
//	t := Timex{Time: time.Now()}
//	isAfter := t.IsAfter(time.Now().AddDate(0, 0, 1)) // Returns true if the current time is after the specified time
func (t *Timex) IsAfter(other time.Time) bool {
	return t.Time.After(other)
}

// IsBetween checks if the current time is between the specified start and end times (inclusive).
//
// Parameters:
//   - `start`: The start time of the range.
//   - `end`: The end time of the range.
//
// Returns:
//   - A boolean value indicating whether the current time is between the specified start and end times.
//
// Example:
//
//	t := Timex{Time: time.Now()}
//	isBetween := t.IsBetween(time.Now().AddDate(0, 0, 1), time.Now().AddDate(0, 0, 2)) // Returns true if the current time is between the specified start and end times
func (t *Timex) IsBetween(start, end time.Time) bool {
	return (t.Equal(start) || t.After(start)) && (t.Equal(end) || t.Before(end))
}

// IsSameDay checks if the current time is on the same day as the specified time.
//
// Parameters:
//   - `other`: The other time to compare against.
//
// Returns:
//   - A boolean value indicating whether the current time is on the same day as the specified time.
//
// Example:
//
//	t := Timex{Time: time.Now()}
//	isSameDay := t.IsSameDay(time.Now().AddDate(0, 0, 1)) // Returns true if the current time is on the same day as the specified time
func (t *Timex) IsSameDay(other time.Time) bool {
	return t.Year() == other.Year() && t.YearDay() == other.YearDay()
}

// IsSameMonth checks if the current time is on the same month as the specified time.
//
// Parameters:
//   - `other`: The other time to compare against.
//
// Returns:
//   - A boolean value indicating whether the current time is on the same month as the specified time.
//
// Example:
//
//	t := Timex{Time: time.Now()}
//	isSameMonth := t.IsSameMonth(time.Now().AddDate(0, 1, 0)) // Returns true if the current time is on the same month as the specified time
func (t *Timex) IsSameMonth(other time.Time) bool {
	return t.Year() == other.Year() && t.Month() == other.Month()
}

// IsSameYear checks if the current time is on the same year as the specified time.
//
// Parameters:
//   - `other`: The other time to compare against.
//
// Returns:
//   - A boolean value indicating whether the current time is on the same year as the specified time.
//
// Example:
//
//	t := Timex{Time: time.Now()}
//	isSameYear := t.IsSameYear(time.Now().AddDate(1, 0, 0)) // Returns true if the current time is on the same year as the specified time
func (t *Timex) IsSameYear(other time.Time) bool {
	return t.Year() == other.Year()
}

// IsToday checks if the current time is today.
//
// Returns:
//   - A boolean value indicating whether the current time is today.
//
// Example:
//
//	t := Timex{Time: time.Now()}
//	isToday := t.IsToday() // Returns true if the current time is today
func (t *Timex) IsToday() bool {
	now := time.Now()
	return t.IsSameDay(now)
}

// IsYesterday checks if the current time is yesterday.
//
// Returns:
//   - A boolean value indicating whether the current time is yesterday.
//
// Example:
//
//	t := Timex{Time: time.Now()}
//	isYesterday := t.IsYesterday() // Returns true if the current time is yesterday
func (t *Timex) IsYesterday() bool {
	yesterday := time.Now().AddDate(0, 0, -1)
	return t.IsSameDay(yesterday)
}

// IsTomorrow checks if the current time is tomorrow.
//
// Returns:
//   - A boolean value indicating whether the current time is tomorrow.
//
// Example:
//
//	t := Timex{Time: time.Now()}
//	isTomorrow := t.IsTomorrow() // Returns true if the current time is tomorrow
func (t *Timex) IsTomorrow() bool {
	tomorrow := time.Now().AddDate(0, 0, 1)
	return t.IsSameDay(tomorrow)
}

// IsPast checks if the current time is in the past.
//
// Returns:
//   - A boolean value indicating whether the current time is in the past.
//
// Example:
//
//	t := Timex{Time: time.Now()}
//	isPast := t.IsPast() // Returns true if the current time is in the past
func (t *Timex) IsPast() bool {
	return t.Before(time.Now())
}

// IsFuture checks if the current time is in the future.
//
// Returns:
//   - A boolean value indicating whether the current time is in the future.
//
// Example:
//
//	t := Timex{Time: time.Now()}
//	isFuture := t.IsFuture() // Returns true if the current time is in the future
func (t *Timex) IsFuture() bool {
	return t.After(time.Now())
}

// parseWithFormat attempts to parse a given date/time string `s` using a series of predefined formats
// specified in the `TimeFormats` slice of the Timex instance. It tries to parse the string in the provided
// `location` and returns the parsed time value or an error if parsing fails.
//
// The function iterates over each format in `TimeFormats`, attempting to parse the string with each format.
// If the parsing is successful (i.e., no error is returned), it returns the parsed time value. If all formats
// fail to parse the string, it returns an error indicating the failure.
//
// Parameters:
//   - `s`: A string representing the date/time to be parsed.
//   - `location`: A pointer to a time.Location value, which specifies the timezone to use during parsing.
//
// Returns:
//   - `v`: A time.Time value representing the parsed date/time if successful.
//   - `err`: An error value indicating any issues that occurred during parsing; if parsing is successful,
//     this will be nil.
//
// Example:
//
//	t := Timex{TimeFormats: []string{"2006-01-02 15:04:05", "2006-01-02"}}
//	parsedTime, err := t.parseWithFormat("2023-10-25 12:30:00", time.UTC) // Attempts to parse the given string.
//	if err != nil {
//		// Handle parsing error
//	}
//
// Note:
//   - The function will return the first successfully parsed time value and ignore any subsequent formats.
func (t *Timex) parseWithFormat(s string, location *time.Location) (v time.Time, layout string, err error) {
	formats := t.getTimeFormats()
	if len(formats) == 0 {
		err = fmt.Errorf("timefy: no time formats configured for parsing %q; ensure Timex is created via With() or NewRule().With()", s)
		return
	}
	for _, format := range formats {
		v, err = time.ParseInLocation(format, s, location)

		if err == nil {
			layout = format
			return
		}
	}
	v = time.Time{}
	err = fmt.Errorf("timefy: can't parse string as time: %q", s)
	return
}
