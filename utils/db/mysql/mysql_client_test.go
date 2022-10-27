package mysql

import (
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

type User struct {
	Id   int64
	Name string
}

func (*User) TableName() string {
	return "user"
}

func TestNewMysqlClient(t *testing.T) {
	convey.Convey("TestNewMysqlClient", t, func() {
		// cfg := mysql_config.DefaultConfig()
		// cfg.DSN = "root:123456@tcp(localhost:13306)/test?charset=utf8mb4&collation=utf8mb4_general_ci&parseTime=True&loc=Local&timeout=1s&readTimeout=3s&writeTimeout=3s"
		// client, err := NewMysqlClient(cfg)
		// convey.So(err, convey.ShouldBeNil)
		// resp := User{}
		// err = client.WithContext(context.Background()).Where("id != ?", 1).Find(&resp).Error
		// convey.So(err, convey.ShouldBeNil)
		// fmt.Printf("resp%#v\n", resp)
		// convey.So(resp, convey.ShouldNotBeNil)
		// convey.So(err, convey.ShouldBeNil)
	})
}
