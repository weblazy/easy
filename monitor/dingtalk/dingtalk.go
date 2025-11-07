package dingtalk

import (
	"context"
	"fmt"

	"github.com/weblazy/easy/http/http_client"
	"github.com/weblazy/easy/http/http_client/http_client_config"
	"github.com/weblazy/easy/monitor"
)

type (
	DingTalk struct {
		monitor.Handler
		Url       string `json:"url"`
		atMobiles []string
		isAtAll   bool
	}

	TextMsg struct {
		Msgtype string `json:"msgtype"`
		Text    struct {
			Content string `json:"content"`
		} `json:"text"`
		At struct {
			AtMobiles []string `json:"atMobiles"`
			IsAtAll   bool     `json:"isAtAll"`
		} `json:"at"`
	}
)

// @desc
// @auth liuguoqiang 2020-12-07
// @param
// @return
func NewDingTalk(url string) *DingTalk {
	return &DingTalk{
		Url:       url,
		atMobiles: []string{},
		isAtAll:   false,
	}
}

// @desc @部分成员
// @auth liuguoqiang 2020-12-07
// @param
// @return
func (dingTalk *DingTalk) WithAtMobiles(atMobiles []string) *DingTalk {
	if atMobiles != nil {
		dingTalk.atMobiles = atMobiles
	}
	return dingTalk
}

// @desc @所有成员
// @auth liuguoqiang 2020-12-07
// @param
// @return
func (dingTalk *DingTalk) WithIsAtAll(isAtAll bool) *DingTalk {
	dingTalk.isAtAll = isAtAll
	return dingTalk
}

// @desc 发送钉钉消息
// @auth liuguoqiang 2020-12-07
// @param
// @return
func (dingTalk *DingTalk) SendMsg(body interface{}) ([]byte, error) {
	cfg := http_client_config.DefaultConfig()
	client := http_client.NewHttpClient(cfg)
	request := client.Request.SetContext(context.Background()).SetBody(body)
	resp, err := request.Post(dingTalk.Url)
	return resp.Body(), err
}

// @desc 发送钉钉文本消息
// @auth liuguoqiang 2020-12-07
// @param
// @return
func (dingTalk *DingTalk) SendTextMsg(ctx context.Context, content string) error {
	if dingTalk.Url == "" {
		return fmt.Errorf("报警地址为空")
	}
	msg := TextMsg{
		Msgtype: "text",
	}
	msg.Text.Content = content
	msg.At.IsAtAll = dingTalk.isAtAll
	msg.At.AtMobiles = dingTalk.atMobiles
	_, err := dingTalk.SendMsg(msg)
	return err
}

// // @desc Request 通用请求
// // @auth liuguoqiang 2020-12-07
// // @param
// // @return
// func Request(url string, body interface{}, headers map[string]string) ([]byte, error) {
// 	client := http_request.New()
// 	req := client.Request
// 	if headers != nil {
// 		req = req.SetHeaders(headers)
// 	}
// 	response, err := req.
// 		SetBody(body).
// 		Post(url)
// 	if err != nil {
// 		return nil, err
// 	}
// 	respByte := response.Body()
// 	return respByte, err
// }
