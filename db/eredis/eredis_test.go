package eredis

import (
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestNewRedisClient(t *testing.T) {
	convey.Convey("TestNewRedisClient", t, func() {
		// cfg := eredis_config.DefaultConfig()
		// cfg.Addr = "127.0.0.1:16379"
		// cfg.Name = "user_redis"
		// client := NewRedisClient(cfg)
		// cmd := client.Get(context.Background(), "test")
		// resp, err := cmd.Result()
		// fmt.Printf("%#v\n", resp)
		// convey.So(resp, convey.ShouldNotBeNil)
		// convey.So(err, convey.ShouldBeNil)
	})
}
