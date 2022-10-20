package http_server

import (
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestNewHttpServer(t *testing.T) {
	convey.Convey("TestNewHttpServer", t, func() {
		// cfg := http_server_config.DefaultConfig()
		// server, err := NewHttpServer(cfg)
		// convey.So(err, convey.ShouldBeNil)
		// err = server.Start()
		// convey.So(err, convey.ShouldBeNil)
	})
}

// func TestNewHttpServerViper(t *testing.T) {
// 	convey.Convey("test config", t, func() {
// 		cfg := viper.New()
// 		cfg.SetConfigType("toml")
// 		s := strings.NewReader(`Name="go"
// [http_server]
// name=6666
// level=50`)

// 		err := cfg.ReadConfig(s)
// 		convey.So(err, convey.ShouldBeNil)
// 		resp, err := NewHttpServerViper("http_server", cfg)
// 		fmt.Printf("%#v\n", resp.Config)
// 		convey.So(resp, convey.ShouldNotBeNil)
// 		convey.So(err, convey.ShouldBeNil)
// 	})
// }
