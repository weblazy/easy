package emysql

import (
	"context"
	"sync"

	"github.com/weblazy/easy/db/emysql/emysql_config"
	"github.com/weblazy/easy/econfig"
	"gorm.io/gorm"
)

var MysqlMap sync.Map

// GetMysql return a MysqlClient
func GetMysql(ctx context.Context, dbName string) *gorm.DB {
	// 判断是否开启了事务，如果存在dbname开启事务直接返回/如果没有开启并启动事务
	tx := checkTransaction(ctx, dbName)
	if tx != nil {
		return tx
	}
	return getMysql(ctx, dbName)
}

// getMysql return a *gorm.DB
func getMysql(ctx context.Context, dbName string) *gorm.DB {
	if v, ok := MysqlMap.Load(dbName); ok {
		return v.(*MysqlClient).DB.WithContext(ctx)
	}
	conf := emysql_config.DefaultConfig()
	econfig.GlobalViper.UnmarshalKey(dbName, conf)
	mysqlClient, err := NewMysqlClient(conf)
	if err != nil {
		return nil
	}

	MysqlMap.Store(dbName, mysqlClient)
	return mysqlClient.DB.WithContext(ctx)
}
