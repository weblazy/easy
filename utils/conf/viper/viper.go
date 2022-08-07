package viper

import (
	"bytes"
	"log"
	"os"
	"path"
	"strings"
	"sync"

	"github.com/BurntSushi/toml"
	"github.com/spf13/viper"
	"github.com/sunmi-OS/gocore/v2/utils"
	"github.com/sunmi-OS/gocore/v2/utils/file"
)

type Viper struct {
	C *viper.Viper
}

var multipleViper sync.Map
var C = viper.New()

func NewConfigToToml(configs string) {
	C.SetConfigType("toml")
	CheckToml(configs)
	err := C.ReadConfig(bytes.NewBuffer([]byte(configs)))
	if err != nil {
		print(err)
	}
}

func MergeConfigToToml(configs string) {
	CheckToml(configs)
	C.SetConfigType("toml")
	err := C.MergeConfig(bytes.NewBuffer([]byte(configs)))
	if err != nil {
		print(err)
	}
}

func CheckToml(configs string) {
	var tmp interface{}
	if _, err := toml.Decode(configs, &tmp); err != nil {
		log.Fatalf("Error decoding TOML: %s", err)
		return
	}
}

func NewConfig(filePath string, fileName string) {
	C = newConfig(filePath, fileName).C
}

func newConfig(filePath string, fileName string) *Viper {
	v := viper.New()
	v.SetConfigName(fileName)
	//filePath支持相对路径和绝对路径 etc:"/a/b" "b" "./b"
	if filePath[:1] != "/" {
		v.AddConfigPath(path.Join(file.GetPath(), filePath))
	} else {
		v.AddConfigPath(filePath)
	}
	v.WatchConfig()
	// 找到并读取配置文件并且 处理错误读取配置文件
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}
	return &Viper{v}
}

func BuildVipers(filePath string, fileName ...string) {
	for _, v := range fileName {
		_, found := multipleViper.Load(v)
		if !found { //can not remap
			A := newConfig(filePath, v)
			multipleViper.Store(v, A)
		}
	}
}

func LoadViperByFilename(filename string) *Viper {
	value, _ := multipleViper.Load(filename)
	if value == nil {
		return nil
	} else {
		return value.(*Viper)
	}
}

func GetEnvConfig(key string) *utils.TypeTransform {
	env := os.Getenv(strings.Replace(strings.ToUpper(key), ".", "_", -1))
	if env != "" {
		return &utils.TypeTransform{Value: env}
	}
	return utils.Transform(C.Get(key))
}
