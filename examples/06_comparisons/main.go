// Package main demonstrates the comparison helpers of Timex: IsBefore,
// IsAfter, IsBetween, IsSameDay/Month/Year, IsWeekend/IsWeekday, Between and
// the relative-to-now checks IsToday, IsYesterday, IsTomorrow, IsPast, IsFuture.
//
// Run: go run ./examples/06_comparisons
//
// Expected output:
//
//	IsBefore(later):     true
//	IsAfter(earlier):    true
//	IsBetween(e, l):     true
//	IsSameDay(later):    true
//	IsSameMonth(Mar 1):  true
//	IsSameYear(2023):    false
//	IsWeekend:           true
//	IsWeekday:           false
//	Between(2024-03-01, 2024-03-31): true
//	now IsToday:         true
//	now+24h IsTomorrow:  true
//	now-24h IsYesterday: true
//	2024 IsPast:         true
//	2999 IsFuture:       true
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
	ref := time.Date(2024, time.March, 16, 10, 0, 0, 0, time.UTC) // a Saturday
	t := timefy.With(ref)
	earlier := ref.Add(-time.Hour)
	later := ref.Add(time.Hour)

	fmt.Println("IsBefore(later):    ", t.IsBefore(later))
	fmt.Println("IsAfter(earlier):   ", t.IsAfter(earlier))
	fmt.Println("IsBetween(e, l):    ", t.IsBetween(earlier, later))
	fmt.Println("IsSameDay(later):   ", t.IsSameDay(later))
	fmt.Println("IsSameMonth(Mar 1): ", t.IsSameMonth(time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC)))
	fmt.Println("IsSameYear(2023):   ", t.IsSameYear(time.Date(2023, 3, 16, 0, 0, 0, 0, time.UTC)))
	fmt.Println("IsWeekend:          ", t.IsWeekend())
	fmt.Println("IsWeekday:          ", t.IsWeekday())

	// Between parses its bounds with the same rules as Parse (strict comparison).
	fmt.Println("Between(2024-03-01, 2024-03-31):", t.Between("2024-03-01", "2024-03-31"))

	// Relative-to-now checks use the wall clock, so they are shown against "now".
	now := timefy.With(time.Now())
	fmt.Println("now IsToday:        ", now.IsToday())
	fmt.Println("now+24h IsTomorrow: ", timefy.With(time.Now().AddDate(0, 0, 1)).IsTomorrow())
	fmt.Println("now-24h IsYesterday:", timefy.With(time.Now().AddDate(0, 0, -1)).IsYesterday())
	fmt.Println("2024 IsPast:        ", t.IsPast())
	fmt.Println("2999 IsFuture:      ", timefy.With(time.Date(2999, 1, 1, 0, 0, 0, 0, time.UTC)).IsFuture())
	return nil
}
