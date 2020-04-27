package wx

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"strconv"
	"time"

	"github.com/sunmi-OS/gocore/httplib"
	"gopkg.in/redis.v5"
)

type (
	Wx struct {
		appId       string
		secret      string
		grantType   string
		accessToken string
		getRedis    func() *redis.Client
	}

	GetUnLimitQRCodeRequest struct {
		Scene     string `json:"scene"`
		Page      string `json:"page"`
		AutoColor bool   `json:"auto_color"`
		IsHyaline bool   `json:"is_hyaline"`
		Width     int64  `json:"width"`
	}
	SendRequest struct {
		Openid          string `json:"touser"`
		TemplateId      string `json:"template_id"`
		Page            string `json:"page"`
		FormId          string `json:"form_id"`
		Data            string `json:"data"`
		EmphasisKeyword string `json:"emphasis_keyword"`
	}
	DecryptDataRequest struct {
		SessionKey string
		Iv         string
		Data       string
	}
	CheckLoginResponse struct {
		OpenId     string `json:"openId"`
		SessionKey string `json:"session_key"`
	}
)

const (
	// access_token 地址
	AccessTokenUrl = "https://api.weixin.qq.com/cgi-bin/token"
	//获取无限制小程序二维码
	CreateUqrcodeUrl = "https://api.weixin.qq.com/wxa/getwxacodeunlimit"
	//授权$code 访问地址
	CodeAccessUrl    = "https://api.weixin.qq.com/sns/jscode2session"
	TemplatedSendUrl = "https://api.weixin.qq.com/cgi-bin/message/wxopen/template/send"
	//获取授权页ticket
	GetticketUrl = "https://api.weixin.qq.com/cgi-bin/ticket/getticket"
	//获取授权页链接
	AuthUrl = "https://api.weixin.qq.com/card/invoice/getauthurl"
	//上传电子发票PDF文件
	SetpdfUrl = "https://api.weixin.qq.com/card/invoice/platform/setpdf"
	//将电子发票卡券插入用户卡包
	InvoiceInsertUrl = "https://api.weixin.qq.com/card/invoice/insert"
)

// @desc 初始化
// @auth liuguoqiang 2020-02-25
// @param
// @return
func NewWx(appId, secret, grantType string, getRedis func() *redis.Client) *Wx {
	return &Wx{
		appId:     appId,
		secret:    secret,
		grantType: grantType,
		getRedis:  getRedis,
	}
}

// @desc 根据access_token值进行授权
// @auth liuguoqiang 2020-02-25
// @param $isFresh 是否刷新access_token
// @return
func (s *Wx) InitAuthToken(isFresh bool) (string, error) {
	//查询缓存
	tokenKey := "tob_wechat:applet:token:" + s.appId
	accessToken := s.getRedis().Get(tokenKey).Val()

	if accessToken != "" && !isFresh {
		s.accessToken = accessToken
		return s.accessToken, nil
	}

	// 获取token
	req := httplib.Get(AccessTokenUrl + "?grant_type=client_credential&appid=" + s.appId + "&secret=" + s.secret)
	data := make(map[string]interface{})
	err := req.ToJSON(&data)
	if err != nil {
		return "", err
	}
	if accessToken, ok := data["access_token"]; !ok {
		return "", errors.New(strconv.FormatFloat(data["errcode"].(float64), 'f', -1, 64) + ":" + data["errmsg"].(string))
	} else {
		err := s.getRedis().Set(tokenKey, accessToken, 7000*time.Second).Err()
		if err != nil {
			return "", err
		}
		s.accessToken = accessToken.(string)
		return s.accessToken, nil
	}
}

// @desc 获取二维码
// @auth liuguoqiang 2020-02-25
// @param
// @return
func (s *Wx) GetUnLimitQRCode(params *GetUnLimitQRCodeRequest, isFresh bool) ([]byte, error) {
	return s.Request(nil, params, CreateUqrcodeUrl, isFresh, true)
}

// @desc 发送模板
// @auth liuguoqiang 2020-04-17
// @param
// @return
func (s *Wx) Send(params *SendRequest, isFresh bool) ([]byte, error) {
	return s.Request(nil, params, TemplatedSendUrl, isFresh, true)
}

