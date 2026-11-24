// Package main demonstrates the package-level helpers that work on a given
// time.Time or on "now": BeginOfDay/EndOfDayOf, PrevBeginOfDay/PrevEndOfDay,
// AddSecond..AddDay, SinceHour/Minute/Second, week start configuration,
// current-period boundaries, Parse/MustParse with locations, and RFC formatting.
//
// Run: go run ./examples/09_package_level
//
// Output that depends on the wall clock is shown as booleans so it is stable.
//
// Expected output:
//
//	BeginOfDay:      2024-03-15 00:00:00
//	EndOfDayOf:      2024-03-15 23:59:59
//	PrevBeginOfDay:  2024-03-13 00:00:00
//	PrevEndOfDay:    2024-03-13 23:59:59
//	AddSecond/Minute/Hour/Day: 14:31:15 15:00:45 2024-03-16 00:30 2024-04-04
//	IsLeapYear(2024), IsLeapYear(2100): true false
//	IsLeapYearOf(ref): true
//	weekdays in range: 5
//	DateParts: [0 45 30 14 15 3 2024]
//	Quarter in 1..4: true
//	SinceHour ~1.5:   true
//	SinceMinute ~90:  true
//	SinceSecond ~5400: true
//	IsWithinTolerance(now+30s), (now+5m): true false
//	current minute   contains now: true
//	current hour     contains now: true
//	current day      contains now: true
//	current week     contains now: true
//	current month    contains now: true
//	current quarter  contains now: true
//	current year     contains now: true
//	week start now: Monday
//	Monday() is Monday: true
//	Friday() is Friday: true
//	EndOfSunday() is Sunday: true
//	Parse: 2024-03-15 09:30:00
//	MustParse: 2024-03-15 09:30:00
//	ParseInLocation: 2024-03-15 09:30:00 JST
//	MustParseInLocation: 2024-03-15 09:30:00 JST
//	Between(2000-01-01, 2999-01-01): true
//	FormatRFCAt:       2024-03-15 14:30:45
//	FormatRFCshortAt:  14:30:45
//	DefaultFormatRFCAt: 2024-03-15 14:30:45
//	FormatRFC has date: true
//	FormatRFCshort has time: true
//	DefaultFormatRFC non-empty: true
package main

import (
	"fmt"
	"os"
	"time"

	"github.com/polarixa/timefy"
)

