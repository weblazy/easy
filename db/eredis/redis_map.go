package eredis

import (
	"sync"

	"github.com/weblazy/easy/db/eredis/eredis_config"
	"github.com/weblazy/easy/econfig"
)

var RedisMap sync.Map

// GetRedis return a RedisClient
func GetRedis(dbName string) *RedisClient {
	if v, ok := RedisMap.Load(dbName); ok {
		return v.(*RedisClient)
	}
	conf := eredis_config.DefaultConfig()
	econfig.GlobalViper.UnmarshalKey(dbName, conf)
	redisClient := NewRedisClient(conf)
	RedisMap.Store(dbName, redisClient)
	return redisClient
}
