package configs

import (
	//"github.com/ghodss/yaml" //todo ghodss/yaml 这个问题？
	"gopkg.in/yaml.v2"
)

func marshal(value interface{}) (out []byte, err error) {
	return yaml.Marshal(value)
}
func unmarshal(in []byte, out interface{}) error {
	return yaml.Unmarshal(in, out)
}
