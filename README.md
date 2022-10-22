# easy
如果大家觉得好用,右上角帮忙点个star吧。(^_^)
> 欢迎感兴趣的小伙伴一同开发,收集日常好用的golang工具包。
# 联系我们
- 技术支持/合作/咨询请联系作者QQ: 2276282419
- 作者邮箱: 2276282419@qq.com
- 即时通讯技术交流QQ群: 33280853
##### 单元测试
```
go test -coverpkg=./... -coverprofile=coverage.data -timeout=5s ./...
go tool cover -html=coverage.data -o coverage.html
````
[![Go Report Card](https://goreportcard.com/badge/github.com/sunmi-OS/gocore)](https://goreportcard.com/report/github.com/sunmi-OS/gocore/v2.0.9)

##### 功能列表

- grpc_server: 
  - 日志插件
  - metric插件
  - recovery插件
  - timeout插件
  - trace插件
- grpc_client: 
  - 日志插件
  - metric插件
  - timeout插件
  - trace插件
- http_server: gin
    - 日志插件
    - metric插件
    - recovery插件
    - timeout插件
    - trace插件
- http_client: resty
  - 日志插件
  - metric插件
  - timeout插件
  - trace插件
- db: gorm
  - 日志插件
  - metric插件
  - timeout插件
  - trace插件
  - 脚手架: orm
- redis: go-redis
  - 日志插件
  - metric插件
  - timeout插件
  - trace插件
- log: zap
- config: viper
- 监控面板: prometheus+grafana
- 告警
- 脚手架: easy-cli


