package mysql

import (
	"gorm.io/gorm"
)

func init() {
	mysqlClient = &DefaultClient{}
}

type DefaultClient struct {
	Client
}

func (*DefaultClient) GetDB() *gorm.DB {
	db := GetORM("conf.DBFreezonecoin")
	return db
}
