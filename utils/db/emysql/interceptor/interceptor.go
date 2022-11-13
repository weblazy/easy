package interceptor

import (
	"gorm.io/gorm"
)

const (
	TypeGorm = "gorm"
)

var (
	// ErrRecordNotFound record not found error, happens when haven't find any matched data when looking up with a struct
	ErrRecordNotFound = gorm.ErrRecordNotFound
	// ErrInvalidTransaction invalid transaction when you are trying to `Commit` or `Rollback`
	ErrInvalidTransaction = gorm.ErrInvalidTransaction
)

// // 确保在生产不要开 debug
// func DebugInterceptor(compName string, dsn *manager.DSN, op string, options *mysql_config.Config) func(mysql_config.Handler) mysql_config.Handler {
// 	return func(next mysql_config.Handler) mysql_config.Handler {
// 		return func(db *gorm.DB) {
// 			beg := time.Now()
// 			next(db)
// 			duration := time.Since(beg)
// 			if db.Error != nil {
// 				elog.ErrorCtx(db.Statement.Context, "fgorm.response", elog.MakeReqResError(1, compName, dsn.Addr+"/"+dsn.DBName, duration, logSQL(db.Statement.SQL.String(), db.Statement.Vars, true), db.Error.Error()))
// 			} else {
// 				elog.InfoCtx(db.Statement.Context, "fgorm.response", elog.MakeReqResInfo(1, compName, dsn.Addr+"/"+dsn.DBName, duration, logSQL(db.Statement.SQL.String(), db.Statement.Vars, true), fmt.Sprintf("%v", db.Statement.Dest)))
// 			}
// 		}
// 	}
// }
