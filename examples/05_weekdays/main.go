// Package main demonstrates the weekday helpers: Monday..Sunday return that
// day of the reference week, the *Safe variants return an error instead of
// panicking, and EndOfSunday returns the last moment of the week's Sunday.
//
// Run: go run ./examples/05_weekdays
//
// Expected output:
//
//	monday:    2024-05-13 Mon
//	tuesday:   2024-05-14 Tue
//	wednesday: 2024-05-15 Wed
//	thursday:  2024-05-16 Thu
//	friday:    2024-05-17 Fri
//	saturday:  2024-05-18 Sat
//	sunday:    2024-05-19 Sun
//	end of sunday: 2024-05-19 23:59:59
//	friday of 2024-12-25: 2024-12-27 Fri
//	MondaySafe     2024-05-13 Mon
//	TuesdaySafe    2024-05-14 Tue
//	WednesdaySafe  2024-05-15 Wed
//	ThursdaySafe   2024-05-16 Thu
//	FridaySafe     2024-05-17 Fri
//	SaturdaySafe   2024-05-18 Sat
//	SundaySafe     2024-05-19 Sun
//	safe error: timefy: failed to parse date for Monday: timefy: can't parse string as time: "not a date"
package main

import (
	"fmt"
	"os"
	"time"

	"github.com/polarixa/timefy"
)

const dateOnly = "2006-01-02 Mon"

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func run() error {
	ref := time.Date(2024, time.May, 15, 14, 30, 0, 0, time.UTC) // a Wednesday
	t := timefy.With(ref)

	// Without arguments the helpers start from the beginning of the reference day.
	fmt.Println("monday:   ", t.Monday().Format(dateOnly))
	fmt.Println("tuesday:  ", t.Tuesday().Format(dateOnly))
	fmt.Println("wednesday:", t.Wednesday().Format(dateOnly))
	fmt.Println("thursday: ", t.Thursday().Format(dateOnly))
	fmt.Println("friday:   ", t.Friday().Format(dateOnly))
	fmt.Println("saturday: ", t.Saturday().Format(dateOnly))
	fmt.Println("sunday:   ", t.Sunday().Format(dateOnly))
	fmt.Println("end of sunday:", t.EndOfSunday().Format("2006-01-02 15:04:05"))

	// With a date string the week of that date is used.
	v, err := t.FridaySafe("2024-12-25")
	if err != nil {
		return fmt.Errorf("friday of Christmas week: %w", err)
	}
	fmt.Println("friday of 2024-12-25:", v.Format(dateOnly))

	// Every weekday has a Safe variant with the same shape.
	safe := []struct {
		name string
		fn   func(...string) (time.Time, error)
	}{
		{"MondaySafe", t.MondaySafe}, {"TuesdaySafe", t.TuesdaySafe},
		{"WednesdaySafe", t.WednesdaySafe}, {"ThursdaySafe", t.ThursdaySafe},
		{"FridaySafe", t.FridaySafe}, {"SaturdaySafe", t.SaturdaySafe},
		{"SundaySafe", t.SundaySafe},
	}
	for _, s := range safe {
		d, err := s.fn("2024-05-15")
		if err != nil {
			return fmt.Errorf("%s: %w", s.name, err)
		}
		fmt.Printf("%-14s %s\n", s.name, d.Format(dateOnly))
	}

	// The Safe variants report unparsable input as an error.
	if _, err = t.MondaySafe("not a date"); err != nil {
		fmt.Println("safe error:", err)
	}
	return nil
}
