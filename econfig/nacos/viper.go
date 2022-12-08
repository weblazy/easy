package nacos

import (
	"context"

	"github.com/BurntSushi/toml"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/weblazy/easy/econfig/eviper"
	"github.com/weblazy/easy/elog"
	"go.uber.org/zap"
)

type configParam struct {
	group   string
	dataIds []string
}

type ViperToml struct {
	Viper             *eviper.Viper
	dataIdOrGroupList []configParam
	callbackList      map[string]func(namespace, group, dataId, data string)
}

var emptyCtx = context.Background()

// GetViper 获取VT实例
func GetViper() *ViperToml {
	return nacosHarder.vt
}

// NacosToViper 同步Nacos读取的配置注入Viper
func (vt *ViperToml) NacosToViper() {
	s, err := vt.GetConfig()
	if err != nil {
		panic(err)
	}
	vt.Viper.MergeViperFromString(s)
}

// SetBaseConfig 注入基础配置
func (vt *ViperToml) SetBaseConfig(configs string) {
	vt.Viper.MergeViperFromString(configs)
}

// GetConfig 获取整套配置文件
func (vt *ViperToml) GetConfig() (string, error) {
	if nacosHarder.local {
		configs, err := nacosHarder.icc.GetConfig(vo.ConfigParam{})
		if err != nil {
			return "", err
		}
		return configs, nil
	}
	var configs = ""
	for _, dataIdOrGroup := range vt.dataIdOrGroupList {
		group := dataIdOrGroup.group
		for _, v := range dataIdOrGroup.dataIds {
			content, err := nacosHarder.icc.GetConfig(vo.ConfigParam{DataId: v, Group: group})
			if err != nil {
				return "", err
			}
			configs += "\r\n" + content
			// 注册回调
			err = nacosHarder.icc.ListenConfig(vo.ConfigParam{
				DataId:   v,
				Group:    group,
				OnChange: vt.callbackList[group+v],
			})
			if err != nil {
				elog.ErrorCtx(emptyCtx, group+"\r\n"+v+"\r\n ListenConfig Error", zap.Error(err))
			}
		}
	}
	return configs, nil
}

// GetConfigParse 获取配置并且绑定结构体
func (vt *ViperToml) GetConfigParse(confPtr interface{}) error {
	config, err := vt.GetConfig()
	if err != nil {
		return err
	}
	_, err = toml.Decode(config, confPtr)
	if err != nil {
		return err
	}
	return nil
}

// SetDataIds 设置需要读取哪些配置
// 配置默认回调方法更新配置
func (vt *ViperToml) SetDataIds(group string, dataIds ...string) {
	vt.dataIdOrGroupList = append(vt.dataIdOrGroupList, configParam{group: group, dataIds: dataIds})
	for _, v := range dataIds {
		vt.callbackList[group+v] = func(namespace, group, dataId, data string) {
			vt.NacosToViper()
			elog.WarnCtx(emptyCtx, namespace+"\r\n"+group+"\r\n"+dataId+"\r\n"+data+"\r\n Update Config")
		}
	}
}

// SetCallBackFunc 配置自定义回调方法
func (vt *ViperToml) SetCallBackFunc(group, dataId string, callback func(namespace, group, dataId, data string)) {
	vt.callbackList[group+dataId] = func(namespace, group, dataId, data string) {
		vt.NacosToViper()
		callback(namespace, group, dataId, data)
	}
}
