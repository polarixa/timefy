// Package main demonstrates the BeginningOf* and EndOf* methods of Timex for
// minute, hour, day, week, month, quarter, half-year and year.
//
// Run: go run ./examples/04_period_boundaries
//
// Expected output:
//
//	minute   2024-05-15 14:30:00 | 2024-05-15 14:30:59.999999999
//	hour     2024-05-15 14:00:00 | 2024-05-15 14:59:59.999999999
//	day      2024-05-15 00:00:00 | 2024-05-15 23:59:59.999999999
//	week     2024-05-12 00:00:00 | 2024-05-18 23:59:59.999999999
//	month    2024-05-01 00:00:00 | 2024-05-31 23:59:59.999999999
//	quarter  2024-04-01 00:00:00 | 2024-06-30 23:59:59.999999999
//	half     2024-01-01 00:00:00 | 2024-06-30 23:59:59.999999999
//	year     2024-01-01 00:00:00 | 2024-12-31 23:59:59.999999999
//	quarter number: 2
package main

import (
	"fmt"
	"os"
	"time"

	"github.com/polarixa/timefy"
)

const layout = "2006-01-02 15:04:05.999999999"

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func run() error {
	ref := time.Date(2024, time.May, 15, 14, 30, 45, 123, time.UTC) // a Wednesday
	t := timefy.NewRule().WithWeekStartDay(time.Sunday).With(ref)

	rows := []struct {
		name       string
		begin, end time.Time
	}{
		{"minute", t.BeginningOfMinute(), t.EndOfMinute()},
		{"hour", t.BeginningOfHour(), t.EndOfHour()},
		{"day", t.BeginningOfDay(), t.EndOfDay()},
		{"week", t.BeginningOfWeek(), t.EndOfWeek()},
		{"month", t.BeginningOfMonth(), t.EndOfMonth()},
		{"quarter", t.BeginningOfQuarter(), t.EndOfQuarter()},
		{"half", t.BeginningOfHalf(), t.EndOfHalf()},
		{"year", t.BeginningOfYear(), t.EndOfYear()},
	}
	for _, r := range rows {
		fmt.Printf("%-8s %s | %s\n", r.name, r.begin.Format("2006-01-02 15:04:05"), r.end.Format(layout))
	}
	fmt.Println("quarter number:", t.Quarter())
	return nil
}
