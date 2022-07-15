package mysql

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/sunmi-OS/gocore/v2/conf/viper"
	"github.com/sunmi-OS/gocore/v2/glog"
	"github.com/sunmi-OS/gocore/v2/utils"
	"github.com/sunmi-OS/gocore/v2/utils/closes"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var defaultName = "dbDefault"

var (
	// ErrRecordNotFound record not found error, happens when haven't find any matched data when looking up with a struct
	ErrRecordNotFound = gorm.ErrRecordNotFound
	// ErrInvalidTransaction invalid transaction when you are trying to `Commit` or `Rollback`
	ErrInvalidTransaction = gorm.ErrInvalidTransaction
)

type Map struct {
	gormMaps      sync.Map
	defaultDbName string
}

var _Gorm *Map
var closeOnce sync.Once

// 初始化Gorm
func NewDB(dbname string) (g *Map) {
	err := NewOrUpdateDB(dbname)
	if err != nil {
		return nil
	}
	closeOnce.Do(func() {
		closes.AddShutdown(closes.ModuleClose{
			Name:     "Gorm Close",
			Priority: closes.GormPriority,
			Func:     Close,
		})
	})
	return _Gorm
}

// SetDefaultName 设置默认DB Name
func SetDefaultName(dbName string) {
	_Gorm.defaultDbName = dbName
}

// NewOrUpdateDB 初始化或更新Gorm
func NewOrUpdateDB(dbname string) error {
	var (
		orm *gorm.DB
		err error
	)
	if _Gorm == nil {
		_Gorm = &Map{defaultDbName: defaultName}
	}
	// second: load gorm client
	oldGorm, _ := _Gorm.gormMaps.Load(dbname)
	// first: open new gorm client
	err = utils.Retry(func() error {
		orm, err = openORM(dbname)
		if err != nil {
			glog.ErrorF("UpdateDB(%s) error:%+v", dbname, err)
			return err
		}
		return nil
	}, 5, 3*time.Second)
	// 如果NEW异常直接panic如果是Update返回error
	if err != nil {
		if oldGorm == nil {
			panic(err)
		}
		return err
	}
	// third: delete old gorm client and store the new gorm client
	_Gorm.gormMaps.Delete(dbname)
	_Gorm.gormMaps.Store(dbname, orm)
	// fourth: if old client is not nil, delete and close connection
	if oldGorm != nil {
		db, _ := oldGorm.(*gorm.DB).DB()
		if db != nil {
			err := db.Close()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// GetORM 获取默认的Gorm实例
// 目前仅支持 不传 或者仅传一个 dbname
func GetORM(dbname ...string) *gorm.DB {
	name := _Gorm.defaultDbName
	if len(dbname) != 0 {
		name = dbname[0]
	}
	v, ok := _Gorm.gormMaps.Load(name)
	if ok {
		return v.(*gorm.DB)
	}
	return nil
}

func Close() {
	_Gorm.gormMaps.Range(func(dbName, orm interface{}) bool {
		glog.WarnF("close db %s", dbName)
		_Gorm.gormMaps.Delete(dbName)
		db, _ := orm.(*gorm.DB).DB()
		if db != nil {
			db.Close()
		}
		return true
	})
}

//----------------------------以下是私有方法--------------------------------

// openORM 私有方法
func openORM(dbname string) (*gorm.DB, error) {
	//默认配置
	viper.C.SetDefault(dbname, map[string]interface{}{
		"Host":         "127.0.0.1",
		"Name":         "gocore",
		"User":         "root",
		"Passwd":       "",
		"Port":         3306,
		"MaxIdleConns": 20,
		"MaxOpenConns": 20,
		"Type":         "mysql",
		"Debug":        false,
	})
	dbHost := viper.GetEnvConfig(dbname + ".Host").String()
	dbName := viper.GetEnvConfig(dbname + ".Name").String()
	dbUser := viper.GetEnvConfig(dbname + ".User").String()
	dbPasswd := viper.GetEnvConfig(dbname + ".Passwd").String()
	dbPort := viper.GetEnvConfig(dbname + ".Port").String()
	dbType := viper.GetEnvConfig(dbname + ".Type").String()
	dbDebug := viper.GetEnvConfig(dbname + ".Debug").Bool()

	dsn := dbUser + ":" + dbPasswd + "@tcp(" + dbHost + ":" + dbPort + ")/" + dbName + "?charset=utf8mb4&parseTime=true&loc=Local"
	lc := logger.Config{
		SlowThreshold: 200 * time.Millisecond, // 慢 SQL 阈值
		LogLevel:      logger.Warn,            // Log level
		Colorful:      false,                  // 禁用彩色打印，日志平台会打印出颜色码，影响日志观察
	}
	if dbDebug {
		lc.LogLevel = logger.Info
	}
	newLogger := logger.New(
		log.New(os.Stdout, "[GORM] >> ", 64|log.Ldate|log.Lmicroseconds), // io writer
		lc,
	)
	switch dbType {
	case "mysql":
		orm, err := gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: newLogger})
		if err != nil {
			return nil, err
		}
		db, err := orm.DB()
		if err != nil {
			return nil, err
		}
		// 连接池的空闲数大小
		db.SetMaxIdleConns(viper.C.GetInt(dbname + ".MaxIdleConns"))
		// 最大打开连接数
		db.SetMaxOpenConns(viper.C.GetInt(dbname + ".MaxOpenConns"))
		return orm, nil
	default:
		return nil, fmt.Errorf("not support sql driver [%s]", dbType)
	}
}
