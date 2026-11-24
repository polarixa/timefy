# timefy

[![Go Reference](https://pkg.go.dev/badge/github.com/polarixa/timefy.svg)](https://pkg.go.dev/github.com/polarixa/timefy)
[![Go Version](https://img.shields.io/github/go-mod/go-version/polarixa/timefy)](go.mod)
[![Release](https://img.shields.io/github/v/tag/polarixa/timefy?label=release)](https://github.com/polarixa/timefy/tags)
[![License](https://img.shields.io/badge/license-Apache%202.0-blue.svg)](LICENSE)

> A Go library of time utilities: period boundaries, weekdays, parsing, timezones, durations and relative time.

## Table of Contents

- [Overview](#overview)
- [Installation](#installation)
- [Quick Start](#quick-start)
- [Usage Examples](#usage-examples)
- [API Reference](#api-reference)
- [Testing](#testing)
- [Contributing](#contributing)
- [License](#license)

<a id="overview"></a>

<details>
<summary><strong>Overview</strong></summary>

**timefy** is a comprehensive Go library designed to enhance productivity by offering a wide range of time-related utilities. Whether you're dealing with date formatting, time zone conversions, or scheduling. Timefy provides a robust toolkit to simplify time management tasks in Go applications.

**Why it exists.** Go's `time` package gives you the primitives but not the everyday helpers: start/end of a period, week-start rules, lenient parsing, relative phrasing. timefy wraps `time.Time` in a `Timex` value (and a configurable `Rule`) so these are one call each, with package-level shortcuts that use the current time.

**Key features**

- Period boundaries from minute to year (`BeginningOf*` / `EndOf*`), including quarter and half-year.
- Weekday helpers (`Monday` … `Sunday`), configurable week start, weekday ranges.
- Flexible string parsing (`Parse`, `MustParse`, `ParseInLocation`) with extensible layouts.
- Timezone conversion, comparisons (`Is*`, `Between`), durations, `TimeAgo` / `TimeUntil`.

**Use it when** you do calendar-style arithmetic and parsing in Go services or CLIs. **Skip it when** you only need `time.Now()` and `Format`, or need a full calendar/recurrence engine (holidays, recurring rules), which timefy does not provide.

**Status.** Pre-1.0 (tags are `v0.0.x`); the API may still change between releases.

</details>

<a id="installation"></a>

<details>
<summary><strong>Installation</strong></summary>

**Prerequisites:** Go version 1.23 or higher.

To install, you can use the following commands based on your preference:

- For a specific version:

  ```bash
  go get github.com/polarixa/timefy@v0.0.3
  ```

- For the latest version:
  ```bash
  go get github.com/polarixa/timefy@latest
  ```

Import it in your code; `go build`/`go run`/`go test` fetch the dependency automatically with [Go's module support](https://go.dev/wiki/Modules#how-to-use-modules):

```go
import "github.com/polarixa/timefy"
```

</details>

<a id="quick-start"></a>

<details>
<summary><strong>Quick Start</strong></summary>

Requires Go 1.23+ on Linux, macOS or Windows. No network, environment variables or external services are needed.

```bash
mkdir demo && cd demo
go mod init demo
go get github.com/polarixa/timefy@latest
```

Create `main.go`:

```go
package main

import (
	"fmt"
	"time"

	"github.com/polarixa/timefy"
)

func main() {
	tx := timefy.New(time.Date(2025, 1, 15, 17, 51, 49, 0, time.UTC))
	fmt.Println(tx.BeginningOfMonth()) // 2025-01-01 00:00:00 +0000 UTC
	fmt.Println(tx.EndOfMonth())       // 2025-01-31 23:59:59.999999999 +0000 UTC
}
```

```bash
go run .
```

</details>

<a id="usage-examples"></a>

<details>
<summary><strong>Usage Examples</strong></summary>

<details>
<summary>Calculating the time based on current time</summary>

```go
package main

import (
	"fmt"
	"time"

	"github.com/polarixa/timefy"
)

func main() {
	// time based on current time
	tx := timefy.With(time.Now())        // 2024-10-25 21:31:00.690342
	fmt.Println(tx.BeginningOfMinute())  // 2024-10-25 21:31:00
	fmt.Println(tx.BeginningOfHour())    // 2024-10-25 21:00:00
	fmt.Println(tx.BeginningOfDay())     // 2024-10-25 00:00:00
	fmt.Println(tx.BeginningOfWeek())    // 2024-10-20 00:00:00
	fmt.Println(tx.BeginningOfMonth())   // 2024-10-01 00:00:00
	fmt.Println(tx.BeginningOfQuarter()) // 2024-10-01 00:00:00
	fmt.Println(tx.BeginningOfHalf())    // 2024-07-01 00:00:00
	fmt.Println(tx.BeginningOfYear())    // 2024-01-01 00:00:00
	fmt.Println(tx.EndOfMinute())        // 2024-10-25 21:31:59.999999999 +0700 +07
	fmt.Println(tx.EndOfHour())          // 2024-10-25 21:59:59.999999999 +0700 +07
	fmt.Println(tx.EndOfDay())           // 2024-10-25 23:59:59.999999999 +0700 +07
	fmt.Println(tx.EndOfWeek())          // 2024-10-26 23:59:59.999999999 +0700 +07
	fmt.Println(tx.EndOfMonth())         // 2024-10-31 23:59:59.999999999 +0700 +07
	fmt.Println(tx.EndOfQuarter())       // 2024-12-31 23:59:59.999999999 +0700 +07
	fmt.Println(tx.EndOfHalf())          // 2024-12-31 23:59:59.999999999 +0700 +07
	fmt.Println(tx.EndOfYear())          // 2024-12-31 23:59:59.999999999 +0700 +07

	fmt.Println(tx.EndOfWeek())          // 2024-10-28 23:59:59.999999999 +0700 +07
}
```

</details>

<details>
<summary>Calculating the time based on another time</summary>

```go
package main

import (
	"fmt"
	"time"

	"github.com/polarixa/timefy"
)

func main() {
	// time based on another time
	t := time.Date(2025, 01, 15, 17, 51, 49, 123456789, time.Local)
	tx := timefy.New(t)           // or timefy.With(t)
	fmt.Println(tx.EndOfMonth()) // 2025-01-31 23:59:59.999999999 +0700 +07
}
```

</details>

<details>
<summary>Calculating the time based on configuration</summary>

```go
package main

import (
	"fmt"
	"time"

	"github.com/polarixa/timefy"
)

func main() {
	// time based on configuration
	rule := timefy.NewRule().WithLocationRFC(timefy.DefaultTimezoneVietnam)
	t := time.Date(2025, 01, 15, 17, 51, 49, 123456789, time.Now().Location())
	fmt.Println(t)                              // 2025-01-15 17:51:49.123456789 +0700 +07
	fmt.Println(rule.With(t).BeginningOfWeek()) // 2025-01-12 00:00:00 +0700 +07
	v, _ := rule.Parse("2025-11-12 22:14:01")
	fmt.Println(v) // 2025-11-12 22:14:01 +0700 +07
}
```

</details>

<details>
<summary>Monday / Sunday</summary>

```go
package main

import (
	"fmt"
	"time"

	"github.com/polarixa/timefy"
)

func main() {
	// Monday / Sunday
	t := time.Date(2025, 01, 15, 17, 51, 49, 123456789, time.Now().Location())
	fmt.Println(timefy.Monday())                   // 2024-10-21 00:00:00 +0700 +07
	fmt.Println(timefy.Monday("15:35:34"))         // 2024-10-21 15:35:34 +0700 +07
	fmt.Println(timefy.Sunday())                   // 2024-10-27 00:00:00 +0700 +07
	fmt.Println(timefy.Sunday("16:45:34"))         // 2024-10-27 16:45:34 +0700 +07
	fmt.Println(timefy.EndOfSunday())              // 2024-10-27 23:59:59.999999999 +0700 +07
	fmt.Println(timefy.With(t).Monday())           // 2025-01-13 00:00:00 +0700 +07
	fmt.Println(timefy.With(t).Monday("15:35:34")) // 2025-01-13 15:35:34 +0700 +07
	fmt.Println(timefy.With(t).Sunday())           // 2025-01-19 00:00:00 +0700 +07
	fmt.Println(timefy.With(t).Sunday("16:45:34")) // 2025-01-19 16:45:34 +0700 +07
	fmt.Println(timefy.With(t).EndOfSunday())      // 2025-01-19 23:59:59.999999999 +0700 +07
}
```

</details>

<details>
<summary>Parse String to Time</summary>

```go
package main

import (
	"fmt"

	"github.com/polarixa/timefy"
)

func main() {
	// String to Time
	t, _ := timefy.Parse("2025")
	fmt.Println(t) // 2025-01-01 00:00:00 +0700 +07
	t, _ = timefy.Parse("2025-02")
	fmt.Println(t) // 2025-02-01 00:00:00 +0700 +07
	t, _ = timefy.Parse("2025-02-15")
	fmt.Println(t) // 2025-02-15 00:00:00 +0700 +07
	t, _ = timefy.Parse("11-15")
	fmt.Println(t) // 2024-11-15 00:00:00 +0700 +07
	t, _ = timefy.Parse("13:15")
	fmt.Println(t) // 2024-10-25 13:15:00 +0700 +07
	t, _ = timefy.Parse("13:15:45")
	fmt.Println(t) // 2024-10-25 13:15:45 +0700 +07
	t, _ = timefy.Parse("23")
	fmt.Println(t) // 2024-10-25 23:00:00 +0700 +07
	// MustParse must parse string to time or it will panic
	t = timefy.MustParse("11")
	fmt.Println(t) // 2024-10-25 11:00:00 +0700 +07
	t = timefy.MustParse("99:99") // panic: can't parse string as time: 99:99
	// fmt.Println(t)

	// Extend timefy to support more formats is quite easy, just update timefy.TimeFormats with other time layouts, e.g:
	timefy.TimeFormats = append(timefy.TimeFormats, "02 Jan 2006 15:04")
	fmt.Println(timefy.TimeFormats)
}
```

</details>

<details>
<summary>Parse in Location</summary>

```go
package main

import (
	"fmt"
	"time"

	"github.com/polarixa/timefy"
)

func main() {
	// ParseInLocation parses in a specific timezone
	t, err := timefy.ParseInLocation(time.UTC, "2023-10-25")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(t) // 2023-10-25 00:00:00 +0000 UTC

	// MustParseInLocation panics on error
	t = timefy.MustParseInLocation(time.UTC, "2023-10-25 14:30:00")
	fmt.Println(t) // 2023-10-25 14:30:00 +0000 UTC
}
```

</details>

<details>
<summary>Day Operations</summary>

```go
package main

import (
	"fmt"
	"time"

	"github.com/polarixa/timefy"
)

func main() {
	now := time.Date(2025, 01, 15, 17, 51, 49, 0, time.UTC)

	// BeginOfDay returns the start of the day (00:00:00)
	fmt.Println(timefy.BeginOfDay(now)) // 2025-01-15 00:00:00 +0000 UTC

	// EndOfDayOf returns the end of the day (23:59:59)
	fmt.Println(timefy.EndOfDayOf(now)) // 2025-01-15 23:59:59 +0000 UTC

	// PrevBeginOfDay returns the start of the day N days ago
	fmt.Println(timefy.PrevBeginOfDay(now, 2)) // 2025-01-13 00:00:00 +0000 UTC

	// PrevEndOfDay returns the end of the day N days ago
	fmt.Println(timefy.PrevEndOfDay(now, 2)) // 2025-01-13 23:59:59 +0000 UTC
}
```

</details>

<details>
<summary>Time Manipulation (Add/Subtract)</summary>

```go
package main

import (
	"fmt"
	"time"

	"github.com/polarixa/timefy"
)

func main() {
	now := time.Date(2025, 01, 15, 12, 0, 0, 0, time.UTC)

	// AddSecond adds seconds (negative to subtract)
	fmt.Println(timefy.AddSecond(now, 30))  // 2025-01-15 12:00:30 +0000 UTC
	fmt.Println(timefy.AddSecond(now, -30)) // 2025-01-15 11:59:30 +0000 UTC

	// AddMinute adds minutes
	fmt.Println(timefy.AddMinute(now, 10))  // 2025-01-15 12:10:00 +0000 UTC
	fmt.Println(timefy.AddMinute(now, -10)) // 2025-01-15 11:50:00 +0000 UTC

	// AddHour adds hours
	fmt.Println(timefy.AddHour(now, 3))  // 2025-01-15 15:00:00 +0000 UTC
	fmt.Println(timefy.AddHour(now, -3)) // 2025-01-15 09:00:00 +0000 UTC

	// AddDay adds days
	fmt.Println(timefy.AddDay(now, 5))  // 2025-01-20 12:00:00 +0000 UTC
	fmt.Println(timefy.AddDay(now, -5)) // 2025-01-10 12:00:00 +0000 UTC
}
```

</details>

<details>
<summary>Timezone Utilities</summary>

```go
package main

import (
	"fmt"
	"time"

	"github.com/polarixa/timefy"
)

func main() {
	now := time.Date(2025, 01, 15, 12, 0, 0, 0, time.UTC)

	// InTimezone converts time to a specific timezone (returns error if invalid)
	nyTime, err := timefy.InTimezone(now, "America/New_York")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(nyTime) // 2025-01-15 07:00:00 -0500 EST

	// InTimezoneOrSame converts time (returns original time on invalid timezone)
	tokyoTime := timefy.InTimezoneOrSame(now, "Asia/Tokyo")
	fmt.Println(tokyoTime) // 2025-01-15 21:00:00 +0900 JST

	invalidTZ := timefy.InTimezoneOrSame(now, "Invalid/Zone")
	fmt.Println(invalidTZ) // 2025-01-15 12:00:00 +0000 UTC (original time returned)
}
```

</details>

<details>
<summary>Leap Year & Tolerance Checks</summary>

```go
package main

import (
	"fmt"
	"time"

	"github.com/polarixa/timefy"
)

func main() {
	// IsLeapYear checks if a year is a leap year
	fmt.Println(timefy.IsLeapYear(2024)) // true
	fmt.Println(timefy.IsLeapYear(2023)) // false
	fmt.Println(timefy.IsLeapYear(1900)) // false
	fmt.Println(timefy.IsLeapYear(2000)) // true

	// IsLeapYearOf checks using a time.Time value
	t := time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC)
	fmt.Println(timefy.IsLeapYearOf(t)) // true

	// IsWithinTolerance checks if a time is within 1 minute of now
	now := time.Now()
	fmt.Println(timefy.IsWithinTolerance(now.Add(30 * time.Second)))  // true
	fmt.Println(timefy.IsWithinTolerance(now.Add(2 * time.Minute)))   // false
}
```

</details>

<details>
<summary>Between (String-based Range Check)</summary>

```go
package main

import (
	"fmt"

	"github.com/polarixa/timefy"
)

func main() {
	// Between checks if the current time is between two time strings
	isWithin := timefy.Between("2020-01-01", "2030-12-31")
	fmt.Println(isWithin) // true (current time is between these dates)

	isWithin = timefy.Between("2030-01-01", "2030-12-31")
	fmt.Println(isWithin) // false (current time is not in 2030 yet)
}
```

</details>

<details>
<summary>Quarter</summary>

```go
package main

import (
	"fmt"
	"time"

	"github.com/polarixa/timefy"
)

func main() {
	// Quarter returns the current quarter (1-4)
	fmt.Println(timefy.Quarter()) // e.g., 1 for January-March

	// Quarter on a Timex instance
	t := time.Date(2025, 7, 15, 0, 0, 0, 0, time.UTC)
	fmt.Println(timefy.New(t).Quarter()) // 3 (July is Q3)
}
```

</details>

<details>
<summary>Since (Elapsed Time Calculations)</summary>

```go
package main

import (
	"fmt"
	"time"

	"github.com/polarixa/timefy"
)

func main() {
	past := time.Now().Add(-2 * time.Hour)

	// SinceHour returns the number of hours passed
	fmt.Println(timefy.SinceHour(past)) // ~2.0

	// SinceMinute returns the number of minutes passed
	fmt.Println(timefy.SinceMinute(past)) // ~120.0

	// SinceSecond returns the number of seconds passed
	fmt.Println(timefy.SinceSecond(past)) // ~7200.0
}
```

</details>

<details>
<summary>Weekdays in Range</summary>

```go
package main

import (
	"fmt"
	"time"

	"github.com/polarixa/timefy"
)

func main() {
	start := time.Date(2025, 1, 13, 0, 0, 0, 0, time.UTC) // Monday
	end := time.Date(2025, 1, 19, 0, 0, 0, 0, time.UTC)   // Sunday

	// WeekdaysInRange returns all weekdays (Mon-Fri) in the range
	weekdays := timefy.WeekdaysInRange(start, end)
	for _, d := range weekdays {
		fmt.Println(d) // 2025-01-13, 2025-01-14, 2025-01-15, 2025-01-16, 2025-01-17
	}
}
```

</details>

<details>
<summary>DateParts (Extract Time Components)</summary>

```go
package main

import (
	"fmt"
	"time"

	"github.com/polarixa/timefy"
)

func main() {
	t := time.Date(2025, 3, 15, 14, 30, 45, 123456789, time.UTC)

	// DateParts returns a slice of time components: [nanosecond, second, minute, hour, day, month, year]
	components := timefy.DateParts(t)
	fmt.Println(components) // [123456789 45 30 14 15 3 2025]
}
```

</details>

<details>
<summary>Formatting (RFC Layouts)</summary>

```go
package main

import (
	"fmt"
	"time"

	"github.com/polarixa/timefy"
)

func main() {
	t := time.Date(2025, 8, 15, 13, 45, 30, 0, time.UTC)

	// FormatRFC formats the current time using a TimeFormatRFC layout
	fmt.Println(timefy.FormatRFC(timefy.TimeFormat20060102T150405))

	// FormatRFCshort formats the current time using a short TimeRFC layout
	fmt.Println(timefy.FormatRFCshort(timefy.TimeRFC01T150405))

	// FormatRFCAt formats a specific time with a TimeFormatRFC layout
	fmt.Println(timefy.FormatRFCAt(t, timefy.TimeFormat20060102T150405)) // 2025-08-15T13:45:30

	// FormatRFCshortAt formats a specific time with a short TimeRFC layout
	fmt.Println(timefy.FormatRFCshortAt(t, timefy.TimeRFC01T150405)) // 13:45:30

	// DefaultFormatRFC formats the current time with the default layout
	fmt.Println(timefy.DefaultFormatRFC())

	// DefaultFormatRFCAt formats a specific time with the default layout
	fmt.Println(timefy.DefaultFormatRFCAt(t))
}
```

</details>

<details>
<summary>Relative Time (TimeAgo / TimeUntil)</summary>

```go
package main

import (
	"fmt"
	"time"

	"github.com/polarixa/timefy"
)

func main() {
	// TimeAgo returns a human-readable string for past time
	fmt.Println(timefy.New(time.Now().Add(-5 * time.Second)).TimeAgo())  // "just now"
	fmt.Println(timefy.New(time.Now().Add(-5 * time.Minute)).TimeAgo()) // "5 minutes ago"
	fmt.Println(timefy.New(time.Now().Add(-2 * time.Hour)).TimeAgo())   // "2 hours ago"
	fmt.Println(timefy.New(time.Now().Add(-48 * time.Hour)).TimeAgo())  // "2 days ago"
	fmt.Println(timefy.New(time.Now().Add(-720 * time.Hour)).TimeAgo()) // "1 month ago"

	// TimeUntil returns a human-readable string for future time
	fmt.Println(timefy.New(time.Now().Add(5 * time.Second)).TimeUntil())  // "in a few seconds"
	fmt.Println(timefy.New(time.Now().Add(5 * time.Minute)).TimeUntil()) // "in 5 minutes"
	fmt.Println(timefy.New(time.Now().Add(2 * time.Hour)).TimeUntil())   // "in 2 hours"
	fmt.Println(timefy.New(time.Now().Add(48 * time.Hour)).TimeUntil())  // "in 2 days"

	// TimeUntil with past time
	fmt.Println(timefy.New(time.Now().Add(-1 * time.Hour)).TimeUntil()) // "in the past"
}
```

</details>

<details>
<summary>Duration Calculations</summary>

```go
package main

import (
	"fmt"
	"time"

	"github.com/polarixa/timefy"
)

func main() {
	now := time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC)
	tx := timefy.New(now)

	// DurationInDays returns the number of days between two times
	future := now.AddDate(0, 0, 10)
	fmt.Println(tx.DurationInDays(future)) // 10

	// DurationInWeeks returns the number of weeks between two times
	future = now.AddDate(0, 0, 21)
	fmt.Println(tx.DurationInWeeks(future)) // 3

	// DurationInMonths returns the number of months between two times
	future = now.AddDate(0, 3, 0)
	fmt.Println(tx.DurationInMonths(future)) // 3

	// DurationInYears returns the number of years between two times
	future = now.AddDate(2, 0, 0)
	fmt.Println(tx.DurationInYears(future)) // 2
}
```

</details>

<details>
<summary>Day Type Checks (Weekend / Weekday)</summary>

```go
package main

import (
	"fmt"
	"time"

	"github.com/polarixa/timefy"
)

func main() {
	saturday := time.Date(2025, 1, 18, 12, 0, 0, 0, time.UTC)
	monday := time.Date(2025, 1, 20, 12, 0, 0, 0, time.UTC)

	// IsWeekend checks if the time is Saturday or Sunday
	fmt.Println(timefy.New(saturday).IsWeekend()) // true
	fmt.Println(timefy.New(monday).IsWeekend())   // false

	// IsWeekday checks if the time is Monday to Friday
	fmt.Println(timefy.New(saturday).IsWeekday()) // false
	fmt.Println(timefy.New(monday).IsWeekday())   // true
}
```

</details>

<details>
<summary>Time Comparison Helpers</summary>

```go
package main

import (
	"fmt"
	"time"

	"github.com/polarixa/timefy"
)

func main() {
	now := time.Now()
	tx := timefy.New(now)

	yesterday := now.AddDate(0, 0, -1)
	tomorrow := now.AddDate(0, 0, 1)

	// IsBefore checks if the time is before another time
	fmt.Println(tx.IsBefore(tomorrow))  // true
	fmt.Println(tx.IsBefore(yesterday)) // false

	// IsAfter checks if the time is after another time
	fmt.Println(tx.IsAfter(yesterday)) // true
	fmt.Println(tx.IsAfter(tomorrow))  // false

	// IsBetween checks if the time falls within a range (inclusive)
	fmt.Println(tx.IsBetween(yesterday, tomorrow)) // true

	// IsSameDay checks if two times are on the same day
	fmt.Println(tx.IsSameDay(now))       // true
	fmt.Println(tx.IsSameDay(yesterday)) // false

	// IsSameMonth checks if two times are in the same month
	nextMonth := now.AddDate(0, 1, 0)
	fmt.Println(tx.IsSameMonth(now))       // true
	fmt.Println(tx.IsSameMonth(nextMonth)) // false

	// IsSameYear checks if two times are in the same year
	nextYear := now.AddDate(1, 0, 0)
	fmt.Println(tx.IsSameYear(now))      // true
	fmt.Println(tx.IsSameYear(nextYear)) // false
}
```

</details>

<details>
<summary>Temporal Checks (Today / Yesterday / Tomorrow / Past / Future)</summary>

```go
package main

import (
	"fmt"
	"time"

	"github.com/polarixa/timefy"
)

func main() {
	now := time.Now()

	// IsToday checks if the time is today
	fmt.Println(timefy.New(now).IsToday()) // true

	// IsYesterday checks if the time is yesterday
	yesterday := now.AddDate(0, 0, -1)
	fmt.Println(timefy.New(yesterday).IsYesterday()) // true

	// IsTomorrow checks if the time is tomorrow
	tomorrow := now.AddDate(0, 0, 1)
	fmt.Println(timefy.New(tomorrow).IsTomorrow()) // true

	// IsPast checks if the time is in the past
	past := now.Add(-1 * time.Hour)
	fmt.Println(timefy.New(past).IsPast()) // true

	// IsFuture checks if the time is in the future
	future := now.Add(1 * time.Hour)
	fmt.Println(timefy.New(future).IsFuture()) // true
}
```

</details>

<details>
<summary>Global Convenience Functions (BeginningOf / EndOf)</summary>

```go
package main

import (
	"fmt"

	"github.com/polarixa/timefy"
)

func main() {
	// These global functions use the current time automatically
	fmt.Println(timefy.BeginningOfMinute())  // Current time at start of minute
	fmt.Println(timefy.BeginningOfHour())    // Current time at start of hour
	fmt.Println(timefy.BeginningOfDay())     // Today at 00:00:00
	fmt.Println(timefy.BeginningOfWeek())    // Start of current week
	fmt.Println(timefy.BeginningOfMonth())   // 1st of current month at 00:00:00
	fmt.Println(timefy.BeginningOfQuarter()) // 1st of current quarter
	fmt.Println(timefy.BeginningOfYear())    // Jan 1st of current year

	fmt.Println(timefy.EndOfMinute())  // Current time at end of minute
	fmt.Println(timefy.EndOfHour())    // Current time at end of hour
	fmt.Println(timefy.EndOfDay())     // Today at 23:59:59.999999999
	fmt.Println(timefy.EndOfWeek())    // End of current week
	fmt.Println(timefy.EndOfMonth())   // Last day of current month at 23:59:59.999999999
	fmt.Println(timefy.EndOfQuarter()) // Last day of current quarter
	fmt.Println(timefy.EndOfYear())    // Dec 31st at 23:59:59.999999999
}
```

</details>

<details>
<summary>Configuration (Rule) — Complete API</summary>

```go
package main

import (
	"fmt"
	"time"

	"github.com/polarixa/timefy"
)

func main() {
	// NewRule creates a rule with default configuration
	rule := timefy.NewRule()

	// WithWeekStartDay sets the week start day
	rule.WithWeekStartDay(time.Monday)

	// WithTimeFormats overrides the time formats
	rule.WithTimeFormats([]string{"2006-01-02 15:04:05"})

	// WithLocation sets the location using *time.Location
	rule.WithLocation(time.UTC)

	// WithLocationRFC sets location using a predefined zone string
	rule.WithLocationRFC(timefy.DefaultTimezoneVietnam)

	// ApplyUTCLoc sets location to UTC
	rule.ApplyUTCLoc()

	// ApplyLocalLoc sets location to Local
	rule.ApplyLocalLoc()

	// AppendTimeFormat appends custom format strings
	rule.AppendTimeFormat("02 Jan 2006 15:04")

	// AppendTimeFormatRFC appends defined RFC format constants
	rule.AppendTimeFormatRFC(timefy.TimeFormat20060102T150405)

	// AppendTimeFormatRFCshort appends defined short RFC format constants
	rule.AppendTimeFormatRFCshort(timefy.TimeRFC01T150405)

	// Rule.With creates a Timex with the rule's configuration
	tx := rule.With(time.Now())
	fmt.Println(tx.BeginningOfWeek())

	// Rule.Parse parses a string using the rule's configuration
	v, _ := rule.Parse("2025-11-12 22:14:01")
	fmt.Println(v)

	// Rule.MustParse panics on error
	v = rule.MustParse("2025-11-12 22:14:01")
	fmt.Println(v)
}
```

</details>

<details>
<summary>Runnable programs in <code>examples/</code></summary>

Runnable programs live in [`examples/`](examples/README.md): focused API demos (`01`–`10`) and longer scenarios (`11`–`15`: billing, global team, audit, reporting, scheduler).

```bash
go run ./examples/01_quickstart
```

</details>

</details>

<a id="api-reference"></a>

<details>
<summary><strong>API / Configuration Reference</strong></summary>

The package-level API, grouped. Run `go doc -all github.com/polarixa/timefy` for full signatures.

| Area                             | Exported names                                                                                                                                                                                                                                                                  |
| -------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Construction                     | `New`, `With`, `NewRule`, `Rule` (`With`, `Parse`, `MustParse`, `WithWeekStartDay`, `WithTimeFormats`, `WithLocation`, `WithLocationRFC`, `ApplyUTCLoc`, `ApplyLocalLoc`, `AppendTimeFormat*`)                                                                                  |
| Period boundaries (current time) | `BeginningOfMinute/Hour/Day/Week/Month/Quarter/Year`, `EndOfMinute/Hour/Day/Week/Month/Quarter/Year`, `Quarter`, `BeginOfDay`, `EndOfDayOf`, `PrevBeginOfDay`, `PrevEndOfDay`                                                                                                   |
| Weekdays                         | `Monday` … `Sunday`, `EndOfSunday`, `WeekdaysInRange`, `GetWeekdaysInRange`, `SetWeekStartDay`, `GetWeekStartDay`, `WeekStartDay`                                                                                                                                               |
| Parsing                          | `Parse`, `MustParse`, `ParseInLocation`, `MustParseInLocation`, `TimeFormats`, `TimeFormatRegexp`                                                                                                                                                                               |
| Arithmetic and elapsed           | `AddSecond`, `AddMinute`, `AddHour`, `AddDay`, `SinceSecond`, `SinceMinute`, `SinceHour`                                                                                                                                                                                        |
| Timezones                        | `InTimezone`, `InTimezoneOrSame`, `AdjustTimezone`, `SetTimezone`, `ZoneRFC` constants such as `DefaultTimezoneVietnam`                                                                                                                                                         |
| Checks                           | `Between`, `IsLeapYear`, `IsLeapYearN`, `IsLeapYearOf`, `IsWithinTolerance`                                                                                                                                                                                                     |
| Formatting                       | `FormatRFC`, `FormatRFCAt`, `FormatRFCshort`, `FormatRFCshortAt`, `DefaultFormatRFC`, `DefaultFormatRFCAt`, `FormatTimex`, `DateParts`, `TimeFormatRFC` / `TimeRFC` constants                                                                                                   |
| `Timex` methods                  | `BeginningOf*`, `EndOf*`, `Monday`…`Sunday`, `Quarter`, `Is*` (`IsBefore`, `IsAfter`, `IsBetween`, `IsSameDay/Month/Year`, `IsToday`, `IsYesterday`, `IsTomorrow`, `IsPast`, `IsFuture`, `IsWeekend`, `IsWeekday`), `DurationInDays/Weeks/Months/Years`, `TimeAgo`, `TimeUntil` |

Deprecated aliases (`FEndOfDay`, `FFormatRFC`, `FFormatRFCshort`, `FDefaultFormatRFC`) remain exported; prefer the non-`F` forms.

</details>

<a id="testing"></a>

<details>
<summary><strong>Testing &amp; Benchmarks</strong></summary>

The repository ships unit and regression tests (`time_test.go`, `regression_test.go`). No benchmarks are defined at the moment.

Run `make help` to list all targets. Common ones:

```bash
make test   # run all tests with coverage
make lint   # go vet
make check  # fmt-check, lint and test
```

</details>

<a id="contributing"></a>

<details>
<summary><strong>Contributing</strong></summary>

To contribute to project, follow these steps:

1. Clone the repository:

   ```bash
   git clone --depth 1 https://github.com/polarixa/timefy.git
   ```

2. Navigate to the project directory:

   ```bash
   cd timefy
   ```

3. Prepare the project environment:
   ```bash
   go mod tidy
   ```

There is no `CONTRIBUTING.md` yet; issues and pull requests are welcome. Run `make check` before submitting.

</details>

<a id="license"></a>

<details>
<summary><strong>License</strong></summary>

Licensed under the [Apache License 2.0](LICENSE).

</details>