// @desc 获取授权页面ticket
// @auth liuguoqiang 2020-02-25
// @param
// @return
func (s *Wx) GetTicket(ticketType string, isFresh bool) ([]byte, error) {
	urlParam := make(map[string]string)
	urlParam["type"] = "wx_card"
	return s.Request(urlParam, nil, GetticketUrl, isFresh, false)
}

// @desc 获取微信授权页链接
// @auth liuguoqiang 2020-02-25
// @param
// @return
func (s *Wx) GetAuthUrl(params map[string]interface{}, isFresh bool) ([]byte, error) {
	return s.Request(nil, params, AuthUrl, isFresh, true)
}

// @desc 通用请求
// @auth liuguoqiang 2020-02-25
// @param
// @return
func (s *Wx) Request(urlParam map[string]string, bodyParams interface{}, paramUrl string, isFresh bool, isPost bool) ([]byte, error) {
	if s.accessToken == "" || isFresh {
		_, err := s.InitAuthToken(isFresh)
		if err != nil {
			return nil, err
		}
	}
	var req *httplib.BeegoHTTPRequest
	url := paramUrl + "?access_token=" + s.accessToken
	if isPost {
		req = httplib.Post(url)
	} else {
		req = httplib.Get(url)
	}
	if urlParam != nil {
		for key, value := range urlParam {
			req = req.Param(key, value)
		}
	}
	req, err := req.JSONBody(bodyParams)
	if err != nil {
		return nil, err
	}
	dataByte, err := req.Bytes()
	if err != nil {
		return nil, err
	}
	data := make(map[string]interface{})
	err = json.Unmarshal(dataByte, &data)
	if err == nil {
		if _, ok := data["errcode"]; ok && data["errcode"].(float64) != 0 {
			if !isFresh {
				dataByte, err = s.Request(urlParam, bodyParams, paramUrl, true, isPost)
				if err != nil {
					return nil, err
				}
			} else {
				return nil, errors.New(strconv.FormatFloat(data["errcode"].(float64), 'f', -1, 64) + ":" + data["errmsg"].(string))
			}
		}
	}
	return dataByte, nil
}

// @desc 根据微信code获取授权信息
// @auth liuguoqiang 2020-04-08
// @param
// @return
func (s *Wx) CheckLogin(code string) (*CheckLoginResponse, error) {
	params := make(map[string]interface{})
	params["appid"] = s.appId
	params["secret"] = s.secret
	params["js_code"] = code
	params["grant_type"] = s.grantType
	req := httplib.Post(CodeAccessUrl)
	req, err := req.JSONBody(params)
	if err != nil {
		return nil, err
	}
	dataByte, err := req.Bytes()
	if err != nil {
		return nil, err
	}
	data := make(map[string]interface{})
	err = json.Unmarshal(dataByte, &data)
	if err == nil {
		if _, ok := data["errcode"]; ok {
			return nil, errors.New(strconv.FormatFloat(data["errcode"].(float64), 'f', -1, 64) + ":" + data["errmsg"].(string))
		}
	}
	return &CheckLoginResponse{
		OpenId:     data["openid"].(string),
		SessionKey: data["session_key"].(string),
	}, nil
}

// @desc 检验数据的真实性，并且获取解密后的明文.
// @auth liuguoqiang 2020-04-08
// @param
// @return
func DecryptData(req *DecryptDataRequest) ([]byte, error) {
	if len(req.SessionKey) != 24 {
		return nil, fmt.Errorf("错误的SessionKey")
	}
	if len(req.Iv) != 24 {
		return nil, fmt.Errorf("错误的Iv")
	}
	aesKey, err := base64.StdEncoding.DecodeString(req.SessionKey)
	if err != nil {
		return nil, err
	}
	iv, err := base64.StdEncoding.DecodeString(req.Iv)
	if err != nil {
		return nil, err
	}
	data, err := base64.StdEncoding.DecodeString(req.Data)
	if err != nil {
		return nil, err
	}
	resp := AesDecrypt(data, aesKey, iv)
	return resp, nil
}

