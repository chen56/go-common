package jsonx

import (
	"encoding/json"
	"github.com/chen56/go-common/timex"
	"strconv"
	"time"
)

type Map struct {
	Map map[string]interface{}
}

func NewMap(mp map[string]interface{}) *Map {
	return &Map{
		Map: mp,
	}
}

func Unmarshal2Map(data []byte) (*Map, error) {
	mp := make(map[string]interface{})
	err := json.Unmarshal(data, &mp)
	if err != nil {
		return nil, err
	}
	return &Map{
		Map: mp,
	}, nil
}

func (x *Map) Marshal() ([]byte, error) {
	return json.Marshal(x.Map)
}

//获取子项
func (x *Map) Get(key string) (*Map, bool) {
	mp, ok := x.Map[key].(map[string]interface{})
	if !ok {
		return nil, false
	}
	return &Map{
		Map: mp,
	}, true
}

//获取值
func (x *Map) GetValue(key string) (interface{}, bool) {
	v, ok := x.Map[key].(interface{})
	return v, ok
}

//设置值
func (x *Map) SetValue(key string, value interface{}) {
	x.Map[key] = value
}

//获取String
func (x *Map) GetString(key string) (string, bool) {
	v, ok := x.Map[key]
	if ok {
		return x.toString(v)
	}
	return "", false
}

//转换String
func (x *Map) toString(val interface{}) (string, bool) {
	s, ok := val.(string)
	if ok {
		return s, true
	}
	f, ok := val.(float64)
	if ok {
		s := strconv.FormatFloat(f, 'f', -1, 64)
		return s, true
	}
	return "", false
}

//获取StringList
func (x *Map) GetStringList(key string) ([]string, bool) {
	l, ok := x.Map[key].([]interface{})
	if !ok {
		return nil, false
	}
	list := make([]string, len(l))
	for _, v := range l {
		s, ok := x.toString(v)
		if !ok {
			return nil, false
		}
		list = append(list, s)
	}
	return list, true
}

//获取Int
func (x *Map) GetInt(key string) (int, bool) {
	v, ok := x.Map[key]
	if ok {
		return x.toInt(v)
	}
	return 0, false
}

//转换Int
func (x *Map) toInt(val interface{}) (int, bool) {
	f, ok := val.(float64)
	if ok {
		return int(f), true
	}
	s, ok := val.(string)
	if ok {
		i, err := strconv.Atoi(s)
		if err == nil {
			return i, true
		}
	}
	return 0, false
}

//获取IntList
func (x *Map) GetIntList(key string) ([]int, bool) {
	l, ok := x.Map[key].([]interface{})
	if !ok {
		return nil, false
	}
	list := make([]int, len(l))
	for _, v := range l {
		s, ok := x.toInt(v)
		if !ok {
			return nil, false
		}
		list = append(list, s)
	}
	return list, true
}

//获取Int32
func (x *Map) GetInt32(key string) (int32, bool) {
	v, ok := x.Map[key]
	if ok {
		return x.toInt32(v)
	}
	return 0, false
}

//转换Int32
func (x *Map) toInt32(val interface{}) (int32, bool) {
	f, ok := val.(float64)
	if ok {
		return int32(f), true
	}
	s, ok := val.(string)
	if ok {
		i, err := strconv.ParseInt(s, 10, 32)
		if err == nil {
			return int32(i), true
		}
	}
	return 0, false
}

//获取Int32List
func (x *Map) GetInt32List(key string) ([]int32, bool) {
	l, ok := x.Map[key].([]interface{})
	if !ok {
		return nil, false
	}
	list := make([]int32, len(l))
	for _, v := range l {
		s, ok := x.toInt32(v)
		if !ok {
			return nil, false
		}
		list = append(list, s)
	}
	return list, true
}

//获取Int64
func (x *Map) GetInt64(key string) (int64, bool) {
	v, ok := x.Map[key]
	if ok {
		return x.toInt64(v)
	}
	return 0, false
}

//转换Int64
func (x *Map) toInt64(val interface{}) (int64, bool) {
	f, ok := val.(float64)
	if ok {
		return int64(f), true
	}
	s, ok := val.(string)
	if ok {
		i, err := strconv.ParseInt(s, 10, 64)
		if err == nil {
			return i, true
		}
	}
	return 0, false
}

//获取Int64List
func (x *Map) GetInt64List(key string) ([]int64, bool) {
	l, ok := x.Map[key].([]interface{})
	if !ok {
		return nil, false
	}
	list := make([]int64, len(l))
	for _, v := range l {
		s, ok := x.toInt64(v)
		if !ok {
			return nil, false
		}
		list = append(list, s)
	}
	return list, true
}

//获取日期时间
func (x *Map) GetTime(key string, layout string) (time.Time, bool) {
	if layout == "unix" {
		f, ok := x.Map[key].(float64)
		if !ok {
			return timex.Zero, false
		}
		i := int64(f)
		if i > 4294967294 { //毫秒
			return time.Unix(int64(i/1000), (i%1000)*1e6), true
		}
		return time.Unix(i, 0), true
	}
	s, ok := x.Map[key].(string)
	if !ok {
		return timex.Zero, false
	}
	t, err := time.ParseInLocation(layout, s, timex.LocationAsiaShanghai)
	if err != nil {
		return timex.Zero, false
	}
	return t, true
}
