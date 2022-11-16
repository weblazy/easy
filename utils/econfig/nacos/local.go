package nacos

import (
	"io/ioutil"

	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

type LocalNacos struct {
	configs string
	config_client.IConfigClient
}

// SetLocalConfigFile 注入本地配置 指定目录
func SetLocalConfigFile(filePath string) {
	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	SetLocalConfig(string(bytes))
}

// SetLocalConfig 注入本地配置
func SetLocalConfig(configs string) {
	localNacos := NewLocalNacos(configs)
	nacosHarder.icc = localNacos
	nacosHarder.local = true
}

func NewLocalNacos(configs string) config_client.IConfigClient {
	return &LocalNacos{configs: configs}
}

func (l *LocalNacos) GetConfig(param vo.ConfigParam) (string, error) {
	str := l.configs
	return str, nil
}

func (l *LocalNacos) PublishConfig(param vo.ConfigParam) (bool, error) {
	return true, nil
}

func (l *LocalNacos) DeleteConfig(param vo.ConfigParam) (bool, error) {
	return true, nil
}

func (l *LocalNacos) ListenConfig(params vo.ConfigParam) (err error) {
	return nil
}
