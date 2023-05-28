package eviper

import (
	"bytes"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"

	"github.com/BurntSushi/toml"
	"github.com/spf13/viper"
)

type Viper struct {
	*viper.Viper
}

var multipleViper sync.Map

func NewViperFromString(configs string) *Viper {
	v := viper.New()
	CheckToml(configs)
	v.SetConfigType("toml")
	err := v.ReadConfig(bytes.NewBuffer([]byte(configs)))
	if err != nil {
		print(err)
	}
	return &Viper{v}
}

func (v *Viper) MergeViperFromString(configs string) {
	CheckToml(configs)
	v.SetConfigType("toml")
	err := v.MergeConfig(bytes.NewBuffer([]byte(configs)))
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

func NewViperFromFile(filePath string, fileName string) *Viper {
	return newConfig(filePath, fileName)
}

func newConfig(filePath string, fileName string) *Viper {
	v := viper.New()
	v.SetConfigName(fileName)
	//filePath支持相对路径和绝对路径 etc:"/a/b" "b" "./b"
	if filePath == "" || filePath[:1] != "/" {
		v.AddConfigPath(path.Join(GetPath(), filePath))
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

// GetPath 获取项目路径
func GetPath() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		print(err.Error())
	}
	path := strings.Replace(dir, "\\", "/", -1)
	return path
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

func (v *Viper) GetEnvConfig(key string) string {
	env := os.Getenv(strings.Replace(strings.ToUpper(key), ".", "_", -1))
	if env != "" {
		return env
	}
	return v.GetString(key)
}
