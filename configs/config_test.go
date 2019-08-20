package configs_test

import (
	"fmt"
	"github.com/chen56/go-common/configs"
	"github.com/chen56/go-common/must"
	"github.com/c2fo/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
	"time"
)

type Log struct {
	Level string `yaml:"level" json:"level"`
}
type Conf struct {
	Log     Log    `yaml:"log"  json:"log"`
	AppInfo string `yaml:"appInfo"  json:"appInfo"`
}

var testdatas = []struct {
	test   string
	conf   []string
	env    map[string]string
	result Conf
}{
	{
		"读取1个配置",
		[]string{`
log:
  level: error
`},
		nil,
		Conf{Log: Log{
			Level: "error",
		}},
	},

	{
		"多配置覆盖",
		[]string{`
log:
  level: error`,
			`
log:
  level: debug`},
		nil,
		Conf{Log: Log{
			Level: "debug",
		}},
	},

	{
		"环境变量覆盖,不区分大小写",
		[]string{`
log:
  level: error
`},
		map[string]string{"log_level": "warn"},
		Conf{Log: Log{
			Level: "warn",
		}},
	},

	{
		"环境变量覆盖,不区分大小写",
		[]string{`
log:
  level: error
appInfo: log=${log_level}
`},
		nil,
		Conf{Log: Log{
			Level: "error",
		}, AppInfo: "log=error"},
	},
}

func TestAllConfig(t *testing.T) {
	for _, d := range testdatas {
		var sources []configs.ConfigSource
		for _, conf := range d.conf {
			sources = append(sources, configs.NewMemoryConfigSource([]byte(conf)))
		}

		for k, v := range d.env {
			os.Setenv(k, v)
		}

		c := configs.NewConfig(configs.NewConfigOptions{ConfigSources: sources})

		actual := Conf{}
		require.NoError(t, c.ScanConfig(&actual))
		//spew.Dump(actual)
		require.Equal(t, d.result, actual)

		for k := range d.env {
			os.Unsetenv(k)
		}

	}
}

func ExampleConfig_ScanConfig() {
	var config1 = configs.NewMemoryConfigSource([]byte(`
log:
  level: error
appInfo: log=${log_level}
`))
	var config2 = configs.NewMemoryConfigSource([]byte(`
log:
  level: debug
`))

	c := configs.NewConfig(configs.NewConfigOptions{ConfigSources: []configs.ConfigSource{config1, config2}})

	conf := Conf{}
	err := c.ScanConfig(&conf)
	must.NoError(err)
	//spew.Dump(c)
	fmt.Printf("%+v", conf)

	// Output: {Log:{Level:debug} AppInfo:log=debug}
}

// 指针值可以表现null
func TestPointerValue(t *testing.T) {
	type TestConfig struct {
		Debug       bool  `yaml:"debug" json:"debug"`
		DebugPtr    *bool `yaml:"debugPtr" json:"debugPtr"`
		DebugPtrNil *bool `yaml:"debugPtrNil" json:"debugPtrNil"`
	}
	var config = []byte(`
debug: true
debugPtr: true
`)

	c := configs.NewConfig(configs.NewConfigOptions{
		ConfigSources: []configs.ConfigSource{configs.NewMemoryConfigSource(config)}},
	)
	conf := TestConfig{}
	err := c.ScanConfig(&conf)
	must.NoError(err)
	//spew.Dump(c)
	assert.NotNil(t, conf.DebugPtr)
	assert.Equal(t, true, *conf.DebugPtr)
	assert.Nil(t, conf.DebugPtrNil)
}
func TestDuration(t *testing.T) {
	type DurationConf struct {
		Second time.Duration `yaml:"second"  json:"second"`
		Minute time.Duration `yaml:"minute"  json:"minute"`
		Hour   time.Duration `yaml:"hour"  json:"hour"`
		Zero   time.Duration `yaml:"zero"  json:"zero"`
	}

	var tests = []struct {
		test   string
		conf   string
		result DurationConf
	}{
		{
			"读取1个配置",
			`
second: 1s
minute: 2m1s
hour: 3h2m1s
`,
			DurationConf{
				Second: 1 * time.Second,
				Minute: 2*time.Minute + 1*time.Second,
				Hour:   3*time.Hour + 2*time.Minute + 1*time.Second,
				Zero:   0,
			},
		},
	}
	for _, d := range tests {
		var sources = []configs.ConfigSource{configs.NewMemoryConfigSource([]byte(d.conf))}
		c := configs.NewConfig(configs.NewConfigOptions{ConfigSources: sources})

		actual := DurationConf{}
		require.NoError(t, c.ScanConfig(&actual))
		//spew.Dump(actual)
		require.Equal(t, d.result, actual)
	}
}
