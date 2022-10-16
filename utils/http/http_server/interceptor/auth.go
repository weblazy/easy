package interceptor

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"runtime"

	"io/ioutil"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sunmi-OS/gocore/v2/utils"
	"github.com/weblazy/crypto/aes"
	"github.com/weblazy/easy/utils/code_err"
	"github.com/weblazy/easy/utils/glog"
	"go.uber.org/zap"
)

const (
	UidHeader       = "X-Uid"
	DebugHeader     = "X-Debug"
	TokenHeader     = "X-Token"
	NonceHeader     = "X-Nonce"
	TimestampHeader = "X-Timestamp"
	SignHeader      = "X-Sign"
	LanguageHeader  = "X-Language"
	TokenPrefix     = "token#"
	UserPrefix      = "user#"
	CodeExpiryTime  = 15 * 60                 //15分钟
	TokenExpiryTime = 7 * 86400 * time.Second //7天
)

// Auth
func Auth(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			GetCurrentGoroutineStack(err)
			Error(c, code_err.ParamsErr, fmt.Errorf("panic"))
		}
	}()
	req := c.Request
	header := req.Header
	debugKey := header.Get(DebugHeader)
	var bodyBytes []byte
	bodyBytes, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		Error(c, code_err.ParamsErr, fmt.Errorf("Invalid request body"))
		return
	}

	if utils.GetRunTime() == "onl" || debugKey != "test" {
		token := header.Get(TokenHeader)
		if token == "" {
			timestamp := header.Get(TimestampHeader)
			nonce := header.Get(NonceHeader)
			token = nonce + timestamp
		}
		encryptKey := Sha256([]byte(token))
		aesObj := aes.NewAes(encryptKey)
		requestBody, err := aesObj.Decrypt(string(bodyBytes))
		if err != nil {
			Error(c, code_err.DecryptErr, err)
			return
		}
		bodyBytes = []byte(requestBody)
	}

	// 新建缓冲区并替换原有Request.body
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(bodyBytes)))
	c.Next()
}

func Sha256ToHex(text []byte) string {
	hash := sha256.New()
	hash.Write(text)
	return hex.EncodeToString(hash.Sum(nil))
}

func Md5(s string) string {
	sum := md5.Sum([]byte(s))
	return hex.EncodeToString(sum[:])
}

func Sha256(text []byte) []byte {
	hash := sha256.New()
	hash.Write(text)
	return hash.Sum(nil)
}

// ValidateSign
func ValidateSign(signParam string, appKey string, bodyParams []byte) (err error) {
	realSign, err := HmacSHA256Sign([]byte(appKey), bodyParams)
	if err != nil {
		return fmt.Errorf("签名错误")
	}
	if string(realSign) != signParam {
		glog.InfoCtx(context.Background(), "signErr", zap.String("realSign:", realSign), zap.String("signParam:", signParam))
		return fmt.Errorf("签名错误: correctSign：%s Token：%s body：%s", realSign, appKey, bodyParams)
	}
	return nil
}

// sign
func HmacSHA256Sign(secret []byte, params []byte) (string, error) {
	mac := hmac.New(sha256.New, secret)
	_, err := mac.Write(params)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(mac.Sum(nil)), nil
}

// GetCurrentGoroutineStack 获取当前Goroutine的调用栈，便于排查panic异常
func GetCurrentGoroutineStack(err interface{}) []byte {
	var buf [4096]byte
	n := runtime.Stack(buf[:], false)
	return buf[:n]
}
