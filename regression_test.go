package timefy_test

import (
	"testing"
	"time"

	"github.com/polarixa/timefy"
)

// F1: strings carrying a time must keep zero minutes/seconds even when a
// trailing zone or fraction is present.
func TestParseKeepsExplicitZeroTimeComponents(t *testing.T) {
	cases := map[string]time.Time{
		"2023-08-15T13:00:30.5Z":        time.Date(2023, 8, 15, 13, 0, 30, 500000000, time.UTC),
		"2023-08-15T13:00:00.5Z":        time.Date(2023, 8, 15, 13, 0, 0, 500000000, time.UTC),
		"Tue, 15 Aug 2023 13:00:30 UTC": time.Date(2023, 8, 15, 13, 0, 30, 0, time.UTC),
	}
	for in, want := range cases {
		got, err := timefy.With(time.Date(2030, 1, 1, 9, 47, 20, 0, time.UTC)).Parse(in)
		if err != nil {
			t.Fatalf("Parse(%q): %v", in, err)
		}
		if !got.Equal(want) {
			t.Errorf("Parse(%q) = %v; want %v", in, got, want)
		}
	}
}

// F2: a failing string must not discard a successful one, in either order.
func TestParseMultiArgOrderIndependent(t *testing.T) {
	want := time.Date(2023, 10, 25, 0, 0, 0, 0, time.UTC)
	tx := timefy.With(time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC))
	for _, args := range [][]string{{"2023-10-25", "garbage"}, {"garbage", "2023-10-25"}} {
		got, err := tx.Parse(args...)
		if err != nil || !got.Equal(want) {
			t.Errorf("Parse(%v) = %v, %v; want %v, nil", args, got, err, want)
		}
	}
	if _, err := tx.Parse(); err == nil {
		t.Error("Parse() with no arguments should return an error")
	}
	if _, err := tx.Parse("garbage"); err == nil {
		t.Error("Parse(garbage) should return an error")
	}
}

// F3: the repeated DST hour must resolve to the hour containing the instant.
func TestBeginningOfHourDSTFallBack(t *testing.T) {
	ny, err := time.LoadLocation("America/New_York")
	if err != nil {
		t.Skip("tzdata unavailable")
	}
	in := time.Date(2024, 11, 3, 6, 30, 0, 0, time.UTC).In(ny) // second 01:30 (EST)
	tx := timefy.With(in)
	b, e := tx.BeginningOfHour(), tx.EndOfHour()
	if b.After(in) || e.Before(in) {
		t.Errorf("hour bounds %v..%v do not contain %v", b, e, in)
	}
	if want := time.Date(2024, 11, 3, 6, 0, 0, 0, time.UTC); !b.Equal(want) {
		t.Errorf("BeginningOfHour = %v; want %v", b, want)
	}
}

// F5: an invalid zone must surface an error instead of falling back to local time.
func TestWithLocationRFCInvalid(t *testing.T) {
	if _, err := timefy.NewRule().WithLocationRFC("Not/AZone").Parse("10:00"); err == nil {
		t.Error("expected an error for an invalid location")
	}
}

// F7: future times are not "just now".
func TestTimeAgoFuture(t *testing.T) {
	if got := timefy.With(time.Now().AddDate(0, 0, 3)).TimeAgo(); got == "just now" {
		t.Errorf("TimeAgo(+3d) = %q", got)
	}
}

// F9: the end day is included regardless of time of day.
func TestGetWeekdaysInRangeIncludesEndDay(t *testing.T) {
	start := time.Date(2024, 10, 28, 8, 0, 0, 0, time.UTC)
	end := time.Date(2024, 11, 1, 0, 0, 0, 0, time.UTC)
	if got := timefy.GetWeekdaysInRange(start, end); len(got) != 5 {
		t.Errorf("got %d weekdays; want 5", len(got))
	}
}

// F10: a nil location must not panic.
func TestParseInLocationNil(t *testing.T) {
	if _, err := timefy.ParseInLocation(nil, "2024-01-01"); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

// S3: appending formats must not write into the caller's slice.
func TestAppendTimeFormatDoesNotAlias(t *testing.T) {
	base := make([]string, 1, 8)
	base[0] = "2006-01-02"
	timefy.NewRule().WithTimeFormats(base).AppendTimeFormat("X1")
	timefy.NewRule().WithTimeFormats(base).AppendTimeFormat("X2")
	if got := base[:2][1]; got != "" {
		t.Errorf("caller slice was mutated: %q", got)
	}
}

// F4/F6: durations are signed other-minus-t and calendar aware.
func TestDurationCalendarSemantics(t *testing.T) {
	ny, err := time.LoadLocation("America/New_York")
	if err != nil {
		t.Skip("tzdata unavailable")
	}
	utc := time.UTC
	now := time.Date(2024, 6, 15, 12, 0, 0, 0, utc)
	tx := timefy.With(now)
	if got := tx.DurationInDays(now.AddDate(0, 0, 1)); got != 1 {
		t.Errorf("days(+1d) = %d; want 1", got)
	}
	if got := tx.DurationInDays(now.AddDate(0, 0, -1)); got != -1 {
		t.Errorf("days(-1d) = %d; want -1", got)
	}
	if got := tx.DurationInWeeks(now.AddDate(0, 0, 7)); got != 1 {
		t.Errorf("weeks(+7d) = %d; want 1", got)
	}
	if got := tx.DurationInMonths(now.AddDate(0, 1, 0)); got != 1 {
		t.Errorf("months(+1m) = %d; want 1", got)
	}
	if got := tx.DurationInYears(now.AddDate(1, 0, 0)); got != 1 {
		t.Errorf("years(+1y) = %d; want 1", got)
	}
	// DST: a 23-hour calendar day still counts as one day.
	s := time.Date(2024, 3, 9, 12, 0, 0, 0, ny)
	if got := timefy.With(s).DurationInDays(time.Date(2024, 3, 10, 12, 0, 0, 0, ny)); got != 1 {
		t.Errorf("DST days = %d; want 1", got)
	}
	// Incomplete month / year are not counted.
	if got := timefy.With(time.Date(2024, 1, 31, 0, 0, 0, 0, utc)).DurationInMonths(time.Date(2024, 2, 1, 0, 0, 0, 0, utc)); got != 0 {
		t.Errorf("Jan31->Feb1 months = %d; want 0", got)
	}
	if got := timefy.With(time.Date(2023, 12, 31, 0, 0, 0, 0, utc)).DurationInYears(time.Date(2024, 1, 1, 0, 0, 0, 0, utc)); got != 0 {
		t.Errorf("Dec31->Jan1 years = %d; want 0", got)
	}
}

// S2: a Timex built from a Rule parses in the rule's location.
func TestRuleWithHonoursLocation(t *testing.T) {
	ny, err := time.LoadLocation("America/New_York")
	if err != nil {
		t.Skip("tzdata unavailable")
	}
	got := timefy.NewRule().WithLocation(ny).With(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)).MustParse("2024-05-05 10:00")
	if got.Location() != ny || got.Hour() != 10 {
		t.Errorf("got %v; want 10:00 in %v", got, ny)
	}
}
