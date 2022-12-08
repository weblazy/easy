package emysql

import (
	"gorm.io/gorm"
)

func init() {
	// MysqlClient = &DefaultClient{}
	// GetDB = func(key string) *gorm.DB {
	// 	db := GetORM(key)
	// 	return db
	// }
}

type DefaultClient struct {
	Client
}

func (*DefaultClient) GetDB(key string) *gorm.DB {
	db := GetORM(key)
	return db
}
