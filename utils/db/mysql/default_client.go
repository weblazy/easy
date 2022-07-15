package mysql

import (
	"gorm.io/gorm"
)

func init() {
	MysqlClient = &DefaultClient{}
	GetDB = func() *gorm.DB {
		db := GetORM("conf.DBFreezonecoin")
		return db
	}
}

type DefaultClient struct {
	Client
}

func (*DefaultClient) GetDB() *gorm.DB {
	db := GetORM("conf.DBFreezonecoin")
	return db
}
