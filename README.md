# easy
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