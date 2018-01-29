package es

import (
	"encoding/json"
	"fmt"

	"github.com/apex/log"
	"gopkg.in/olivere/elastic.v5"
)

//Deprecated
//废弃，应该在自己的项目中创建
func NewESClient(url string) (client *elastic.Client, err error) {
	return elastic.NewClient(
		elastic.SetURL(url),
		elastic.SetSniff(false),
		elastic.SetHealthcheck(false),
		elastic.SetErrorLog(&esLogAdapter{log.ErrorLevel}),
		elastic.SetTraceLog(&esLogAdapter{log.DebugLevel}),
		elastic.SetInfoLog(&esLogAdapter{log.InfoLevel}),
	)
}

type esLogAdapter struct {
	level log.Level
}

func (l *esLogAdapter) Printf(format string, v ...interface{}) {
	log.Debugf(format, v...)
	switch l.level {
	case log.DebugLevel:
		log.Debugf(format, v)
	case log.InfoLevel:
		log.Infof(format, v)
	case log.WarnLevel:
		log.Warnf(format, v)
	case log.ErrorLevel:
		log.Errorf(format, v)
	case log.FatalLevel:
		log.Fatalf(format, v)
	default:
		panic(fmt.Sprintf("not support log level :%s", l.level))
	}
}

func PrintSearchSource(source *elastic.SearchSource) {
	src, err := source.Source()
	if err != nil {
		panic(err)
	}
	data, err := json.Marshal(src)
	if err != nil {
		panic(err)
	}
	got := string(data)
	fmt.Println(got)
}
