# 业务指标监控接入

`fmetric` 只是对 `prometheus` 的简单封装, 仅仅将初始化和注册指标包装成一个函数而已. 如果有疑惑请先学习 `prometheus` 相关知识.

接入需要注意的是:

1. 非业务相关 label 不需要关注, 例如: app, projectEnv (k8s prometheus operator 会在收集时自动注入 k8s pod 相关信息)
2. 注意 labels 组合种类总数**必须**是常数级别有限的

对于业务指标, 绝大多数情况 `Counter` 类型就够了, 别的类型也是同理(请确保你在充分了解不同指标类型后选择合适的类型使用). 下面用此类型为例.

```go

// 1. 在项目里初始化指标
var SomeBizCounter = fmetric.CounterVecOpts{
	    // 指标名称命名要清晰, 并且是下划线形式
		Name:      "monitor_some_biz_total",
		// help 字段增加指标说明
		Help:      "Total number of some biz logic on the server",
		// 业务 label 名称
		Labels:    []string{"channel", "status"},
	}.Build()


// 2. 在适当的业务逻辑中上报指标
// WithLabelValues 值顺序必须和 Labels 声明顺序一致
// 真正的业务使用 label 值可以定义成全局常量, 避免拼写错误
SomeBizCounter.WithLabelValues("channel1", "OK").Inc()
// other case
SomeBizCounter.WithLabelValues("channel2", "ERROR").Inc()
```
