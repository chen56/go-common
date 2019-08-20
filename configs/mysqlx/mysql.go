package mysqlx

import (
	"code.cloudfoundry.org/bytefmt"
	"database/sql"
	"github.com/chen56/go-common/must"
	"github.com/chen56/go-common/timex"
	"github.com/go-sql-driver/mysql"
	"time"
)

/*
MysqlConf mysql配置
*/
type MysqlConf struct {
	User   string `yaml:"user" json:"user"`
	Passwd string `yaml:"passwd" json:"passwd"`
	Addr   string `yaml:"addr" json:"addr"`
	DBName string `yaml:"dbName" json:"dbName"`

	//MaxAllowedPacket string `yaml:"maxAllowedPacket"`
	//ParseTime        bool   `yaml:"parseTime"`
	//Timeout          string `yaml:"timeout"`
	//Collation        string `yaml:"collation"`

	// sql.DB#SetMaxOpenConns()
	MaxOpenConns int `yaml:"maxOpenConns" json:"maxOpenConns"`
	// sql.DB#SetMaxIdleConns()
	MaxIdleConns int `yaml:"maxIdleConns" json:"maxIdleConns"`
	// sql.DB#SetConnMaxLifetime() 填写
	// 填写规则同time.Duration：1s , 1m , 1h
	ConnMaxLifetime time.Duration `yaml:"connMaxLifetime" json:"connMaxLifetime"`
}

// ToDSN 转换为dsn
func (c MysqlConf) ToDSN() string {
	mustToBytes := func(s string) uint64 {
		maxAllowedPacket, err := bytefmt.ToBytes(s) // 4 MiB
		must.NoError(err)
		return maxAllowedPacket
	}

	result := mysql.Config{
		Addr:   c.Addr,
		User:   c.User,
		Passwd: c.Passwd,
		DBName: c.DBName,

		Net:                  "tcp",
		Collation:            "utf8mb4_unicode_ci",
		Loc:                  timex.LocationAsiaShanghai,
		MaxAllowedPacket:     int(mustToBytes("4M")),
		AllowNativePasswords: true,
		ParseTime:            true,
		Timeout:              30 * time.Second,
		ReadTimeout:          30 * time.Second,
		WriteTimeout:         5 * time.Second,
	}
	return result.FormatDSN()
}

// MustConnect 连接
// 阿里云实例最大连接数，一般4000左右
func (c MysqlConf) MustConnect() *sql.DB {
	db, err := sql.Open("mysql", c.ToDSN())
	must.NoError(err, "dsn: %s", c.ToDSN())
	db.SetMaxOpenConns(c.MaxOpenConns)
	db.SetMaxIdleConns(c.MaxIdleConns)
	db.SetConnMaxLifetime(c.ConnMaxLifetime)
	return db
}
