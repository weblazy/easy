package grpc_server

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
	"github.com/weblazy/easy/grpc/proto/user"
)

func TestNewGrpcServer(t *testing.T) {
	convey.Convey("TestNewGrpcServer", t, func() {
		// cfg := grpc_server_config.DefaultConfig()
		// server := NewGrpcServer(cfg, &elog.LogConf{})
		// user.RegisterUserServiceServer(server, &User{})
		// err := server.Init()
		// convey.So(err, convey.ShouldBeNil)
		// err = server.Start()
		// convey.So(err, convey.ShouldBeNil)
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
