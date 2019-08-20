package elastic5x

import (
	"github.com/chen56/go-common/must"
	"github.com/sirupsen/logrus"
	"gopkg.in/olivere/elastic.v5"
)

//Conf https://github.com/olivere/elastic/wiki/Sniffing
type Conf struct {
	URL         string `yaml:"url"          json:"url"`
	Sniff       bool   `yaml:"sniff"        json:"sniff"`
	Healthcheck bool   `yaml:"healthcheck"  json:"healthcheck"`

	Errorlog elastic.Logger // error log for critical messages
	Infolog  elastic.Logger // information log for e.g. response times
	Tracelog elastic.Logger // trace log for debugging
}

//Client elastic.Client 扩展
type Client struct {
	*elastic.Client
}

//MustNewClient 新建Client
func (c Conf) MustNewClient() *Client {
	var errorlog = c.Errorlog
	if errorlog == nil {
		errorlog = &esLogAdapter{logrus.ErrorLevel}
	}
	var infolog = c.Infolog
	if infolog == nil {
		infolog = &esLogAdapter{logrus.InfoLevel}
	}
	var tracelog = c.Tracelog
	if tracelog == nil {
		tracelog = &esLogAdapter{logrus.TraceLevel}
	}
	client, err := elastic.NewClient(
		elastic.SetSniff(c.Sniff),
		elastic.SetURL(c.URL),
		elastic.SetHealthcheck(c.Healthcheck),

		//定制日志转发
		elastic.SetErrorLog(errorlog),
		elastic.SetInfoLog(infolog),
		elastic.SetTraceLog(tracelog),
	)

	must.NoError(err)
	return &Client{
		client,
	}
}

type esLogAdapter struct {
	level logrus.Level
}

func (l *esLogAdapter) Printf(format string, v ...interface{}) {
	logrus.StandardLogger().Logf(l.level, format, v...)
}
