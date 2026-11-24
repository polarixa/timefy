// Package main demonstrates InTimezone (strict) and InTimezoneOrSame
// (lenient) together with the DefaultTimezone* constants.
//
// Run: go run ./examples/10_timezones
//
// Expected output:
//
//	America/New_York   2024-07-01 10:00 EDT
//	Europe/London      2024-07-01 15:00 BST
//	Asia/Tokyo         2024-07-01 23:00 JST
//	Asia/Ho_Chi_Minh   2024-07-01 21:00 +07
//	strict:   unknown time zone Mars/Olympus
//	lenient:  2024-07-01 14:00 UTC
package main

import (
	"fmt"
	"os"
	"time"

	"github.com/polarixa/timefy"
)

const layout = "2006-01-02 15:04 MST"

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func run() error {
	meeting := time.Date(2024, time.July, 1, 14, 0, 0, 0, time.UTC)

	zones := []timefy.ZoneRFC{
		timefy.DefaultTimezoneNewYork,
		timefy.DefaultTimezoneLondon,
		timefy.DefaultTimezoneTokyo,
		timefy.DefaultTimezoneVietnam,
	}
	for _, z := range zones {
		local, err := timefy.InTimezone(meeting, string(z))
		if err != nil {
			return fmt.Errorf("convert to %s: %w", z, err)
		}
		fmt.Printf("%-18s %s\n", z, local.Format(layout))
	}

	// InTimezone reports an invalid zone; InTimezoneOrSame returns the input unchanged.
	if _, err := timefy.InTimezone(meeting, "Mars/Olympus"); err != nil {
		fmt.Println("strict:  ", err)
	}
	fmt.Println("lenient: ", timefy.InTimezoneOrSame(meeting, "Mars/Olympus").Format(layout))
	return nil
}
