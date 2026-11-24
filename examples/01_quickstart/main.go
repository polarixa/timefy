// Package main demonstrates the smallest working use of timefy: wrap a
// time.Time, move to period boundaries, and format the result.
//
// Run: go run ./examples/01_quickstart
//
// Expected output:
//
//	start of day: 2024-03-15 00:00:00
//	end of day:   2024-03-15 23:59:59.999999999
//	formatted:    2024-03-15 14:30:45
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
	ref := time.Date(2024, time.March, 15, 14, 30, 45, 0, time.UTC)
	t := timefy.With(ref)

	fmt.Println("start of day:", t.BeginningOfDay().Format("2006-01-02 15:04:05"))
	fmt.Println("end of day:  ", t.EndOfDay().Format("2006-01-02 15:04:05.999999999"))
	fmt.Println("formatted:   ", t.FormatRFC(timefy.TimeFormat20060102150405))
	return nil
}
