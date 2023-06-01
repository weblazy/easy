package emysql

import (
	"database/sql/driver"
	"encoding/json"
)

type Map map[string]interface{}

// Value 实现方法
func (m Map) Value() (driver.Value, error) {
	return json.Marshal(m)
}

// Scan 实现方法
func (m *Map) Scan(input interface{}) error {
	return json.Unmarshal(input.([]byte), m)
}
