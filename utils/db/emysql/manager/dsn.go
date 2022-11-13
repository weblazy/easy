package manager

import (
	"gorm.io/gorm"
)

type DSN struct {
	User     string            // Username
	Password string            // Password (requires User)
	Net      string            // Network type
	Addr     string            // Network address (requires Net)
	DBName   string            // Database name
	Params   map[string]string // Connection parameters
}

type DSNParser interface {
	GetDialector(dsn string) gorm.Dialector
	ParseDSN(dsn string) (cfg *DSN, err error)
	Scheme() string
}
