package mysql

import "gorm.io/gorm"

var mysqlClient Client

type Client interface {
	GetDB() *gorm.DB
}

func SetMysqlClient(c Client) {
	mysqlClient = c
}

func GetMysqlClient() Client {
	return mysqlClient
}
