package dsn

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMysqlDSNParser_ParseDSN(t *testing.T) {
	dsn := "user:password@tcp(localhost:9910)/dbname?charset=utf8&parseTime=True"
	parser := &MysqlDSNParser{}
	cfg, err := parser.ParseDSN(dsn)
	assert.NoError(t, err)
	assert.Equal(t, "user", cfg.User)
	assert.Equal(t, "password", cfg.Password)
	assert.Equal(t, "dbname", cfg.DBName)
	assert.Equal(t, "localhost:9910", cfg.Addr)
	assert.Equal(t, "tcp", cfg.Net)
	assert.Equal(t, "utf8", cfg.Params["charset"])
	assert.Equal(t, "True", cfg.Params["parseTime"])
	fmt.Println(cfg)
}
