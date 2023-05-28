package econfig

import (
	"os"
	"strings"

	"github.com/weblazy/easy/econfig/eviper"
	"github.com/weblazy/easy/econfig/nacos"
)

var GlobalViper *eviper.Viper

func InitGlobalViper(config interface{}, localConfig ...string) {
	switch os.Getenv(EasyConfigType) {
	case LocalType:
		GlobalViper = eviper.NewViperFromString(localConfig[0])
	case FielType:
		GlobalViper = eviper.NewViperFromFile("", os.Getenv(EasyConfigFile))
	case NacosType:
		nacos.NewNacosEnv()
		vt := nacos.GetViper()
		vt.SetDataIds(os.Getenv("ServiceName"), os.Getenv("DataId"))
		// 注册配置更新回调
		vt.NacosToViper()
		GlobalViper = vt.Viper
	default:
		GlobalViper = eviper.NewViperFromString(localConfig[0])
	}
	GlobalViper.Unmarshal(&config)
}

func GetEnvConfig(key string) string {
	env := os.Getenv(strings.Replace(strings.ToUpper(key), ".", "_", -1))
	if env != "" {
		return env
	}
	return GlobalViper.GetString(key)
}
