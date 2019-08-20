package configs

import (
	"fmt"
	"io/ioutil"
)

//ConfigSource 配置源，配置可以从字符串/文件/网络读取
type ConfigSource interface {
	fmt.Stringer
	Read() ([]byte, error)
}

//NewFileConfigSource 文件配置源
func NewFileConfigSource(path string) ConfigSource {
	return &fileConfigSource{path: path}
}

type fileConfigSource struct {
	path string
}

func (t *fileConfigSource) Read() ([]byte, error) {
	return ioutil.ReadFile(t.path)
}

func (t *fileConfigSource) String() string {
	return "file:" + t.path
}

//NewMemoryConfigSource 内存配置源
func NewMemoryConfigSource(data []byte) ConfigSource {
	return &memoryConfigSource{data: data}
}

type memoryConfigSource struct {
	data []byte
}

func (t *memoryConfigSource) Read() ([]byte, error) {
	return t.data, nil
}
func (t *memoryConfigSource) String() string {
	return "memory"
}

//NewMemoryConfigSource 内存配置源
func NewInterfaceConfigSource(conf interface{}) ConfigSource {
	return &interfaceConfigSource{conf: conf}
}

type interfaceConfigSource struct {
	conf interface{}
}

func (t *interfaceConfigSource) Read() ([]byte, error) {
	return marshal(t.conf)
}

func (t *interfaceConfigSource) String() string {
	return fmt.Sprintf("interface: %T", t.conf)
}
