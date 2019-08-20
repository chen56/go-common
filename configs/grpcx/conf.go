package grpcx

import (
	"fmt"
)

// GrpcConf 配置
type GrpcConf struct {
	Port        int `yaml:"port"  json:"port"`
	GatewayPort int `yaml:"gatewayPort"  json:"gatewayPort"`
}

// GrpcListenAddress `:Port`
func (x GrpcConf) GrpcListenAddress() string {
	return fmt.Sprintf(":%d", x.Port)
}

// GrpcListenAddress `:Port`
func (x GrpcConf) GrpcEndpoint() string {
	return fmt.Sprintf("localhost:%d", x.Port)
}

// GatewayListenAddress `:Port`
func (x GrpcConf) GatewayListenAddress() string {
	return fmt.Sprintf(":%d", x.GatewayPort)
}
