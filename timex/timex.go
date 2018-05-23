package timex

import (
	"time"

	"github.com/chen56/go-common/assert"
)

var Shanghai = MustLoadLocation("Asia/Shanghai")

func MustParse(layout, value string) time.Time {
	t, err := time.Parse(layout, value)
	assert.NoErr(err)
	return t
}
func MustParseInLocation(layout, value string,loc *time.Location) time.Time {
	t, err := time.ParseInLocation(layout, value,loc)
	assert.NoErr(err)
	return t
}
func MustParseDuration(s string) time.Duration {
	t, err := time.ParseDuration(s)
	assert.NoErr(err)
	return t
}
func MustLoadLocation(name string) *time.Location {
	t, err := time.LoadLocation(name)
	assert.NoErr(err)
	return t
}

func BeginOfDay(now time.Time) time.Time {
	result := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	return result
}
func BeginningOfMonth(now time.Time) time.Time {
	result := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	return result
}

type TimeDiff struct {
	year, month, day, hour, min, sec int
}

func Diff(a, b time.Time) time.Time {
	if a.Location() != b.Location() {
		b = b.In(a.Location())
	}
	if a.After(b) {
		a, b = b, a
	}
	y1, M1, d1 := a.Date()
	y2, M2, d2 := b.Date()

	h1, m1, s1 := a.Clock()
	h2, m2, s2 := b.Clock()

	var year, day, hour, min, sec int
	var month time.Month
	year = int(y2 - y1)
	month = time.Month(M2 - M1)
	day = int(d2 - d1)
	hour = int(h2 - h1)
	min = int(m2 - m1)
	sec = int(s2 - s1)

	// Normalize negative values
	if sec < 0 {
		sec += 60
		min--
	}
	if min < 0 {
		min += 60
		hour--
	}
	if hour < 0 {
		hour += 24
		day--
	}
	if day < 0 {
		// days in month:
		t := time.Date(y1, M1, 32, 0, 0, 0, 0, time.UTC)
		day += 32 - t.Day()
		month--
	}
	if month < 0 {
		month += 12
		year--
	}
	return time.Date(year, month, day, hour, min, sec, 0, a.Location())
}
