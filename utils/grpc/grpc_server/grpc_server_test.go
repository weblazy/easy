package grpc_server

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
	"github.com/weblazy/easy/utils/elog"
	"github.com/weblazy/easy/utils/grpc/grpc_server/grpc_server_config"
	"github.com/weblazy/easy/utils/grpc/proto/user"
)

func TestNewGrpcServer(t *testing.T) {
	convey.Convey("TestNewGrpcServer", t, func() {
		cfg := grpc_server_config.DefaultConfig()
		server := NewGrpcServer(cfg, &elog.LogConf{})
		user.RegisterUserServiceServer(server.Server, &User{})
		// server.RegisterService(server.Server)
		err := server.Init()
		convey.So(err, convey.ShouldBeNil)
		err = server.Start()
		convey.So(err, convey.ShouldBeNil)

	})
}

type User struct {
	user.UserServiceServer
}

func (*User) GetUserInfo(ctx context.Context, req *user.GetUserInfoRequest) (*user.GetUserInfoResponse, error) {
	return &user.GetUserInfoResponse{
		Detail: &user.User{
			Name: "lazy",
		},
	}, nil
}
