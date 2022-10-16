# grpc client

## 配置说明

```go
type Config struct {
	Debug            bool          // 是否开启调试，默认不开启, 开启可以打印请求日志
	Addr             string        // 连接地址，直连为 127.0.0.1:9001，服务发现为 nacos:///appname
	BalancerName     string        // 负载均衡方式，默认round robin
	DialTimeout      time.Duration // 连接超时，默认3s
	ReadTimeout      time.Duration // 读超时，默认1s
	SlowLogThreshold time.Duration // 慢日志记录的阈值，默认600ms
	EnableBlock      bool          // 是否开启阻塞，默认开启
	// EnableOfficialGrpcLog        bool          // 是否开启官方grpc日志，默认关闭 // blog 和 zap 类型不兼容, 没法做
	EnableWithInsecure           bool // 是否开启非安全传输，默认开启
	EnableMetricInterceptor      bool // 是否开启监控，默认开启
	EnableTraceInterceptor       bool // 是否开启链路追踪，默认开启
	EnableAppNameInterceptor     bool // 是否开启传递应用名，默认开启
	EnableTimeoutInterceptor     bool // 是否开启超时传递，默认开启
	EnableAccessInterceptor      bool // 是否开启记录请求数据，默认不开启
	EnableAccessInterceptorReq   bool // 是否开启记录请求参数，默认不开启
	EnableAccessInterceptorRes   bool // 是否开启记录响应参数，默认不开启
	EnableServiceConfig          bool // 是否开启服务配置，默认关闭
	EnableFailOnNonTempDialError bool
}
```

## 连接服务问题

默认情况下(我们组件逻辑), grpc 连接会设置 3s 超时, 超时没连接上就会 `panic`.

这样做是为了保证依赖的服务正常, 也是为了尽早暴露错误, `fail fast`.

但是测试环境不稳定, 或者我们允许循环依赖情况出现时, 默认配置没法满足需求.

如果你需要支持上述场景, 需要增加配置:

```toml
enableBlock = false
```

这种情况 grpc 连接不会 block 程序启动. 后续依赖服务正常后 grpc client 功能也会正常, 不需要做重启等操作.

### EnableBlock 和 OnFail 参数区别

更新: OnFail 参数已删除, 不再支持 `onFail = "error"` 这种行为.

先说结论, 大多数时候你应该使用 `EnableBlock = false` 配置.

`OnFail` 是我们 component 通用参数, 基本意义就是开发者是否将当前 component 视为强依赖.

`OnFail` 参数控制 `grpc.Dial` 建立连接时出错时的处理方式, 默认为 `panic` 结束程序.

当 `OnFail` 设置为 `error` 时, 仅仅会在连接失败时打印错误日志, 但是 conn 返回值其实是 nil, 所以当你的程序后续依赖这个连接时, 程序依旧会 `panic`. 所以这种情况仅适合程序不依赖这个 grpc client 逻辑时使用.

`EnableBlock` 设置为 false 时, `grpc.Dial` 不会的返回错误, 所以 `OnFail` 参数其实会没有效果.
