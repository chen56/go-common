package jsonx

import (
	"encoding/json"
	"github.com/chen56/go-common/must"
	_ "github.com/go-sql-driver/mysql"
)

func MustMarshal2String(v interface{}) string {
	b, err := json.Marshal(v)
	must.NoError(err)
	return string(b)
}

func MustMarshal(v interface{}) []byte {
	b, err := json.Marshal(v)
	must.NoError(err)
	return b
}
