package kvdbmysql

import (
	"database/sql"

	"fmt"

	"io"

	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/xiaonanln/goworld/engine/gwlog"
	"github.com/xiaonanln/goworld/engine/kvdb/types"
)

const (
	_MAX_KEY_LENGTH = 256
)

type mysqlKVDB struct {
	driverName     string
	dataSourceName string
	db             *sql.DB
}

// OpenMySQLKVDB opens SQL driver for KVDB backend
func OpenMySQLKVDB(dataSourceName string) (kvdbtypes.KVDBEngine, error) {
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	// try to create the __kv__ table if not exists
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS `__kv__`(`key` VARCHAR(" + strconv.Itoa(_MAX_KEY_LENGTH) + ") NOT NULL PRIMARY KEY, `val` BLOB NOT NULL)")
	if err != nil {
		return nil, err
	}

	return &mysqlKVDB{
		driverName:     "mysql",
		dataSourceName: dataSourceName,
		db:             db,
	}, nil
}

func (sqlkvdb *mysqlKVDB) String() string {
	return fmt.Sprintf("%s<%s>", sqlkvdb.driverName, sqlkvdb.dataSourceName)
}

func (sqlkvdb *mysqlKVDB) Get(key string) (val string, err error) {
	row := sqlkvdb.db.QueryRow("SELECT `val` FROM `__kv__` WHERE `key` = ?", key)
	err = row.Scan(&val)
	if err == sql.ErrNoRows {
		err = nil // not found, use default val ""
	}

	return
}

func (sqlkvdb *mysqlKVDB) Put(key string, val string) (err error) {
	_, err = sqlkvdb.db.Exec("INSERT INTO `__kv__`(`key`, `val`) VALUES(?, ?) ON DUPLICATE KEY UPDATE `val`=?", key, val, val)
	return
}

type sqlKVDBIterator struct {
	rows *sql.Rows
}

func (it *sqlKVDBIterator) Next() (kvdbtypes.KVItem, error) {
	if it.rows.Next() {
		var item kvdbtypes.KVItem
		err := it.rows.Scan(&item.Key, &item.Val)
		return item, err
	} else {
		return kvdbtypes.KVItem{}, io.EOF
	}
}

func (sqlkvdb *mysqlKVDB) Find(beginKey string, endKey string) (kvdbtypes.Iterator, error) {
	rows, err := sqlkvdb.db.Query("SELECT `key`, `val` FROM `__kv__` WHERE `key` >= ? AND `key` < ?", beginKey, endKey)
	if err != nil {
		return nil, err
	}

	return &sqlKVDBIterator{
		rows: rows,
	}, nil
}

func (sqlkvdb *mysqlKVDB) Close() {
	if err := sqlkvdb.db.Close(); err != nil {
		gwlog.Errorf("%s: close error: %s", sqlkvdb.String(), err)
	}
}
func (sqlkvdb *mysqlKVDB) IsConnectionError(err error) bool {
	return true
}
