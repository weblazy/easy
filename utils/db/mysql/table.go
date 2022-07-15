package mysql

import (
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

func GetSumInt64(sql string, args ...interface{}) (int64, error) {
	type sum struct {
		Num int64 `json:"num" gorm:"column:num"`
	}
	var obj sum
	return obj.Num, GetDB().Raw(sql, args...).Scan(&obj).Error
}
func GetSumFloat64(sql string, args ...interface{}) (float64, error) {
	type sum struct {
		Num float64 `json:"num" gorm:"column:num"`
	}
	var obj sum
	return obj.Num, GetDB().Raw(sql, args...).Scan(&obj).Error
}

func GetSumDecimal(sql string, args ...interface{}) (decimal.Decimal, error) {
	type sum struct {
		Num decimal.Decimal `json:"num" gorm:"column:num"`
	}
	var obj sum
	return obj.Num, GetDB().Raw(sql, args...).Scan(&obj).Error
}

func BatchInsert(db *gorm.DB, table string, fields []string, params []map[string]interface{}) error {
	if db == nil {
		db = GetDB()
	}
	return BulkInsert(db, table, fields, params)
}
func BatchSave(db *gorm.DB, table string, fields []string, params []map[string]interface{}) error {
	if db == nil {
		db = GetDB()
	}
	return BulkSave(db, table, fields, params)
}
