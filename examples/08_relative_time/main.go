// Package main demonstrates TimeAgo and TimeUntil, which describe a Timex
// relative to the current time in human-readable form.
//
// Run: go run ./examples/08_relative_time
//
// Expected output:
//
//	5s        ago   -> just now
//	5m0s      ago   -> 5 minutes ago
//	3h0m0s    ago   -> 3 hours ago
//	72h0m0s   ago   -> 3 days ago
//	480h0m0s  ago   -> 2 weeks ago
//	9600h0m0s ago   -> 1 year ago
//	30s       ahead -> in a few seconds
//	5m0s      ahead -> in 5 minutes
//	3h0m0s    ahead -> in 3 hours
//	72h0m0s   ahead -> in 3 days
//	future via TimeAgo: in 2 hours
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
	now := time.Now()

	// A small margin keeps the boundaries stable while the program runs.
	past := []time.Duration{5 * time.Second, 5 * time.Minute, 3 * time.Hour, 3 * 24 * time.Hour, 20 * 24 * time.Hour, 400 * 24 * time.Hour}
	for _, d := range past {
		t := timefy.With(now.Add(-d - time.Second))
		fmt.Printf("%-9s ago   -> %s\n", d, t.TimeAgo())
	}

	future := []time.Duration{30 * time.Second, 5 * time.Minute, 3 * time.Hour, 3 * 24 * time.Hour}
	for _, d := range future {
		t := timefy.With(now.Add(d + time.Second))
		fmt.Printf("%-9s ahead -> %s\n", d, t.TimeUntil())
	}

	// TimeAgo on a future time falls back to the TimeUntil wording.
	fmt.Println("future via TimeAgo:", timefy.With(now.Add(2*time.Hour+time.Minute)).TimeAgo())
	return nil
}
