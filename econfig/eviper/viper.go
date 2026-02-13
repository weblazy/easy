package eviper

import (
	"bytes"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"github.com/BurntSushi/toml"
	"github.com/spf13/viper"
	"github.com/weblazy/crypto/aes"
	"github.com/weblazy/crypto/mode"
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

	// 处理加密值
	processEncryptedValues(v)

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

// aesDecrypt AES解密函数
func aesDecrypt(encryptedText string, key []byte) (string, error) {
	// 使用 weblazy/crypto 库进行解密
	aesInstance := aes.NewAes(key).WithMode(&mode.ECBMode{})
	decrypted, err := aesInstance.Decrypt(encryptedText)
	if err != nil {
		return "", err
	}
	return decrypted, nil
}

// processEncryptedValues 处理配置中的加密值
func processEncryptedValues(v *viper.Viper) {
	// 从环境变量获取AES密钥
	aesKey := os.Getenv("AES_KEY")
	if aesKey == "" {
		// 如果没有设置密钥，跳过解密
		return
	}

	// 编译正则表达式匹配 ${xxx} 格式
	re := regexp.MustCompile(`\$\{([^}]+)\}`)

	// 获取所有配置项
	settings := v.AllSettings()

	// 递归处理配置值
	processMap(settings, re, []byte(aesKey))

	// 重新设置配置
	for key, value := range flattenMap("", settings) {
		v.Set(key, value)
	}
}

// processMap 递归处理map中的加密值
func processMap(m map[string]interface{}, re *regexp.Regexp, key []byte) {
	for k, v := range m {
		switch value := v.(type) {
		case string:
			// 检查是否包含 ${xxx} 格式
			if re.MatchString(value) {
				// 替换所有 ${xxx} 为解密后的值
				decryptedValue := re.ReplaceAllStringFunc(value, func(match string) string {
					// 提取 xxx 部分（去掉 ${ 和 }）
					encryptedText := match[2 : len(match)-1]
					// 尝试解密
					if decrypted, err := aesDecrypt(encryptedText, key); err == nil {
						return decrypted
					} else {
						log.Printf("Failed to decrypt value for key %s: %v", k, err)
						return match // 解密失败则保留原值
					}
				})
				m[k] = decryptedValue
			}
		case map[string]interface{}:
			// 递归处理嵌套的map
			processMap(value, re, key)
		case []interface{}:
			// 处理数组
			processSlice(value, re, key)
		}
	}
}

// processSlice 处理数组中的加密值
func processSlice(s []interface{}, re *regexp.Regexp, key []byte) {
	for i, v := range s {
		switch value := v.(type) {
		case string:
			// 检查是否包含 ${xxx} 格式
			if re.MatchString(value) {
				// 替换所有 ${xxx} 为解密后的值
				decryptedValue := re.ReplaceAllStringFunc(value, func(match string) string {
					// 提取 xxx 部分（去掉 ${ 和 }）
					encryptedText := match[2 : len(match)-1]
					// 尝试解密
					if decrypted, err := aesDecrypt(encryptedText, key); err == nil {
						return decrypted
					} else {
						log.Printf("Failed to decrypt array value: %v", err)
						return match // 解密失败则保留原值
					}
				})
				s[i] = decryptedValue
			}
		case map[string]interface{}:
			processMap(value, re, key)
		case []interface{}:
			processSlice(value, re, key)
		}
	}
}

// flattenMap 将嵌套map扁平化为viper可以使用的格式
func flattenMap(prefix string, m map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range m {
		key := k
		if prefix != "" {
			key = prefix + "." + k
		}

		switch value := v.(type) {
		case map[string]interface{}:
			// 递归处理嵌套map
			nested := flattenMap(key, value)
			for nk, nv := range nested {
				result[nk] = nv
			}
		default:
			result[key] = v
		}
	}
	return result
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
