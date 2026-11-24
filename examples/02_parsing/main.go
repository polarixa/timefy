// Package main demonstrates Timex.Parse, Timex.MustParse and Rule.Parse.
// A time-only input is completed with the date of the reference time held by
// the Timex.
//
// Run: go run ./examples/02_parsing
//
// Expected output:
//
//	full datetime: 2024-03-15 09:30:00
//	time only:     2024-03-15 18:45:00
//	invalid:       timefy: can't parse string as time: "not a date"
//	no input:      timefy: no date string provided
//	must parse:    2024-03-15 09:30:00
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
	ref := time.Date(2024, time.March, 15, 12, 0, 0, 0, time.UTC)
	t := timefy.With(ref)

	v, err := t.Parse("2024-03-15 09:30:00")
	if err != nil {
		return fmt.Errorf("parse full datetime: %w", err)
	}
	fmt.Println("full datetime:", v.Format(layout))

	// Time only: the date comes from the reference time.
	v, err = t.Parse("18:45:00")
	if err != nil {
		return fmt.Errorf("parse time only: %w", err)
	}
	fmt.Println("time only:    ", v.Format(layout))

	// Errors are returned, never swallowed.
	if _, err = t.Parse("not a date"); err != nil {
		fmt.Println("invalid:      ", err)
	}
	if _, err = t.Parse(); err != nil {
		fmt.Println("no input:     ", err)
	}

	// MustParse panics on bad input: use it only with trusted constants.
	fmt.Println("must parse:   ", t.MustParse("2024-03-15 09:30:00").Format(layout))
	return nil
}
