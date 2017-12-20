package jsonx

import (
	_ "github.com/go-sql-driver/mysql"
	"encoding/json"
	"github.com/chen56/go-common/assert"
)

func MustMarshal2String(v interface{})string{
	b,err:=json.Marshal(v)
	assert.NoErr(err)
	return  string(b)
}

func MustMarshal(v interface{})[]byte{
	b,err:=json.Marshal(v)
	assert.NoErr(err)
	return b
}
