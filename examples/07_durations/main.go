// Package main demonstrates whole-unit calendar durations between two times:
// DurationInDays, DurationInWeeks, DurationInMonths and DurationInYears.
// Results are whole units (calendar-aware for months and years) and negative
// when the other time is earlier.
//
// Run: go run ./examples/07_durations
//
// Expected output:
//
//	days:   408
//	weeks:  58
//	months: 13
//	years:  1
//	days backwards: -408
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
	start := timefy.With(time.Date(2023, time.January, 31, 9, 0, 0, 0, time.UTC))
	end := time.Date(2024, time.March, 15, 8, 0, 0, 0, time.UTC)

	fmt.Println("days:  ", start.DurationInDays(end))
	fmt.Println("weeks: ", start.DurationInWeeks(end))
	fmt.Println("months:", start.DurationInMonths(end))
	fmt.Println("years: ", start.DurationInYears(end))

	// Reversing the direction negates the result.
	back := timefy.With(end)
	fmt.Println("days backwards:", back.DurationInDays(start.Time))
	return nil
}
