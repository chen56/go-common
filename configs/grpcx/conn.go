package grpcx

import (
	"google.golang.org/grpc"
)

type ConnConf struct {
	HostPort string `yaml:"hostPort" json:"hostPort"`
}

func NewConnConf() *ConnConf {
	return &ConnConf{
		HostPort: "localhost:2080",
	}
}

func (x ConnConf) MustConnect() (conn *grpc.ClientConn) {
	conn, err := grpc.Dial(x.HostPort, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	return
}
