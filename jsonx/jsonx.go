package jsonx

import (
	"encoding/json"
	"github.com/chen56/go-common/must"
)

func MustMarshal(v interface{}) []byte {
	b, err := json.Marshal(v)
	must.NoError(err)
	return b
}

func MustMarshal2String(v interface{}) string {
	b, err := json.Marshal(v)
	must.NoError(err)
	return string(b)
}

func MustMarshalIndent(v interface{}, prefix, indent string) []byte {
	b, err := json.MarshalIndent(v, prefix, indent)
	must.NoError(err)
	return b
}

func MustMarshalIndent2String(v interface{}, prefix, indent string) string {
	b, err := json.MarshalIndent(v, prefix, indent)
	must.NoError(err)
	return string(b)
}

// MustUnmarshal > json.Unmarshal
func MustUnmarshal(data []byte, v interface{}) {
	err := json.Unmarshal(data, v)
	must.NoError(err)
}

//SafePretty 打印出有缩进的json, 只为测试目的，忽略出错
func SafePretty(v interface{}) string {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err.Error()
	}
	return string(b)
}
