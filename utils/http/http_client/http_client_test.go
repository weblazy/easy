package http_client

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
	"github.com/weblazy/easy/utils/http/http_client/http_client_config"
)

func TestNewHttpClient(t *testing.T) {
	convey.Convey("TestNewHttpClient", t, func() {
		cfg := http_client_config.DefaultConfig()
		client := NewHttpClient(cfg)
		request := client.Request.SetContext(context.Background())
		resp, err := request.Get("https://www.baidu.com/")
		body := string(resp.Body())
		convey.So(body, convey.ShouldNotBeNil)
		convey.So(err, convey.ShouldBeNil)
	})
}
