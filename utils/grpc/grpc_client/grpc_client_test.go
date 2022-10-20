package grpc_client

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
	"github.com/weblazy/easy/utils/grpc/grpc_client_config"
	"github.com/weblazy/easy/utils/grpc/proto/user"
)

func TestNewGrpcClient(t *testing.T) {
	convey.Convey("TestNewGrpcClient", t, func() {
		cfg := grpc_client_config.DefaultConfig()
		client := NewGrpcClient(cfg)
		userClient := user.NewUserServiceClient(client)
		resp, err := userClient.GetUserInfo(context.Background(), &user.GetUserInfoRequest{})
		convey.So(err, convey.ShouldBeNil)
		convey.So(resp, convey.ShouldNotBeNil)
	})
}
