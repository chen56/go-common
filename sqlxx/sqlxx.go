package sqlxx

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

// DB 数据库扩展
type DB struct {
	*sqlx.DB
}

// Tx 事务扩展
type Tx struct {
	*sqlx.Tx
}

// NewDb 新建
func NewDb(db *sql.DB, driverName string) *DB {
	return &DB{
		DB: sqlx.NewDb(db, driverName),
	}
}

// InTxContext 运行在一个事务上下文
func (db *DB) InTxContext(ctx context.Context, opts *sql.TxOptions, txFunc func(tx *Tx) error) (err error) {
	tx, err := db.DB.BeginTxx(ctx, opts)
	if err != nil {
		return
	}
	return inTx(&Tx{Tx: tx}, txFunc)
}

//InTx ref: https://stackoverflow.com/questions/16184238/database-sql-tx-detecting-commit-or-rollback
func (db *DB) InTx(txFunc func(tx *Tx) error) (err error) {
	tx, err := db.DB.Beginx()
	if err != nil {
		return
	}
	return inTx(&Tx{Tx: tx}, txFunc)
}

func inTx(sqlxTx *Tx, txFunc func(tx *Tx) error) (err error) {
	defer func() {
		if r := recover(); r != nil {
			//已在一个panic中，不理会Rollback的err
			err := sqlxTx.Tx.Rollback()
			logrus.Error(err)
			// re-throw panic after Rollback
			panic(r)
		}
		if err != nil {
			//err已非空，不理会Rollback的err
			err := sqlxTx.Tx.Rollback()
			logrus.Error(err)
			return
		}
		// err==nil, commit; 如果commit失败，则返回err
		err = sqlxTx.Tx.Commit()
	}()

	err = txFunc(sqlxTx)

	return err
}

func ToJson(rows *sql.Rows) (string, error) {

	columnTypes, err := rows.ColumnTypes()
	if err != nil {
		return "", err
	}

	columns, err := rows.Columns()
	if err != nil {
		return "", err
	}
	count := len(columns)
	tableData := make([]map[string]interface{}, 0)
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)
	for rows.Next() {
		for i := 0; i < count; i++ {
			valuePtrs[i] = &values[i]
		}
		err := rows.Scan(valuePtrs...)
		if err != nil {
			return "", err
		}
		entry := make(map[string]interface{})
		for i, col := range columns {
			var v interface{}
			val := values[i]
			b, ok := val.([]byte)
			if ok {
				v = string(b)
				fmt.Println(columnTypes[i])
			} else {
				v = val
			}
			entry[col] = v

		}
		tableData = append(tableData, entry)
	}
	jsonData, err := json.Marshal(tableData)
	if err != nil {
		return "", err
	}
	fmt.Println(string(jsonData))
	return string(jsonData), nil
}
