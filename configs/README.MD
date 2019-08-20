# conifgs

基于配置的第三方辅助工具

yaml -> go struct -> 具体业务

配置通常和某环境(dev,prod,qa)相关，通过多个配置的merge,比如default + dev，构成某目标环境的具体配置

目前适用范围：

- mysql
- thrift
- redis
- grpc server | client
- kafka
- ldap
- fluent
- elastic
