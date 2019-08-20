package timex

import (
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	"time"
)

//解析日期时间字符串并转换为proto.Timestamp
//比如：timex.ParseProtoTimestamp("2006-01-02T15:04:05.000Z", stringDateTimeRFC3339),
func ParseProtoTimestamp(layout, value string) *timestamp.Timestamp {
	tm, err := time.Parse(layout, value)
	if err != nil {
		tm = Zero
	}
	ts, _ := ptypes.TimestampProto(tm)
	return ts
}

//转换为proto.Timestamp
func ToProtoTimestamp(tm time.Time) *timestamp.Timestamp {
	ts, _ := ptypes.TimestampProto(tm)
	return ts
}

//转换为time.Time
func FromProtoTimestamp(ts *timestamp.Timestamp) time.Time {
	tm, _ := ptypes.Timestamp(ts)
	return tm.In(LocationAsiaShanghai)
}
