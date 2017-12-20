package timex

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLocation(t *testing.T) {
	a := assert.New(t)

	x := time.Date(1980, 1, 1, 0, 0, 0, 0, Shanghai)

	a.Equal("1980-01-01T00:00:00+08:00", x.Format(time.RFC3339))

}

func TestB(t *testing.T) {
	birthday := time.Date(1980, 1, 2, 3, 30, 0, 0, time.UTC)

	fmt.Printf("You are %v years.",
		Diff(birthday, time.Now()).Format(time.RFC3339))

}
func TestA(t *testing.T) {
	a := assert.New(t)
	type testData struct {
		start, end string
		expected   string
	}
	tests := []testData{
		testData{
			start:    "2010-01-01T00:00:00Z",
			end:      "2000-01-01T00:00:00Z",
			expected: "0009-11-30T00:00:00Z",
		},
		testData{
			start:    "2010-01-02T00:00:00Z",
			end:      "2000-01-04T00:00:00Z",
			expected: "0009-11-29T00:00:00Z",
		},
	}

	for i, test := range tests {
		actual := Diff(toTime(test.start), toTime(test.end)).Format(time.RFC3339)
		a.Equal(test.expected, actual, "%v - %v", i, test)
	}
}

func toTime(t string) time.Time {
	result, err := time.Parse(time.RFC3339, t)
	if err != nil {
		panic(fmt.Sprintf("[%s]format error:%s", t, err.Error()))
	}
	return result
}