// @desc aes加密
// @auth liuguoqiang 2020-04-21
// @param
// @return
func AesEncrypt(origData []byte, k []byte, iv []byte) string {
	// 分组秘钥
	// NewCipher该函数限制了输入k的长度必须为16, 24或者32
	block, _ := aes.NewCipher(k)
	// 获取秘钥块的长度
	blockSize := block.BlockSize()
	// 补全码
	origData = PKCS7Padding(origData, blockSize)
	// 加密模式
	blockMode := cipher.NewCBCEncrypter(block, iv)
	// 创建数组
	cryted := make([]byte, len(origData))
	// 加密
	blockMode.CryptBlocks(cryted, origData)
	return base64.StdEncoding.EncodeToString(cryted)
}

// @desc aes解密
// @auth liuguoqiang 2020-04-21
// @param
// @return
func AesDecrypt(crytedByte []byte, key []byte, iv []byte) []byte {
	// 分组秘钥
	block, _ := aes.NewCipher(key)
	// 加密模式
	blockMode := cipher.NewCBCDecrypter(block, iv)
	// 创建数组
	orig := make([]byte, len(crytedByte))
	// 解密
	blockMode.CryptBlocks(orig, crytedByte)
	// 去补全码
	orig = PKCS7UnPadding(orig)
	return orig
}

//补码
//AES加密数据块分组长度必须为128bit(byte[16])，密钥长度可以是128bit(byte[16])、192bit(byte[24])、256bit(byte[32])中的任意一个。
// @desc
// @auth liuguoqiang 2020-04-21
// @param
// @return
func PKCS7Padding(ciphertext []byte, blocksize int) []byte {
	padding := blocksize - len(ciphertext)%blocksize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

// @desc 去码
// @auth liuguoqiang 2020-04-21
// @param
// @return
func PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

// @desc 电子发票插入微信卡包
// @auth liuguoqiang 2020-04-24
// @param
// @return
func (s *Wx) InvoiceInsert(params map[string]interface{}, isFresh bool) ([]byte, error) {
	params["appid"] = s.appId
	return s.Request(nil, params, InvoiceInsertUrl, isFresh, true)
}

// @desc 上传pdf文件
// @auth liuguoqiang 2020-02-25
// @param
// @return
func (s *Wx) SetPdf(pdfPath string, isFresh bool) ([]byte, error) {
	if s.accessToken == "" || isFresh {
		_, err := s.InitAuthToken(isFresh)
		if err != nil {
			return nil, err
		}
	}
	if pdfPath == "" {
		return nil, fmt.Errorf("pdfPath is empty")
	}
	resp, err := http.Get(pdfPath)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, errors.New("resp status:" + fmt.Sprint(resp.StatusCode))
	}
	bin, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	buf := new(bytes.Buffer)
	bodyWriter := multipart.NewWriter(buf)
	fileName := fmt.Sprintf("%d.pdf", time.Now().UnixNano()/1000)
	fileWriter, err := bodyWriter.CreateFormFile("pdf", fileName)
	if err != nil {
		return nil, err
	}
	_, err = fileWriter.Write(bin)
	if err != nil {
		return nil, err
	}
	bodyWriter.Close()
	req1, err := http.NewRequest("POST", SetpdfUrl+"?access_token="+s.accessToken, buf)
	req1.Header.Set("Content-Type", bodyWriter.FormDataContentType())
	if err != nil {
		return nil, err
	}
	client := http.Client{}
	resp1, err := client.Do(req1)
	if err != nil {
		return nil, err
	}
	defer resp1.Body.Close()

	dataByte, err := ioutil.ReadAll(resp1.Body)
	if err != nil {
		return nil, err
	}
	data := make(map[string]interface{})
	err = json.Unmarshal(dataByte, &data)
	if err == nil {
		if _, ok := data["errcode"]; ok && data["errcode"].(float64) != 0 {
			if !isFresh {
				dataByte, err = s.SetPdf(pdfPath, true)
				if err != nil {
					return nil, err
				}
			} else {
				return nil, errors.New(strconv.FormatFloat(data["errcode"].(float64), 'f', -1, 64) + ":" + data["errmsg"].(string))
			}
		}
	}
	return dataByte, nil
}
