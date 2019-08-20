package thriftx

import (
	"context"
	"github.com/apache/thrift/lib/go/thrift"
)

type PoolClientConf struct {
	HostPort string `yaml:"hostPort" json:"hostPort"`
}
type poolClient struct {
	conf             PoolClientConf
	transportFactory thrift.TTransportFactory
	protocolFactory  thrift.TProtocolFactory
}

// thrift.TStandardClient是非线程安全的，没有连接池机制
// 我们自己搞一个线程安全的，池化的client
func (x PoolClientConf) NewPoolClient() (client thrift.TClient) {
	var protocolFactory thrift.TProtocolFactory = thrift.NewTBinaryProtocolFactoryDefault()
	var transportFactory thrift.TTransportFactory = thrift.NewTBufferedTransportFactory(8192)
	transportFactory = thrift.NewTFramedTransportFactory(transportFactory)

	return &poolClient{
		protocolFactory:  protocolFactory,
		conf:             x,
		transportFactory: transportFactory,
	}
}

func (x *poolClient) Call(ctx context.Context, method string, args, result thrift.TStruct) error {
	// todo 暂时先每请求开一个client, 需要搞个连接池自动重连，简化操作
	var transport thrift.TTransport
	transport, err := thrift.NewTSocket(x.conf.HostPort)
	if err != nil {
		return err
	}

	transport, err = x.transportFactory.GetTransport(transport)
	if err != nil {
		return err
	}
	defer transport.Close()
	if err := transport.Open(); err != nil {
		return err
	}
	iprot := x.protocolFactory.GetProtocol(transport)
	oprot := x.protocolFactory.GetProtocol(transport)

	stdCient := thrift.NewTStandardClient(iprot, oprot)

	return stdCient.Call(ctx, method, args, result)
}
