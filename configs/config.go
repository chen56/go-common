package configs

import (
	"fmt"
	"github.com/chen56/go-common/must"
	"github.com/chen56/go-common/reflectx"

	"github.com/pkg/errors"
	"os"
	"reflect"
	"strings"
)

//NewConfigOptions 选项
type NewConfigOptions struct {
	ConfigSources []ConfigSource
}

//Config raw config
type Config struct {
	config map[string]interface{}
}

//NewConfig 新建
func NewConfig(options NewConfigOptions) Config {
	result := Config{}
	must.True(len(options.ConfigSources) >= 1, "最少一个ConfigSource, but: %v", options.ConfigSources)

	var configs []map[string]interface{}
	for _, source := range options.ConfigSources {
		var config map[interface{}]interface{}
		bytes, err := source.Read()
		must.NoError(err)
		err = unmarshal(bytes, &config)
		must.NoError(err)
		configs = append(configs, cloneInterfaceMap(config))
	}

	//第一个认为是主配置，依次用附加配置覆盖之
	var config = configs[0]

	// 把附加的config marge到第一个主配置中
	{
		var others = configs[1:]
		for _, other := range others {
			margeMap(config, other, []string{})
		}
	}

	// 把环境变量marge到config中
	{
		env := map[string]string{}
		collectEnvToKV(env)
		margeEnvValue(config, env, []string{})
	}

	// config + env = 混合kv
	{
		configAndEnv := map[string]string{}
		collectConfigTreeToKV(config, configAndEnv, []string{})
		collectEnvToKV(configAndEnv)

		// 用混合后的kv对config中的所有字符串值做变量替换
		expandEnv(config, configAndEnv)
	}

	result.config = config

	return result
}

// ScanConfig 求出强类型配置
// todo 用 github.com/ghodss/yaml 来避免json/yaml转换问题
func (x *Config) ScanConfig(conf interface{}) error {
	val := reflect.ValueOf(conf)
	if val.Kind() != reflect.Ptr {
		return errors.Errorf("conf must be a pointer to a struct or map[string]interface{} when calling ScanConfig, but: (%s) ", reflect.TypeOf(conf))
	}
	if reflect.Indirect(val).Kind() == reflect.Map {
		if _, ok := conf.(*map[string]interface{}); !ok {
			return errors.Errorf("conf must be a pointer to a struct or map[string]interface{} when calling ScanConfig, but: (%s) ", reflect.TypeOf(conf))
		}
	}
	margedYaml, err := marshal(x.config)
	if err != nil {
		return err
	}
	err = unmarshal(margedYaml, conf)
	if err != nil {
		return err
	}

	if asMap, ok := conf.(*map[string]interface{}); ok {
		val = reflect.Indirect(reflect.ValueOf(conf))
		cloned := cloneStringMap(*asMap)
		val.Set(reflect.ValueOf(cloned))
		return nil
	}
	return nil
}

// MustScanConfig 求出强类型配置
func (x *Config) MustScanConfig(i interface{}) {
	err := x.ScanConfig(i)
	must.NoError(err)
}

func cloneInterfaceMap(config map[interface{}]interface{}) map[string]interface{} {
	var result = map[string]interface{}{}
	for key, value := range config {
		//若值为map,则递归往复
		if valueAsMap, ok := value.(map[interface{}]interface{}); ok {
			result[fmt.Sprint(key)] = cloneInterfaceMap(valueAsMap)
		} else {
			result[fmt.Sprint(key)] = value
		}
	}
	return result
}
func cloneStringMap(config map[string]interface{}) map[string]interface{} {
	for key, value := range config {
		if valueAsMap, ok := value.(map[interface{}]interface{}); ok {
			config[key] = cloneInterfaceMap(valueAsMap)
		}
	}
	return config
}

func collectEnvToKV(env map[string]string) {
	for _, kv := range os.Environ() {
		for i := 0; i < len(kv); i++ {
			if kv[i] == '=' {
				k := kv[:i]
				v := kv[i+1:]
				//key 全转换为大写，方便后面case不敏感获取
				env[strings.ToUpper(k)] = v
			}
		}
	}
}