const layout = "2006-01-02 15:04:05"

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func run() error {
	ref := time.Date(2024, time.March, 15, 14, 30, 45, 0, time.UTC)

	fmt.Println("BeginOfDay:     ", timefy.BeginOfDay(ref).Format(layout))
	fmt.Println("EndOfDayOf:     ", timefy.EndOfDayOf(ref).Format(layout))
	fmt.Println("PrevBeginOfDay: ", timefy.PrevBeginOfDay(ref, 2).Format(layout))
	fmt.Println("PrevEndOfDay:   ", timefy.PrevEndOfDay(ref, 2).Format(layout))
	fmt.Println("AddSecond/Minute/Hour/Day:",
		timefy.AddSecond(ref, 30).Format("15:04:05"),
		timefy.AddMinute(ref, 30).Format("15:04:05"),
		timefy.AddHour(ref, 10).Format("2006-01-02 15:04"),
		timefy.AddDay(ref, 20).Format("2006-01-02"))

	// Leap years.
	fmt.Println("IsLeapYear(2024), IsLeapYear(2100):", timefy.IsLeapYear(2024), timefy.IsLeapYear(2100))
	fmt.Println("IsLeapYearOf(ref):", timefy.IsLeapYearOf(ref))

	// Weekdays in a range (Mon-Fri, inclusive).
	days := timefy.WeekdaysInRange(ref, ref.AddDate(0, 0, 6))
	fmt.Println("weekdays in range:", len(days))

	// DateParts: [nanosecond, second, minute, hour, day, month, year].
	fmt.Println("DateParts:", timefy.DateParts(ref))

	// Quarter number of the current time.
	q := timefy.Quarter()
	fmt.Println("Quarter in 1..4:", q >= 1 && q <= 4)

	// Elapsed time since a moment.
	hourAgo := time.Now().Add(-90 * time.Minute)
	fmt.Println("SinceHour ~1.5:  ", roughly(timefy.SinceHour(hourAgo), 1.5))
	fmt.Println("SinceMinute ~90: ", roughly(timefy.SinceMinute(hourAgo), 90))
	fmt.Println("SinceSecond ~5400:", roughly(timefy.SinceSecond(hourAgo), 5400))

	// IsWithinTolerance: within one minute of now.
	fmt.Println("IsWithinTolerance(now+30s), (now+5m):",
		timefy.IsWithinTolerance(time.Now().Add(30*time.Second)),
		timefy.IsWithinTolerance(time.Now().Add(5*time.Minute)))

	// Current-period boundaries are ordered: beginning <= now <= end.
	now := time.Now()
	periods := []struct {
		name       string
		begin, end time.Time
	}{
		{"minute", timefy.BeginningOfMinute(), timefy.EndOfMinute()},
		{"hour", timefy.BeginningOfHour(), timefy.EndOfHour()},
		{"day", timefy.BeginningOfDay(), timefy.EndOfDay()},
		{"week", timefy.BeginningOfWeek(), timefy.EndOfWeek()},
		{"month", timefy.BeginningOfMonth(), timefy.EndOfMonth()},
		{"quarter", timefy.BeginningOfQuarter(), timefy.EndOfQuarter()},
		{"year", timefy.BeginningOfYear(), timefy.EndOfYear()},
	}
	for _, p := range periods {
		fmt.Printf("current %-8s contains now: %v\n", p.name, !p.begin.After(now) && !p.end.Before(now))
	}

	// Week start is global: set it, read it, and restore the previous value.
	previous := timefy.GetWeekStartDay()
	timefy.SetWeekStartDay(time.Monday)
	fmt.Println("week start now:", timefy.GetWeekStartDay())
	timefy.SetWeekStartDay(previous)

	// Weekday helpers for the current week.
	fmt.Println("Monday() is Monday:", timefy.Monday().Weekday() == time.Monday)
	fmt.Println("Friday() is Friday:", timefy.Friday().Weekday() == time.Friday)
	fmt.Println("EndOfSunday() is Sunday:", timefy.EndOfSunday().Weekday() == time.Sunday)

	// Parse, MustParse and the location-aware variants.
	v, err := timefy.Parse("2024-03-15 09:30:00")
	if err != nil {
		return fmt.Errorf("parse: %w", err)
	}
	fmt.Println("Parse:", v.Format(layout))
	fmt.Println("MustParse:", timefy.MustParse("2024-03-15 09:30:00").Format(layout))

	tokyo, err := time.LoadLocation(string(timefy.DefaultTimezoneTokyo))
	if err != nil {
		return fmt.Errorf("load location: %w", err)
	}
	v, err = timefy.ParseInLocation(tokyo, "2024-03-15 09:30:00")
	if err != nil {
		return fmt.Errorf("parse in location: %w", err)
	}
	fmt.Println("ParseInLocation:", v.Format(layout+" MST"))
	fmt.Println("MustParseInLocation:", timefy.MustParseInLocation(tokyo, "2024-03-15 09:30:00").Format(layout+" MST"))

	// Between: is now strictly inside the range? (bounds are parsed strings)
	fmt.Println("Between(2000-01-01, 2999-01-01):", timefy.Between("2000-01-01", "2999-01-01"))

	// Formatting with RFC layouts.
	fmt.Println("FormatRFCAt:      ", timefy.FormatRFCAt(ref, timefy.TimeFormat20060102150405))
	fmt.Println("FormatRFCshortAt: ", timefy.FormatRFCshortAt(ref, timefy.TimeRFC01T150405))
	fmt.Println("DefaultFormatRFCAt:", timefy.DefaultFormatRFCAt(ref))
	fmt.Println("FormatRFC has date:", len(timefy.FormatRFC(timefy.TimeFormat20060102)) == len("2006-01-02"))
	fmt.Println("FormatRFCshort has time:", len(timefy.FormatRFCshort(timefy.TimeRFC01T150405)) == len("15:04:05"))
	fmt.Println("DefaultFormatRFC non-empty:", timefy.DefaultFormatRFC() != "")
	return nil
}

// roughly reports whether got is within 1% of want.
func roughly(got, want float64) bool {
	d := got - want
	if d < 0 {
		d = -d
	}
	return d <= want*0.01
}
