package interceptor

import (
	"bytes"
	"fmt"
	"io/ioutil"

	"github.com/gin-gonic/gin"
	"github.com/weblazy/easy/code_err"
	"github.com/weblazy/easy/econfig"
)

// Token
func Token(c *gin.Context) {
	debugKey := c.Request.Header.Get(DebugHeader)
	if !econfig.GlobalViper.GetBool("BaseConfig.Debug") || debugKey != econfig.GlobalViper.GetString("BaseConfig.XDebugKey") {
		token := c.Request.Header.Get(TokenHeader)
		if token == "" {
			Error(c, code_err.TokenErr, fmt.Errorf("token 不存在"))
			return
		}
	}
	c.Next()
}

// Sign
func Sign(userIdHeader string, validateToken func(token string) (uid string, err error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		req := c.Request
		header := req.Header
		debugKey := header.Get(DebugHeader)
		var bodyBytes []byte
		bodyBytes, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			Error(c, code_err.ParamsErr, fmt.Errorf("Invalid request body"))
			return
		}
		// 新建缓冲区并替换原有Request.body
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(bodyBytes)))
		if !econfig.GlobalViper.GetBool("BaseConfig.Debug") || debugKey != econfig.GlobalViper.GetString("BaseConfig.XDebugKey") {
			sign := header.Get(SignHeader)
			token := header.Get(TokenHeader)
			timestamp := header.Get(TimestampHeader)
			nonce := header.Get(NonceHeader)
			if token == "" {
				token = nonce + timestamp
			} else {
				uid, err := validateToken(token)
				if err != nil {
					Error(c, code_err.TokenErr, err)
					return
				}
				header.Set(userIdHeader, uid)
			}
			err = ValidateSign(sign, token, []byte(string(bodyBytes)+timestamp+nonce))
			if err != nil {
				Error(c, code_err.SignErr, err)
				return
			}
		}
	}
}
