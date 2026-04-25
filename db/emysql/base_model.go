package emysql

import (
	"context"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// Tabler 接口定义 Model 必须实现的方法
type Tabler interface {
	TableName() string
	GetDB(ctx context.Context) *gorm.DB
}

// BaseModel 提供通用 CRUD 方法的泛型基类
// T 必须实现 Tabler 接口
type BaseModel[T Tabler] struct{}

// getDB 获取数据库连接（内部辅助方法）
func (b *BaseModel[T]) getDB(ctx context.Context) *gorm.DB {
	var model T
	return model.GetDB(ctx)
}

// Insert 插入单条记录
func (b *BaseModel[T]) Insert(ctx context.Context, data *T) error {
	return (*data).GetDB(ctx).Create(data).Error
}

// BulkInsert 批量插入
func (b *BaseModel[T]) BulkInsert(ctx context.Context, fields []string, params []map[string]any) error {
	var model T
	return BulkInsert(b.getDB(ctx), model.TableName(), fields, params)
}

// BulkSave 批量保存（upsert）
func (b *BaseModel[T]) BulkSave(ctx context.Context, fields []string, params []map[string]any) error {
	var model T
	return BulkSave(b.getDB(ctx), model.TableName(), fields, params)
}

// Delete 条件删除
func (b *BaseModel[T]) Delete(ctx context.Context, where string, args ...any) error {
	var model T
	return b.getDB(ctx).Where(where, args...).Delete(&model).Error
}

// Updates 条件更新，返回影响行数
func (b *BaseModel[T]) Updates(ctx context.Context, data map[string]any, where string, args ...any) (int64, error) {
	var model T
	db := b.getDB(ctx).Model(&model).Where(where, args...).Updates(data)
	return db.RowsAffected, db.Error
}

// GetOne 查询单条记录
func (b *BaseModel[T]) GetOne(ctx context.Context, where string, args ...any) (*T, error) {
	var obj T
	return &obj, b.getDB(ctx).Where(where, args...).Take(&obj).Error
}

// GetFirst 排序后查询第一条
func (b *BaseModel[T]) GetFirst(ctx context.Context, order string, where string, args ...any) (*T, error) {
	var obj T
	return &obj, b.getDB(ctx).Where(where, args...).Order(order).Take(&obj).Error
}

// GetList 查询列表
func (b *BaseModel[T]) GetList(ctx context.Context, where string, args ...any) ([]*T, error) {
	var list []*T
	return list, b.getDB(ctx).Where(where, args...).Find(&list).Error
}

// GetListWithLimit 带限制的查询
func (b *BaseModel[T]) GetListWithLimit(ctx context.Context, limit int, where string, args ...any) ([]*T, error) {
	var list []*T
	return list, b.getDB(ctx).Where(where, args...).Limit(limit).Find(&list).Error
}

// GetListOrder 排序查询
func (b *BaseModel[T]) GetListOrder(ctx context.Context, order string, where string, args ...any) ([]*T, error) {
	var list []*T
	return list, b.getDB(ctx).Where(where, args...).Order(order).Find(&list).Error
}

// GetListOrderLimit 排序+限制查询（带保护）
func (b *BaseModel[T]) GetListOrderLimit(ctx context.Context, order string, limit int, where string, args ...any) ([]*T, error) {
	var list []*T
	// 保护：防止查询过多数据
	if limit == 0 || limit > 10000 {
		limit = 10
	}
	return list, b.getDB(ctx).Where(where, args...).Order(order).Limit(limit).Find(&list).Error
}

// GetListPage 分页查询
func (b *BaseModel[T]) GetListPage(ctx context.Context, pageNum, limit int, where string, args ...any) ([]*T, error) {
	var list []*T
	offset := (pageNum - 1) * limit
	return list, b.getDB(ctx).Where(where, args...).Offset(offset).Limit(limit).Find(&list).Error
}

// GetCount 计数
func (b *BaseModel[T]) GetCount(ctx context.Context, where string, args ...any) (int64, error) {
	var model T
	var count int64
	return count, b.getDB(ctx).Model(&model).Where(where, args...).Count(&count).Error
}

// GetSumInt64 求和（int64）
func (b *BaseModel[T]) GetSumInt64(ctx context.Context, sql string, args ...any) (int64, error) {
	type sum struct {
		Num int64 `json:"num" gorm:"column:num"`
	}
	var obj sum
	return obj.Num, b.getDB(ctx).Raw(sql, args...).Scan(&obj).Error
}

// GetSumFloat64 求和（float64）
func (b *BaseModel[T]) GetSumFloat64(ctx context.Context, sql string, args ...any) (float64, error) {
	type sum struct {
		Num float64 `json:"num" gorm:"column:num"`
	}
	var obj sum
	return obj.Num, b.getDB(ctx).Raw(sql, args...).Scan(&obj).Error
}

// GetSumDecimal 求和（decimal）
func (b *BaseModel[T]) GetSumDecimal(ctx context.Context, sql string, args ...any) (decimal.Decimal, error) {
	type sum struct {
		Num decimal.Decimal `json:"num" gorm:"column:num"`
	}
	var obj sum
	return obj.Num, b.getDB(ctx).Raw(sql, args...).Scan(&obj).Error
}
