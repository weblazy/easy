package alipay

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/weblazy/easy/utils/httpx"

	"github.com/weblazy/gocore/httplib"
	"gopkg.in/redis.v5"
)

type (
	Alipay struct {
		appId              string
		rsaPrivateKey      string
		alipayrsaPublicKey string
		signType           string
		format             string
		postCharset        string
		apiVersion         string
		getRedis           func() *redis.Client
	}

	GetUserInfoResponse struct {
		Scene     string `json:"scene"`
		Page      string `json:"page"`
		AutoColor bool   `json:"auto_color"`
		IsHyaline bool   `json:"is_hyaline"`
		Width     int64  `json:"width"`
	}
)

const (
	//网关
	gatewayUrl = "https://openapi.alipay.com/gateway.do"
	TimeLayout = "2006-01-02 15:04:05"
)

// @desc 初始化
// @auth liuguoqiang 2020-02-25
// @param
// @return
func NewAlipay(appId, rsaPrivateKey, alipayrsaPublicKey, signType string, getRedis func() *redis.Client) *Alipay {
	return &Alipay{
		appId:              appId,
		rsaPrivateKey:      rsaPrivateKey,
		alipayrsaPublicKey: alipayrsaPublicKey,
		signType:           signType,
		format:             "json",
		postCharset:        "UTF-8",
		apiVersion:         "1.0",
		getRedis:           getRedis,
	}
}

// @desc
// @auth liuguoqiang 2020-04-09
// @param
// @return
func (s *Alipay) GetUserInfo(code string) {

}

// @desc
// @auth liuguoqiang 2020-04-09
// @param
// @return
func (s *Alipay) Request(url, apiParams map[string]interface{}, authToken string, appInfoAuthtoken string) error {
	//组装系统参数
	params := make(map[string]interface{})
	params["app_id"] = s.appId
	params["version"] = s.apiVersion
	params["format"] = s.format
	params["sign_type"] = s.signType
	params["method"] = url
	params["timestamp"] = time.Now().Format(TimeLayout)
	params["auth_token"] = authToken
	params["alipay_sdk"] = "alipay-sdk-php-20180705"
	params["terminal_type"] = nil
	params["terminal_info"] = nil
	params["prod_code"] = nil
	params["notify_url"] = nil
	params["charset"] = s.postCharset
	params["app_auth_token"] = appInfoAuthtoken
	//签名
	params["sign"] = s.generateSign()
	//系统参数放入GET请求串

	requestUrl, err := httpx.MapToQuery(params)
	if err != nil {
		return err
	}
	requestUrl = gatewayUrl + "?" + requestUrl
	req := httplib.Post(requestUrl)
	req, err = req.JSONBody(apiParams)
	if err != nil {
		return err
	}
	dataByte, err := req.Bytes()
	if err != nil {
		return err
	}
	data := make(map[string]interface{})
	err = json.Unmarshal(dataByte, &data)
	if err != nil {
		return err
	}
	fmt.Printf("%#v", data)
	return nil
}

// @desc
// @auth liuguoqiang 2020-04-09
// @param
// @return
func (s *Alipay) generateSign() string {
	return ""
}
