package grpc_client

import (
	"context"
	"fmt"
	"time"

	"github.com/weblazy/easy/utils/elog"
	"go.uber.org/zap"
	"google.golang.org/grpc/credentials/insecure"

	"google.golang.org/grpc"
)

var emptyCtx = context.Background()

// PackageName 设置包名
const PackageName = "client.fgrpc"

const grpcServiceConfig = `{"loadBalancingPolicy":"%s"}`

type GrpcClient struct {
	config *Config
	*grpc.ClientConn
	err error
}

func NewGrpcClient(config *Config) *GrpcClient {
	var ctx = context.Background()

	if config == nil {
		config = DefaultConfig()
	}
	config.BuildDialOptions()

	var dialOptions = config.dialOptions
	// 默认配置使用block
	if config.EnableBlock {
		if config.DialTimeout > time.Duration(0) {
			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(ctx, config.DialTimeout)
			defer cancel()
		}

		dialOptions = append(dialOptions, grpc.WithBlock())
	}

	if config.EnableWithInsecure {
		dialOptions = append(dialOptions, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	if config.keepAlive != nil {
		dialOptions = append(dialOptions, grpc.WithKeepaliveParams(*config.keepAlive))
	}

	//// 因为默认是开启这个配置
	//// 并且开启后，在grpc 1.40以上会导致dns多一次解析txt内容（目测是为了做grpc的load balance策略，但我们实际上不会用到）
	//// 因为这个service config dns域名通常是没有设置dns解析，所以会跳过k8s的dns，穿透到上一级的dns，而如果dns配置有问题或者不存在，那么会查询非常长的时间（通常在20s或者更长）
	//// 那么为false的时候，禁用他，可以加快我们的启动时间或者提升我们的性能
	//if !config.EnableServiceConfig {
	//	dialOptions = append(dialOptions, grpc.WithDisableServiceConfig())
	//}

	// 直接使用 default server config
	dialOptions = append(dialOptions, grpc.WithDefaultServiceConfig(fmt.Sprintf(grpcServiceConfig, config.BalancerName)), grpc.FailOnNonTempDialError(config.EnableFailOnNonTempDialError))

	startTime := time.Now()
	cc, err := grpc.DialContext(ctx, config.Addr, dialOptions...)

	client := &GrpcClient{
		config:     config,
		ClientConn: cc,
	}

	if err != nil {
		elog.ErrorCtx(emptyCtx, "dial grpc server", elog.FieldError(err), elog.FieldName(config.Name), zap.String("addr", config.Addr), elog.FieldCost(time.Since(startTime)))
		return client
	}

	elog.InfoCtx(emptyCtx, "start grpc client", elog.FieldName(config.Name), elog.FieldCost(time.Since(startTime)))
	return client
}

// Error 错误信息
func (c *GrpcClient) Error() error {
	return c.err
}
