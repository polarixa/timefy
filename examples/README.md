# timefy examples

Runnable programs that use only the public `github.com/polarixa/timefy` API.
Each folder is a `package main`; every file starts with a header comment that
states what it shows and the exact output it prints.

## Prerequisites

- Go 1.23 or newer (see `go.mod`).
- No network, environment variables or external services.

## Focused examples (start here)

| #   | Case                                                 | What it shows                                                                                | Run                                      |
| --- | ---------------------------------------------------- | -------------------------------------------------------------------------------------------- | ---------------------------------------- |
| 01  | [01_quickstart](01_quickstart/main.go)               | Wrap a `time.Time`, get day boundaries, format.                                              | `go run ./examples/01_quickstart`        |
| 02  | [02_parsing](02_parsing/main.go)                     | `Parse`, `MustParse`, time-only input, error results.                                        | `go run ./examples/02_parsing`           |
| 03  | [03_rule_options](03_rule_options/main.go)           | `Rule`: formats, week start, locations, invalid-location error.                              | `go run ./examples/03_rule_options`      |
| 04  | [04_period_boundaries](04_period_boundaries/main.go) | `BeginningOf*` / `EndOf*` from minute to year, `Quarter`.                                    | `go run ./examples/04_period_boundaries` |
| 05  | [05_weekdays](05_weekdays/main.go)                   | `Monday`..`Sunday`, all `*Safe` variants, `EndOfSunday`.                                     | `go run ./examples/05_weekdays`          |
| 06  | [06_comparisons](06_comparisons/main.go)             | `Is*` checks and `Between`.                                                                  | `go run ./examples/06_comparisons`       |
| 07  | [07_durations](07_durations/main.go)                 | `DurationInDays/Weeks/Months/Years`, negative direction.                                     | `go run ./examples/07_durations`         |
| 08  | [08_relative_time](08_relative_time/main.go)         | `TimeAgo` and `TimeUntil`.                                                                   | `go run ./examples/08_relative_time`     |
| 09  | [09_package_level](09_package_level/main.go)         | Package-level helpers: `BeginOfDay`, `Add*`, `Since*`, week start, `Parse*`, RFC formatting. | `go run ./examples/09_package_level`     |
| 10  | [10_timezones](10_timezones/main.go)                 | `InTimezone` (strict), `InTimezoneOrSame` (lenient), `DefaultTimezone*`.                     | `go run ./examples/10_timezones`         |

## Scenario examples

Longer, domain-flavoured programs. Run each with `go run ./examples/<name>`.

| Case                               | Scenario                                                    |
| ---------------------------------- | ----------------------------------------------------------- |
| [audit](audit/main.go)             | Enterprise audit trail and GDPR breach-notification timing. |
| [billing](billing/main.go)         | SaaS subscription billing period calculations.              |
| [global_team](global_team/main.go) | Multi-timezone global team coordination.                    |
| [reporting](reporting/main.go)     | Financial quarter and fiscal-year reporting.                |
| [scheduler](scheduler/main.go)     | Business-day-aware job scheduler.                           |

## Running everything

```sh
for d in examples/*/; do go run "./$d" || exit 1; done
```

## Not demonstrated on purpose

The deprecated aliases `FEndOfDay`, `FFormatRFC`, `FFormatRFCshort`,
`FDefaultFormatRFC`, `IsLeapYearN`, `FormatTimex`, `SetTimezone`,
`AdjustTimezone` and `GetWeekdaysInRange` are kept for compatibility. The
examples use their replacements (`EndOfDayOf`, `FormatRFCAt`,
`FormatRFCshortAt`, `DefaultFormatRFCAt`, `IsLeapYearOf`, `DateParts`,
`InTimezone`, `InTimezoneOrSame`, `WeekdaysInRange`).
