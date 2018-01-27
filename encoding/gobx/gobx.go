package gobx

import (
	"bytes"
	"encoding/gob"
)

func Encode(x interface{}) ([]byte, error) {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	err := encoder.Encode(x)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func Decode(b []byte, result interface{}) error {
	reader := bytes.NewBuffer(b)
	decoder := gob.NewDecoder(reader)
	return decoder.Decode(result)
}
