package grpc_server

import (
	"context"
	"strings"

	"github.com/weblazy/easy/utils/elog"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/reflection"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/test/bufconn"

	"net"
)

var emptyCtx = context.Background()

const (
	// PackageName 包名
	PackageName       = "server.fgrpc"
	networkTypeBufNet = "bufnet"
)

// Component ...
type GrpcServer struct {
	config  *Config
	logConf *elog.LogConf
	*grpc.Server
	listener net.Listener
	quit     chan struct{}
}

func NewGrpcServer(config *Config, logConf *elog.LogConf) *GrpcServer {
	if config == nil {
		config = DefaultConfig()
	}
	config.BuildServerOptions()

	newServer := grpc.NewServer(config.serverOptions...)

	if config.EnableServerReflection {
		elog.InfoCtx(emptyCtx, "enable grpc server reflection")
		reflection.Register(newServer)
	}

	if config.EnableHealth {
		elog.InfoCtx(emptyCtx, "enable grpc health")
		healthpb.RegisterHealthServer(newServer, health.NewServer())
	}

	return &GrpcServer{
		config:   config,
		logConf:  logConf,
		Server:   newServer,
		listener: nil,
		quit:     make(chan struct{}),
	}
}

// Name 配置名称
func (c *GrpcServer) Name() string {
	return c.config.Name
}

// PackageName 包名
func (c *GrpcServer) PackageName() string {
	return PackageName
}

// Init 初始化
func (c *GrpcServer) Init() error {
	var (
		listener net.Listener
		err      error
	)
	// gRPC测试listener
	if c.config.Network == networkTypeBufNet {
		listener = bufconn.Listen(1024 * 1024)
		c.listener = listener
		return nil
	}
	// 正式listener
	listener, err = net.Listen(c.config.Network, c.config.Address())
	if err != nil {
		elog.ErrorCtx(emptyCtx, "new grpc server err", elog.FieldError(err))
	}
	c.config.Port = listener.Addr().(*net.TCPAddr).Port

	c.listener = listener
	return nil
}

// Start implements server.Component interface.
func (c *GrpcServer) Start() error {
	err := c.Server.Serve(c.listener)
	return err
}

// Stop implements server.Component interface
// it will terminate echo server immediately
func (c *GrpcServer) Stop() error {
	c.Server.Stop()
	return nil
}

// GracefulStop implements server.Component interface
// it will stop echo server gracefully
func (c *GrpcServer) GracefulStop(ctx context.Context) error {
	go func() {
		c.Server.GracefulStop()
		close(c.quit)
	}()

	select {
	case <-ctx.Done():
		elog.WarnCtx(ctx, "grpc graceful shutdown timeout")
		return ctx.Err()
	case <-c.quit:
		elog.InfoCtx(ctx, "grpc graceful shutdown success")
		return nil
	}
}

// Address 服务地址
func (c *GrpcServer) Address() string {
	return c.config.Address()
}

// Listener listener信息
func (c *GrpcServer) Listener() net.Listener {
	return c.listener
}

// IsBufNet 返回是不是测试网络类型
func (c *GrpcServer) IsBufNet() bool {
	return c.config.Network == networkTypeBufNet
}

// getPeerIP 获取对端ip
func getPeerIP(ctx context.Context) string {
	// 从grpc里取对端ip
	pr, ok2 := peer.FromContext(ctx)
	if !ok2 {
		return ""
	}
	if pr.Addr == net.Addr(nil) {
		return ""
	}
	addSlice := strings.Split(pr.Addr.String(), ":")
	if len(addSlice) > 1 {
		return addSlice[0]
	}
	return ""
}
