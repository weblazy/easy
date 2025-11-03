package emysql

import (
	"context"
	"sync"

	"emperror.dev/errors"
	"gorm.io/gorm"
)

type TransactionKeyType string

const TxOpen TransactionKeyType = "TxOpen"
const TxDBMap TransactionKeyType = "TxDBMap"

// 注意:该事务不是并发安全的
func Transaction(ctx context.Context, f func(context.Context) error) (err error) {
	// context 打事务启动标
	ctx = context.WithValue(ctx, TxOpen, true)
	// 初始化事务组件
	ctx = context.WithValue(ctx, TxDBMap, &sync.Map{})

	err = f(ctx)
	if err == nil {
		// 提交事务
		return commit(ctx)
	} else {
		// 回滚事务
		rollbackErr := rollback(ctx)
		if rollbackErr != nil {
			return errors.Wrap(err, rollbackErr.Error())
		}
		return err
	}
}

func checkTransaction(ctx context.Context, dbName string) *gorm.DB {
	// 判断是否开启了事务
	if ctx.Value(TxOpen) == nil {
		return nil
	}

	txDBMap, ok := ctx.Value(TxDBMap).(*sync.Map)
	if !ok {
		return nil
	}

	tx, ok := txDBMap.Load(dbName)
	if ok && tx != nil {
		db, ok := tx.(*gorm.DB)
		if ok {
			return db
		}
	}

	// 如果没有开启事务需要开启事务 获取实例开启事务
	db := getMysql(ctx, dbName)
	if db == nil {
		return nil
	}
	db = db.Begin()
	// 将tx放入context
	txDBMap.Store(dbName, db)
	return db
}

func commit(ctx context.Context) (err error) {
	// 判断是否开启了事务
	if ctx.Value(TxOpen) == nil || ctx.Value(TxDBMap) == nil {
		return nil
	}
	// 获取所有的tx
	txDBMap, ok := ctx.Value(TxDBMap).(*sync.Map)
	if !ok {
		return errors.New("TxDBMapTypeErr")
	}
	txDBMap.Range(func(key, value any) bool {
		tx, ok := value.(*gorm.DB)
		if !ok {
			err = errors.New("TxDBTypeErr")
			return true
		}
		tx.Commit()
		return true
	})
	return err
}

func rollback(ctx context.Context) (err error) {
	// 判断是否开启了事务
	if ctx.Value(TxOpen) == nil || ctx.Value(TxDBMap) == nil {
		return nil
	}
	// 获取所有的tx
	txDBMap, ok := ctx.Value(TxDBMap).(*sync.Map)
	if !ok {
		return errors.New("TxDBMapTypeErr")
	}
	txDBMap.Range(func(key, value any) bool {
		tx, ok := value.(*gorm.DB)
		if !ok {
			err = errors.New("TxDBTypeErr")
			return true
		}
		tx.Rollback()
		return true
	})
	return err
}
