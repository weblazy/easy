package sqlx

import (
	"fmt"
	"reflect"
	"strings"

	gorm "github.com/jinzhu/gorm"
)

// sql = INSERT INTO `users` VALUES (?,?,?),(?,?,?)
func BulkInsert(db *gorm.DB, table string, fields []string, params []map[string]interface{}) error {
	sql := "INSERT INTO `" + table + "` (" + strings.Join(fields, ",") + ") VALUES "
	args := make([]interface{}, 0)
	valueArr := make([]string, 0)
	varArr := make([]string, 0)
	for _, obj := range params {
		varArr = varArr[:0]
		varStr := "("
		for _, value := range fields {
			varArr = append(varArr, "?")
			args = append(args, obj[value])
		}
		varStr += strings.Join(varArr, ",") + ")"
		valueArr = append(valueArr, varStr)
	}
	sql += strings.Join(valueArr, ",")
	fmt.Println(sql)
	fmt.Println(args)
	err := db.Exec(sql, args...).Error
	return err
}

func BulkSave(db *gorm.DB, table string, fields []string, params []map[string]interface{}) error {
	sql := "INSERT INTO `" + table + "` (" + strings.Join(fields, ",") + ") VALUES "
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
			varArr = append(varArr, "?")
			args = append(args, obj[value])
		}
		varStr += strings.Join(varArr, ",") + ")"
		valueArr = append(valueArr, varStr)
	}
	sql += strings.Join(valueArr, ",")
	sql += " ON DUPLICATE KEY UPDATE " + strings.Join(updateArr, ",")
	fmt.Println(sql)
	fmt.Println(args)
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
