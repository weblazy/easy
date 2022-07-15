package mysql

import "gorm.io/gorm"

var MysqlClient Client

type GetMysqlDB func() *gorm.DB

var GetDB GetMysqlDB

type Client interface {
	GetDB() *gorm.DB
}
