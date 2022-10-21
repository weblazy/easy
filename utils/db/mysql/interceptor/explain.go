package interceptor

import (
	"fmt"

	"gorm.io/gorm"
)

type Explain struct {
	Id           int64   `gorm:"column:id"`
	SelectType   string  `gorm:"column:select_type"`   // 查询行为类型 simple primary union...
	Table        string  `gorm:"column:table"`         // tableName
	Partitions   string  `gorm:"column:partitions"`    // 分区
	Type         string  `gorm:"column:type"`          // 引擎层查询数据行为类型 system const ref index index_merge all ...
	PossibleKeys string  `gorm:"column:possible_keys"` // 可能用到的所有索引
	Key          string  `gorm:"column:key"`           // 真正用到的所有索引
	KeyLen       int32   `gorm:"column:key_len"`       // 查询时用到的索引长度
	Ref          string  `gorm:"column:ref"`           // 哪些列或常量与key所使用的字段进行比较
	Rows         int32   `gorm:"column:rows"`          // 预估需要扫描的行数
	Filtered     float32 `gorm:"column:filtered"`      // 根据条件过滤后剩余的行数百分比（预估）
	Extra        string  `gorm:"column:extra"`
}

type ExplainPlugin struct{}

func NewExplainPlugin() *ExplainPlugin {
	return &ExplainPlugin{}
}

func (e *ExplainPlugin) Name() string {
	return "explain"
}

func (e *ExplainPlugin) Initialize(db *gorm.DB) error {
	return db.Callback().Query().After("gorm:query").Register("explain", checkIndex)
}

func checkIndex(db *gorm.DB) {
	result := &Explain{}
	session := &gorm.Session{
		NewDB:   true,
		Context: db.Statement.Context,
	}
	err := db.Session(session).Raw("EXPLAIN "+db.Statement.SQL.String(), db.Statement.Vars...).Scan(result).Error
	if err != nil {
		return
	}

	// 命中索引
	if result.Key != "" {
		fmt.Printf("hits index: %s", result.Key)
	}
}
