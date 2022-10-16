package eredis

import (
	"context"
	"crypto/tls"
	"log"
	"strings"
	"sync"

	"github.com/go-redis/redis/v8"
	"github.com/sunmi-OS/gocore/v2/utils/hash"
	"github.com/weblazy/easy/utils/closes"
	"github.com/weblazy/easy/utils/conf/viper"
	"github.com/weblazy/easy/utils/glog"
)

var Map sync.Map
var closeOnce sync.Once

// NewRedis new Redis Client
func NewRedis(dbName string) {
	rc, err := newRedis(dbName)
	if err != nil {
		panic(err)
	}

	p := rc.Ping(context.Background())
	if p.Err() != nil {
		log.Panicln(p.Err())
	}
	Map.Store(dbName, rc)
	closeOnce.Do(func() {
		closes.AddShutdown(closes.ModuleClose{
			Name:     "Redis Close",
			Priority: closes.RedisPriority,
			Func:     Close,
		})
	})
}

func newRedis(db string) (rc *redis.ClusterClient, err error) {
	redisName, _ := splitDbName(db)
	host := viper.GetEnvConfig(redisName + ".host").String()
	port := viper.GetEnvConfig(redisName + ".port").String()
	auth := viper.GetEnvConfig(redisName + ".auth").String()
	encryption := viper.GetEnvConfig(redisName + ".encryption").Int64()
	if encryption == 1 {
		auth = hash.MD5(auth)
	}
	addr := host + port
	if !strings.Contains(addr, ":") {
		addr = host + ":" + port
	}
	if auth != "" {
		rc = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:    []string{addr},
			Password: auth,
			TLSConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		})
	} else {
		rc = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs: []string{addr},
		})
	}
	return rc, nil
}

// NewOrUpdateRedis 新建或更新redis客户端
func NewOrUpdateRedis(dbName string) error {
	rc, err := newRedis(dbName)
	if err != nil {
		return err
	}

	v, _ := Map.Load(dbName)
	Map.Delete(dbName)
	Map.Store(dbName, rc)

	if v != nil {
		err := v.(*redis.ClusterClient).Close()
		if err != nil {
			return err
		}
	}
	return nil
}

func Close() {
	Map.Range(func(dbName, rc interface{}) bool {
		glog.WarnF("close db %s", dbName)
		Map.Delete(dbName)
		err := rc.(*redis.ClusterClient).Close()
		return err == nil
	})
}

func splitDbName(db string) (redisName, dbName string) {
	kv := strings.Split(db, ".")
	if len(kv) == 2 {
		return kv[0], kv[1]
	}
	if len(kv) == 1 {
		return "redisServer", kv[0]
	}
	panic("redis dbName Mismatch")
}

// Client returns a universal redis client(ClusterClient, StubClient or SentinelClient), it depends on you config.
func UniversalClient(dbName string) redis.UniversalClient {
	if v, ok := Map.Load(dbName); ok {
		return v.(redis.UniversalClient)
	}
	return nil
}

// Cluster try to get a redis.ClusterClient
func Cluster(dbName string) *redis.ClusterClient {
	if v, ok := Map.Load(dbName); ok {
		return v.(*redis.ClusterClient)
	}
	return nil
}

// Client try to get a redis.Client
func Client(dbName string) *redis.Client {
	if v, ok := Map.Load(dbName); ok {
		return v.(*redis.Client)
	}
	return nil
}
