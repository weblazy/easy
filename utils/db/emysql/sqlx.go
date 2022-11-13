package emysql

import (
	"fmt"
	"reflect"
	"strings"

	"gorm.io/gorm"
)

var (
	FieldsError = fmt.Errorf("fileds length is 0")
)

// @desc
// @auth liuguoqiang 2020-11-27
// @param
// @return
func BulkInsert(db *gorm.DB, table string, fields []string, params []map[string]interface{}) error {
	if len(params) == 0 {
		return nil
	}
	if len(fields) == 0 {
		return FieldsError
	}
	sql := "INSERT INTO `" + table + "` (`" + strings.Join(fields, "`,`") + "`) VALUES "
	args := make([]interface{}, 0)
	valueArr := make([]string, 0)
	varArr := make([]string, 0)
	for _, obj := range params {
		varArr = varArr[:0]
		varStr := "("
		for _, value := range fields {
			if _, ok := obj[value]; !ok {
				return fmt.Errorf("%s:not found in fields", value)
			}
			varArr = append(varArr, "?")
			args = append(args, obj[value])
		}
		varStr += strings.Join(varArr, ",") + ")"
		valueArr = append(valueArr, varStr)
	}
	sql += strings.Join(valueArr, ",")
	err := db.Exec(sql, args...).Error
	return err
}

// @desc 批量插入
// @auth liuguoqiang 2020-11-27
// @param
// @return
func BulkSave(db *gorm.DB, table string, fields []string, params []map[string]interface{}) error {
	if len(params) == 0 {
		return nil
	}
	if len(fields) == 0 {
		return FieldsError
	}
	sql := "INSERT INTO `" + table + "` (`" + strings.Join(fields, "`,`") + "`) VALUES "
	updateArr := make([]string, 0)
	args := make([]interface{}, 0)
	valueArr := make([]string, 0)
	varArr := make([]string, 0)
	for _, value := range fields {
		updateArr = append(updateArr, value+"=VALUES("+value+")")
	}
	for _, obj := range params {
		varArr = varArr[:0]
		varStr := "("
		for _, value := range fields {
			if _, ok := obj[value]; !ok {
				return fmt.Errorf("%s字段在map中不存在", value)
			}
			varArr = append(varArr, "?")
			args = append(args, obj[value])
		}
		varStr += strings.Join(varArr, ",") + ")"
		valueArr = append(valueArr, varStr)
	}
	sql += strings.Join(valueArr, ",")
	sql += " ON DUPLICATE KEY UPDATE " + strings.Join(updateArr, ",")
	err := db.Exec(sql, args...).Error
	return err
}

// @desc
// @auth liuguoqiang 2020-04-08
// @param
// @return
func Validate(data, model interface{}) bool {
	if _, ok := data.(map[string]interface{}); ok {
		return true
	}
	if reflect.TypeOf(data).Kind() == reflect.TypeOf(model).Kind() {
		return true
	}
	return false
}
