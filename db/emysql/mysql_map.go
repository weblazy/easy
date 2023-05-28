package emysql

import (
	"sync"

	"github.com/weblazy/easy/db/emysql/emysql_config"
	"github.com/weblazy/easy/econfig"
)

var MysqlMap sync.Map

// GetMysql return a MysqlClient
func GetMysql(dbName string) *MysqlClient {
	if v, ok := MysqlMap.Load(dbName); ok {
		return v.(*MysqlClient)
	}
	conf := emysql_config.DefaultConfig()
	econfig.GlobalViper.UnmarshalKey(dbName, conf)
	mysqlClient, err := NewMysqlClient(conf)
	if err != nil {
		return nil
	}
	MysqlMap.Store(dbName, mysqlClient)
	return mysqlClient
}
