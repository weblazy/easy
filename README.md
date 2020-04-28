# easy
golang工具包

##### 作者邮箱:227622419@qq.com
##### 单元测试
```
go test -coverpkg=./... -coverprofile=coverage.data -timeout=5s ./...
go tool cover -html=coverage.data -o coverage.html
````