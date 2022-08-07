package nacos

import (
	"os"

	"github.com/sunmi-OS/gocore/v2/utils"

	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/spf13/cast"
)

type nacos struct {
	icc   config_client.IConfigClient
	vt    *ViperToml
	local bool
}

const (
	LogDebug = "debug"
	LogWarn  = "warn"
	LogError = "error"
	LogInfo  = "info"

	_NacosNamespaceId = "NACOS_NAMESPACE_ID"
	_NacosEndpoint    = "NACOS_ENDPOINT"
	_NacosAccessKey   = "NACOS_ACCESS_KEY"
	_NacosSecretKey   = "NACOS_SECRET_KEY"
	_NacosRegionId    = "NACOS_REGION_ID"

	_NamespaceId = "NAMESPACE_ID"
	_Endpoint    = "ENDPOINT"
	_AccessKey   = "ACCESS_KEY"
	_SecretKey   = "SECRET_KEY"
	_RegionId    = "REGION_ID"

	_DefaultRegionId = "cn-hangzhou"

	_NacosScheme      = "NACOS_SCHEME"
	_NacosContextPath = "NACOS_CONTEXT_PATH"
	_NacosHost        = "NACOS_HOST"
	_NacosPort        = "NACOS_PORT"
)

var nacosHarder = &nacos{
	vt: &ViperToml{
		callbackList: make(map[string]func(namespace, group, dataId, data string)),
	},
}

// NewNacosEnv 注入Nacos配置文件
// 兼容endpoint 和 service 两种方式
func NewNacosEnv() {
	// 获取nacos鉴权、地区、命名空间
	namespaceId := utils.Either(os.Getenv(_NacosNamespaceId), os.Getenv(_NamespaceId))
	accessKey := utils.Either(os.Getenv(_NacosAccessKey), os.Getenv(_AccessKey))
	secretKey := utils.Either(os.Getenv(_NacosSecretKey), os.Getenv(_SecretKey))
	regionID := utils.Either(os.Getenv(_NacosRegionId), os.Getenv(_RegionId))
	if namespaceId == "" || accessKey == "" || secretKey == "" {
		panic("The configuration file cannot be empty.")
	}
	if regionID == "" {
		regionID = _DefaultRegionId
	}
	ccConfig := &constant.ClientConfig{
		NamespaceId: namespaceId,
		AccessKey:   accessKey,
		SecretKey:   secretKey,
		RegionId:    regionID,
		LogLevel:    LogError,
	}

	// 读取service地址，如果有service优先使用service连接方式
	nacosHost := os.Getenv(_NacosHost)
	nacosPort := os.Getenv(_NacosPort)
	if nacosHost != "" && nacosPort != "" {
		err := NewNacos(ccConfig, constant.ServerConfig{
			IpAddr:      nacosHost,
			Port:        cast.ToUint64(nacosPort),
			Scheme:      os.Getenv(_NacosScheme),
			ContextPath: os.Getenv(_NacosContextPath),
		})
		if err != nil {
			panic(err)
		}
		return
	}

	// 如果未使用service检查是否配置了endpoint兼容acm
	endpoint := utils.Either(os.Getenv(_NacosEndpoint), os.Getenv(_Endpoint))
	if endpoint == "" {
		panic("The configuration file cannot be empty.")
	}
	ccConfig.Endpoint = endpoint
	err := NewNacos(ccConfig)
	if err != nil {
		panic(err)
	}
}

// NewNacos 注入Nacos配置文件
func NewNacos(ccConfig *constant.ClientConfig, scConfigs ...constant.ServerConfig) error {
	defaultClientConfig(ccConfig)
	configClient, err := clients.NewConfigClient(vo.NacosClientParam{
		ClientConfig:  ccConfig,
		ServerConfigs: scConfigs,
	})
	if err != nil {
		return err
	}
	nacosHarder.icc = configClient
	return nil
}

// defaultClientConfig 使用nacos时的默认配置
func defaultClientConfig(ccConfig *constant.ClientConfig) {
	if ccConfig.TimeoutMs == 0 {
		ccConfig.TimeoutMs = 5000
	}
	if ccConfig.BeatInterval == 0 {
		ccConfig.BeatInterval = 5 * 1000
	}
	if ccConfig.LogDir == "" {
		ccConfig.LogDir = "./nacos/logs"
	}
	if ccConfig.CacheDir == "" {
		ccConfig.CacheDir = "./nacos/cache"
	}
}

// GetConfig 获取单条配置
func GetConfig(group string, dataIds string) string {
	content, err := nacosHarder.icc.GetConfig(vo.ConfigParam{
		DataId: dataIds,
		Group:  group})
	if err != nil {
		return ""
	}
	return content
}

// CallBackFunc 参数更新回调方法
func CallBackFunc(group, dataId string, callback func(namespace, group, dataId, data string)) error {
	err := nacosHarder.icc.ListenConfig(vo.ConfigParam{
		DataId:   dataId,
		Group:    group,
		OnChange: callback,
	})
	if err != nil {
		return err
	}

	return nil
}
