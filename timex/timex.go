package timex

import (
	"github.com/chen56/go-common/must"
	"time"
)

var LocationAsiaShanghai = MustLoadLocation("Asia/Shanghai")
var Zero = time.Date(0, 1, 1, 0, 0, 0, 0, time.UTC)

func MustParse(layout, value string) time.Time {
	t, err := time.Parse(layout, value)
	must.NoError(err)
	return t
}
func MustParseInLocation(layout, value string, loc *time.Location) time.Time {
	t, err := time.ParseInLocation(layout, value, loc)
	must.NoError(err)
	return t
}
func MustParseDuration(s string) time.Duration {
	t, err := time.ParseDuration(s)
	must.NoError(err)
	return t
}
func MustLoadLocation(name string) *time.Location {
	t, err := time.LoadLocation(name)
	must.NoError(err)
	return t
}

func BeginOfDay(now time.Time) time.Time {
	result := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	return result
}

// 获取周一
func BeginOfWeek(now time.Time) time.Time {
	NumOfWeek := NumOfWeek(now)
	res := now.AddDate(0, 0, 1-NumOfWeek)
	result := time.Date(res.Year(), res.Month(), res.Day(), 0, 0, 0, 0, now.Location())
	return result
}
func BeginOfMonth(now time.Time) time.Time {
	result := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	return result
}

// 获取周中数字
func NumOfWeek(now time.Time) int {
	Num := 0
	switch now.Weekday().String() {
	case "Monday":
		Num = 1
	case "Tuesday":
		Num = 2
	case "Wednesday":
		Num = 3
	case "Thursday":
		Num = 4
	case "Friday":
		Num = 5
	case "Saturday":
		Num = 6
	case "Sunday":
		Num = 7
	}
	return Num
}

//计算两个日期是否相等
func IsSameDate(t1, t2 time.Time) bool {
	y1, m1, d1 := t1.Date()
	y2, m2, d2 := t2.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}

//计算两个日期相差多少天
func DiffDays(t1, t2 time.Time) int {
	y1, m1, d1 := t1.Date()
	y2, m2, d2 := t2.Date()
	t1 = time.Date(y1, m1, d1, 0, 0, 0, 0, time.Local)
	t2 = time.Date(y2, m2, d2, 0, 0, 0, 0, time.Local)
	return int(t2.Sub(t1).Hours() / 24)
}
