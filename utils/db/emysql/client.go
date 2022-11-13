package emysql

import "gorm.io/gorm"

// var MysqlClient Client

type GetMysqlDB func(key string) *gorm.DB

var GetDB GetMysqlDB

type Client interface {
	GetDB(key string) *gorm.DB
}
