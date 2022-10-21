package mysql

import (
	"context"
	"fmt"

	_ "github.com/weblazy/easy/utils/db/mysql/internal/dsn"
	"github.com/weblazy/easy/utils/db/mysql/manager"
	"github.com/weblazy/easy/utils/elog"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var emptyCtx = context.Background()

type MysqlClient struct {
	*gorm.DB
	dsnParser manager.DSNParser
}

// Option 可选项
type Option func(c *Config)

// NewMysqlClient ...
func NewMysqlClient(config *Config, options ...Option) (*MysqlClient, error) {
	mysqlClient := MysqlClient{}

	gormCfg := gorm.Config{}
	// 不开启 raw debug 时, 关闭 gorm 原生日志
	if !config.RawDebug {
		gormCfg.Logger = logger.Discard
	}

	if config.Debug {
		options = append(options, WithInterceptor(debugInterceptor))
	}

	if config.EnableTraceInterceptor {
		options = append(options, WithInterceptor(traceInterceptor))
	}

	if config.EnableMetricInterceptor {
		options = append(options, WithInterceptor(metricInterceptor))
	}

	for _, option := range options {
		option(config)
	}

	// todo 设置补齐超时时间, 解析重写 config.DSN 参数
	// timeout 1s
	// readTimeout 5s
	// writeTimeout 5s
	err := mysqlClient.setDSNParser(config.Dialect)
	if err != nil {
		elog.ErrorCtx(emptyCtx, "setDSNParser err", zap.Error(err), zap.String("dialect", config.Dialect))
	}

	config.DsnCfg, err = mysqlClient.dsnParser.ParseDSN(config.DSN)
	fmt.Printf("%#v\n", config.DsnCfg)

	if err == nil {
		elog.InfoCtx(emptyCtx, "start db", zap.String("addr", config.DsnCfg.Addr), zap.String("name", config.DsnCfg.DBName))
	} else {
		elog.ErrorCtx(emptyCtx, "start db", zap.Error(err))
	}

	db, err := gorm.Open(mysqlClient.dsnParser.GetDialector(config.DSN), &gormCfg)
	if err != nil {
		return nil, err
	}

	if config.RawDebug {
		db = db.Debug()
	}

	gormDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// 设置默认连接配置
	gormDB.SetMaxIdleConns(config.MaxIdleConns)
	gormDB.SetMaxOpenConns(config.MaxOpenConns)

	if config.ConnMaxLifetime != 0 {
		gormDB.SetConnMaxLifetime(config.ConnMaxLifetime)
	}

	var lastErr error
	replace := func(processor Processor, callbackName string, interceptors ...Interceptor) {
		handler := processor.Get(callbackName)
		for _, interceptor := range config.interceptors {
			handler = interceptor(config.Name, config.DsnCfg, callbackName, config)(handler)
		}

		err := processor.Replace(callbackName, handler)
		if err != nil {
			lastErr = err
		}
	}

	replace(db.Callback().Create(), "gorm:create", config.interceptors...)
	replace(db.Callback().Update(), "gorm:update", config.interceptors...)
	replace(db.Callback().Delete(), "gorm:delete", config.interceptors...)
	replace(db.Callback().Query(), "gorm:query", config.interceptors...)
	// replace(db.Callback().Row(), "gorm:row", config.interceptors...)
	replace(db.Callback().Raw(), "gorm:raw", config.interceptors...)

	if lastErr != nil {
		return nil, lastErr
	}
	mysqlClient.DB = db
	return &mysqlClient, nil
}

// WithContext ...
func (m *MysqlClient) WithContext(ctx context.Context) *MysqlClient {
	m.Statement.Context = ctx
	return m
}

func (c *MysqlClient) setDSNParser(dialect string) error {
	dsnParser := manager.Get(dialect)
	if dsnParser == nil {
		return fmt.Errorf("invalid support Dialect: %s", dialect)
	}
	c.dsnParser = dsnParser
	return nil
}

// WithInterceptor 设置自定义拦截器
func WithInterceptor(is ...Interceptor) Option {
	return func(c *Config) {
		if c.interceptors == nil {
			c.interceptors = make([]Interceptor, 0)
		}
		c.interceptors = append(c.interceptors, is...)
	}
}
