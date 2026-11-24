# Changelog

All notable changes to this project are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2026-11-24

### Added

- Added descriptive package-level helpers for the operations that previously
  required the legacy names: `DateParts`, `EndOfDayOf`, `FormatRFCAt`,
  `FormatRFCshortAt`, `DefaultFormatRFCAt`, `IsLeapYearOf`, `InTimezone`,
  `InTimezoneOrSame` and `WeekdaysInRange`. The new names state what they
  operate on and read consistently next to the `*Of`/`*At` family.
- Added 15 numbered, runnable programs under `examples/` (`01_quickstart` to
  `15_scheduler`) plus `examples/README.md`. Each program is a standalone
  `package main` that documents its exact expected output and runs with
  `go run ./examples/<name>`, so users can copy working code instead of
  reading tests.

### Changed

- **Breaking: `DurationInDays`, `DurationInWeeks`, `DurationInMonths` and
  `DurationInYears` now measure from the receiver to `other` and return whole
  calendar units.** `other` in the future is now positive; previously the sign
  was inverted (`tx.DurationInDays(future)` returned a negative number). Months
  and years now require the same day and time of day to be reached before the
  unit counts, so January 31 to February 1 is `0` months; previously only the
  month/year fields were compared. DST-affected days still count as one day.
  Migration: swap the arguments or negate the result, and re-check code that
  relied on partial months/years counting as a full unit.
- Renamed the exported helpers to the new names listed under Added. The old
  names (`FEndOfDay`, `FFormatRFC`, `FFormatRFCshort`, `FDefaultFormatRFC`,
  `IsLeapYearN`, `FormatTimex`, `SetTimezone`, `AdjustTimezone`,
  `GetWeekdaysInRange`) are kept as deprecated wrappers, so existing calls keep
  compiling; switch to the new names to avoid deprecation notices.
- `Parse` now evaluates every supplied string: a string that fails to parse is
  skipped instead of discarding a value parsed from another argument, so
  argument order no longer matters. `Parse()` with no arguments now returns an
  error instead of a zero `time.Time` and a nil error (**breaking** for callers
  that relied on the zero-value result).
- Rewrote `README.md` and `docs/guides/getting-started.md` around the public API, with
  a table of contents and links to the runnable examples, and corrected the
  documented weekday behavior (the returned day is inside the current
  Monday-to-Sunday week; the code itself is unchanged).
- Added self-documenting `make` targets, including `help` and CI-parity targets
  (`build-all`, `test-ci`, `check`, `race`, `bench`), so local runs match the
  pipeline.

### Fixed

- Fixed `Parse` replacing explicit zero components with the current time when
  the input also carried a zone or fractional seconds. For example,
  `2023-08-15T13:00:30.5Z` now keeps minute `0` and second `30` instead of
  inheriting the current minute.
- Fixed `BeginningOfHour` on the repeated hour of a DST fall-back. The boundary
  is now derived from the wall clock of the instant, so the second `01:30` in
  `America/New_York` maps to the `01:00` that contains it instead of the first
  occurrence, and the monotonic clock reading is stripped.
- Fixed `BeginningOfMinute` truncating the absolute instant instead of the wall
  clock, which left a non-zero second count in zones whose offset is not a whole
  number of minutes. It now zeroes the wall-clock seconds and nanoseconds.
- Fixed `WeekdaysInRange` excluding the end date when it fell earlier in the
  day than the start time. Dates are compared by calendar day, so the range
  `2024-10-28 08:00` to `2024-11-01 00:00` now returns all 5 weekdays.
- Fixed `TimeAgo` reporting a future time as `just now`; it now returns the
  forward-looking phrase produced by `TimeUntil`.
- Fixed `ParseInLocation` and `MustParseInLocation` panicking on a nil
  `*time.Location`; they now fall back to `time.Local`.
- Fixed `Parse` ignoring the `Rule` location when the `Timex` wrapped a time in
  another zone. Strings are now interpreted in the configured location
  regardless of how the `Timex` was built.
- Fixed `WithLocationRFC` silently falling back to the wrapped time's location
  when the zone name was invalid. The failure is recorded and returned by
  `Parse` as an error, so a bad zone can no longer produce times in an
  unexpected location.
- Fixed `WithTimeFormats`, `AppendTimeFormat`, `AppendTimeFormatRFC` and
  `AppendTimeFormatRFCshort` writing into the caller's backing array, which
  could corrupt a format slice reused across rules.

[1.0.0]: https://github.com/polarixa/timefy/compare/v0.0.8...v1.0.0
