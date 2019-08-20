package ldapx

import (
	"github.com/chen56/go-common/must"
	ldap "gopkg.in/ldap.v3"
)

type LdapConf struct {
	Addr string `yaml:"addr"          json:"addr"`
}

type Conn struct {
	*ldap.Conn `yaml:"addr"          json:"addr"`
}

func (conf LdapConf) NewConnect() Conn {
	con, err := ldap.Dial("tcp", conf.Addr)
	must.NoError(err)

	return Conn{con}
}
