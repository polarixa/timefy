// Package main demonstrates Rule: custom time formats, week start day and
// locations, including how an invalid location is reported by Parse.
//
// Run: go run ./examples/03_rule_options
//
// Expected output:
//
//	custom format:   2024-03-15 00:00:00 +0000 UTC
//	appended format: 2024-03-15 00:00:00 +0000 UTC
//	week start Mon:  2024-03-11
//	week start Sun:  2024-03-10
//	tokyo:           2024-03-15 09:00:00 +0900 JST
//	paris:           2024-03-15 09:00:00 +0100 CET
//	local zone used: true
//	utc:             2024-03-15 09:00:00 +0000 UTC
//	bad location:    timefy: invalid location "Mars/Olympus": unknown time zone Mars/Olympus
package main

import (
	"fmt"
	"os"
	"time"

	"github.com/polarixa/timefy"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func run() error {
	ref := time.Date(2024, time.March, 15, 12, 0, 0, 0, time.UTC) // a Friday

	// WithTimeFormats replaces the accepted layouts.
	custom := timefy.NewRule().ApplyUTCLoc().WithTimeFormats([]string{"02/01/2006"})
	v, err := custom.Parse("15/03/2024")
	if err != nil {
		return fmt.Errorf("custom format: %w", err)
	}
	fmt.Println("custom format:  ", v)

	// AppendTimeFormat, AppendTimeFormatRFC and AppendTimeFormatRFCshort extend the defaults.
	extended := timefy.NewRule().ApplyUTCLoc().
		AppendTimeFormat("2006.01.02").
		AppendTimeFormatRFC(timefy.TimeFormat200601).
		AppendTimeFormatRFCshort(timefy.TimeRFC06T150405)
	v, err = extended.Parse("2024.03.15")
	if err != nil {
		return fmt.Errorf("appended format: %w", err)
	}
	fmt.Println("appended format:", v)

	// WithWeekStartDay changes what BeginningOfWeek returns.
	mon := timefy.NewRule().WithWeekStartDay(time.Monday).With(ref)
	sun := timefy.NewRule().WithWeekStartDay(time.Sunday).With(ref)
	fmt.Println("week start Mon: ", mon.BeginningOfWeek().Format("2006-01-02"))
	fmt.Println("week start Sun: ", sun.BeginningOfWeek().Format("2006-01-02"))

	// WithLocationRFC makes Parse interpret strings in that zone.
	tokyo := timefy.NewRule().WithLocationRFC(timefy.DefaultTimezoneTokyo)
	v, err = tokyo.Parse("2024-03-15 09:00:00")
	if err != nil {
		return fmt.Errorf("tokyo: %w", err)
	}
	fmt.Println("tokyo:          ", v)

	// WithLocation takes a *time.Location directly.
	paris, err := time.LoadLocation(string(timefy.DefaultTimezoneParis))
	if err != nil {
		return fmt.Errorf("load location: %w", err)
	}
	v, err = timefy.NewRule().WithLocation(paris).Parse("2024-03-15 09:00:00")
	if err != nil {
		return fmt.Errorf("paris: %w", err)
	}
	fmt.Println("paris:          ", v)

	// ApplyLocalLoc selects the machine's local zone (output depends on the machine).
	v, err = timefy.NewRule().ApplyLocalLoc().Parse("2024-03-15 09:00:00")
	if err != nil {
		return fmt.Errorf("local: %w", err)
	}
	fmt.Println("local zone used:", v.Location() == time.Local)

	// ApplyUTCLoc does the same for UTC.
	v, err = timefy.NewRule().ApplyUTCLoc().Parse("2024-03-15 09:00:00")
	if err != nil {
		return fmt.Errorf("utc: %w", err)
	}
	fmt.Println("utc:            ", v)

	// An invalid zone is remembered and reported by the next Parse call.
	if _, err = timefy.NewRule().WithLocationRFC("Mars/Olympus").Parse("2024-03-15"); err != nil {
		fmt.Println("bad location:   ", err)
	}
	return nil
}