func collectConfigTreeToKV(config map[string]interface{}, keyValues map[string]string, parentPath []string) {
	for key, value := range config {

		childPath := append(parentPath, strings.ToUpper(fmt.Sprint(key)))

		//若值为map,则递归往复
		if valueAsMap, ok := value.(map[string]interface{}); ok {
			collectConfigTreeToKV(valueAsMap, keyValues, childPath)
			continue
		}

		reflectValue := reflect.TypeOf(value)
		if reflectValue == nil {
			continue
		}
		if !reflectx.IsPrimitive(reflectValue.Kind()) {
			continue
		}
		envKey := pathToEnvKey(childPath)
		keyValues[envKey] = fmt.Sprint(value)
	}
}

// 用`_`连接树节点路径，转为环境变量的key,
func pathToEnvKey(path []string) string {
	return strings.ToUpper(strings.Join(path, "_"))
}

// 把右边的map合并到左边
// 左边：原值，右边：覆盖值
func margeMap(left map[string]interface{}, right map[string]interface{}, parentPath []string) {
	for key, rightValue := range right {
		childPath := append(parentPath, fmt.Sprint(key))

		//如果不存在，就补上
		leftValue, ok := left[key]
		if !ok {
			left[key] = rightValue
			continue
		}

		//如果右值(覆盖值)为map,就递归往复
		if rightValueAsMap, ok := rightValue.(map[string]interface{}); ok {
			//左value(原值)不是map,则应是配置bug
			leftValueAsMap, ok := leftValue.(map[string]interface{})
			if !ok {
				panic(fmt.Sprintf("配置[%s]位类型不符，left:%v, override:%v", strings.Join(childPath, "."), leftValue, rightValue))
			}
			margeMap(leftValueAsMap, rightValueAsMap, childPath)
			continue
		}

		left[key] = rightValue
	}
}

// 合并环境变量到config中
func margeEnvValue(config map[string]interface{}, env map[string]string, parentPath []string) {
	for key, value := range config {
		childPath := append(parentPath, strings.ToUpper(fmt.Sprint(key)))
		//若值为map,则递归往复
		if valueAsMap, ok := value.(map[string]interface{}); ok {
			margeEnvValue(valueAsMap, env, childPath)
			continue
		}
		envKey := pathToEnvKey(childPath)
		envValue, ok := env[envKey]
		if ok {
			config[key] = envValue
		}
	}
}

// 变量替换,不区分大小写.
// 例，原yaml:
//   `appInfo: ${HOME},${Home}`
// expandEnv后:
//   `appInfo: /Users/chenpeng,/Users/chenpeng`
func expandEnv(config map[string]interface{}, env map[string]string) {
	err := visit(config, []string{}, func(parentNode map[string]interface{}, nodeKey string, nodePath []string, node interface{}) error {
		switch tnode := node.(type) {
		case string:
			newValue := os.Expand(tnode, func(key string) string {
				return env[strings.ToUpper(key)]
			})
			parentNode[nodeKey] = newValue
		}
		return nil
	})
	must.NoError(err, "不应该发生")
}

//Vistor Visit的访问器
type Vistor func(parentNode map[string]interface{}, nodeKey string, nodePath []string, node interface{}) error

//Visit 遍历
func (x *Config) Visit(visitor Vistor) error {
	return visit(x.config, []string{}, visitor)
}

func (x *Config) ToYaml() (out []byte, err error) {
	return marshal(x.config)
}
func (x *Config) MustToYaml() (out []byte) {
	out, err := marshal(x.config)
	must.NoError(err)
	return out
}

func visit(parentNode map[string]interface{}, parentPath []string, visitor Vistor) error {
	for key, node := range parentNode {
		path := append(parentPath, fmt.Sprint(key))
		err := visitor(parentNode, key, path, node)
		if err != nil {
			return err
		}
		//若值为map,则递归往复
		if mapNode, ok := node.(map[string]interface{}); ok {
			err := visit(mapNode, path, visitor)
			if err != nil {
				return err
			}
			continue
		}
	}
	return nil
}
