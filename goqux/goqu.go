package goqux

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/doug-martin/goqu/v8"
	"github.com/sirupsen/logrus"
	//非常重要的声明，要不然goqu 无法正确工作
	_ "github.com/doug-martin/goqu/v8/dialect/mysql"
)

// DB 数据库扩展
type Database struct {
	db   *sql.DB
	Goqu *goqu.Database
}

// Tx 事务扩展
type Tx struct {
	Tx     *sql.Tx
	GoquTx *goqu.TxDatabase
}

// NewDb 新建
func NewDb(db *sql.DB) *Database {
	var goquDB = goqu.New("mysql", db)
	goquDB.Logger(&GoquLoggerAdapter{})
	return &Database{
		db:   db,
		Goqu: goquDB,
	}
}

type GoquLoggerAdapter struct {
}

func (x *GoquLoggerAdapter) Printf(format string, v ...interface{}) {
	fmt.Println("=========================")
	logrus.Infof(format, v...)
}

// InTxContext 运行在一个事务上下文
func (db *Database) InTxContext(ctx context.Context, opts *sql.TxOptions, txFunc func(tx *Tx) error) (err error) {
	tx, err := db.db.BeginTx(ctx, opts)
	if err != nil {
		return
	}

	return inTx(&Tx{tx, goqu.NewTx("mysql", tx)}, txFunc)
}

//InTx ref: https://stackoverflow.com/questions/16184238/database-sql-tx-detecting-commit-or-rollback
func (db *Database) InTx(txFunc func(tx *Tx) error) (err error) {
	tx, err := db.db.Begin()
	if err != nil {
		return
	}
	return inTx(&Tx{tx, goqu.NewTx("mysql", tx)}, txFunc)
}

func inTx(tx *Tx, txFunc func(tx *Tx) error) (err error) {
	defer func() {
		if r := recover(); r != nil {
			//已在一个panic中，不理会Rollback的err
			err := tx.Tx.Rollback()
			logrus.Error(err)
			// re-throw panic after Rollback
			panic(r)
		}
		if err != nil {
			//err已非空，不理会Rollback的err
			err := tx.Tx.Rollback()
			logrus.Error(err)
			return
		}
		// err==nil, commit; 如果commit失败，则返回err
		err = tx.Tx.Commit()
	}()

	err = txFunc(tx)

	return err
}
