package larkbot

import (
	"context"

	"github.com/weblazy/easy/elog"
	"github.com/weblazy/easy/http/http_client"
	"github.com/weblazy/easy/http/http_client/http_client_config"
	"github.com/weblazy/easy/monitor"
	"go.uber.org/zap"
)

type Larkbot struct {
	monitor.Handler
	Url    string `json:"url"`
	Header string `json:"header"`
}

type Message struct {
	MsgType string       `json:"msg_type,omitempty"`
	Content interface{}  `json:"content,omitempty"`
	Card    *MessageCard `json:"card,omitempty"`
}

// card "github.com/larksuite/oapi-sdk-go/v3/card"
type MessageCardPlainText struct {
	Tag     string                    `json:"tag,omitempty"`
	Content string                    `json:"content,omitempty"`
	Lines   int                       `json:"lines,omitempty"`
	I18n    *MessageCardPlainTextI18n `json:"i18n,omitempty"`
}
type MessageCardPlainTextI18n struct {
	ZhCN string `json:"zh_cn,omitempty"`
	EnUS string `json:"en_us,omitempty"`
	JaJP string `json:"ja_jp,omitempty"`
}

type MessageCardHeader struct {
	Template string                `json:"template,omitempty"`
	Title    *MessageCardPlainText `json:"title,omitempty"`
}

type MessageCardDiv struct {
	Tag    string              `json:"tag,omitempty"`
	Text   *MessageCardText    `json:"text,omitempty"`
	Fields []*MessageCardField `json:"fields,omitempty"`
	Extra  interface{}         `json:"extra,omitempty"`
}

type MessageCardField struct {
	IsShort bool             `json:"is_short,omitempty"`
	Text    *MessageCardText `json:"text,omitempty"`
}

type MessageCardImage struct {
	Alt          *MessageCardPlainText  `json:"alt,omitempty"`
	Title        *MessageCardText       `json:"title,omitempty"`
	ImgKey       string                 `json:"img_key,omitempty"`
	CustomWidth  int                    `json:"custom_width,omitempty"`
	CompactWidth bool                   `json:"compact_width,omitempty"`
	Mode         *MessageCardImageModel `json:"mode,omitempty"`
	Preview      bool                   `json:"preview,omitempty"`
}

type MessageCardImageModel string

const (
	MessageCardImageModelFitHorizontal MessageCardImageModel = "fit_horizontal"
	MessageCardImageModelCropCenter    MessageCardImageModel = "crop_center"
)

type MessageCardConfig struct {
	EnableForward  bool `json:"enable_forward,omitempty"`
	UpdateMulti    bool `json:"update_multi,omitempty"`
	WideScreenMode bool `json:"wide_screen_mode,omitempty"`
}

const (
	TemplateBlue      = "blue"
	TemplateWathet    = "wathet"
	TemplateTurquoise = "turquoise"
	TemplateGreen     = "green"
	TemplateYellow    = "yellow"
	TemplateOrange    = "orange"
	TemplateRed       = "red"
	TemplateCarmine   = "carmine"
	TemplateViolet    = "violet"
	TemplatePurple    = "purple"
	TemplateIndigo    = "indigo"
	TemplateGrey      = "grey"
)

type MessageCardI18nElements struct {
	ZhCN []string `json:"zh_cn,omitempty"`
	EnUS []string `json:"en_us,omitempty"`
	JaJP []string `json:"ja_jp,omitempty"`
}

type MessageCard struct {
	Config       *MessageCardConfig       `json:"config,omitempty"`
	Header       *MessageCardHeader       `json:"header,omitempty"`
	Elements     []interface{}            `json:"elements,omitempty"`
	I18nElements *MessageCardI18nElements `json:"i18n_elements,omitempty"`
	CardLink     *MessageCardURL          `json:"card_link,omitempty"`
}

type MessageCardText struct {
	Tag     string `json:"tag,omitempty"`
	Content string `json:"content,omitempty"`
}
type MessageCardURL struct {
	URL        string `json:"url,omitempty"`
	AndroidURL string `json:"android_url,omitempty"`
	IOSURL     string `json:"ios_url,omitempty"`
	PCURL      string `json:"pc_url,omitempty"`
}

func (l *Larkbot) SendTextMsg(content string) {
	ctx := context.Background()
	// 卡片消息体
	messageCard := MessageCard{
		Config: &MessageCardConfig{
			WideScreenMode: true,
		},
		Header: &MessageCardHeader{
			Template: "turquoise",
			Title: &MessageCardPlainText{
				Tag:     "plain_text",
				Content: l.Header,
			},
		},
		Elements: []interface{}{
			MessageCardDiv{
				Tag: "div",
				Text: &MessageCardText{
					Tag:     "plain_text",
					Content: content,
				},
			},
		},
	}

	request := http_client.NewHttpClient(http_client_config.DefaultConfig()).
		Request.SetContext(ctx).
		SetBody(&Message{
			MsgType: "interactive",
			Card:    &messageCard,
		})
	resp, err := request.Post(l.Url)
	if err != nil {
		elog.ErrorCtx(ctx, "SendLarkMsg", zap.String("resp", string(resp.Body())), zap.Error(err))

	}
	elog.ErrorCtx(ctx, "SendLarkMsg", zap.String("resp", string(resp.Body())))
}

func (l *Larkbot) SendCardMsg(fields []*MessageCardField) {
	ctx := context.Background()
	// 卡片消息体
	messageCard := MessageCard{
		Config: &MessageCardConfig{
			WideScreenMode: true,
		},
		Header: &MessageCardHeader{
			Template: "turquoise",
			Title: &MessageCardPlainText{
				Tag:     "plain_text",
				Content: l.Header,
			},
		},
		Elements: []interface{}{
			MessageCardDiv{
				Tag:    "div",
				Fields: fields,
			},
		},
	}

	request := http_client.NewHttpClient(http_client_config.DefaultConfig()).
		Request.SetContext(ctx).
		SetBody(&Message{
			MsgType: "interactive",
			Card:    &messageCard,
		})
	resp, err := request.Post(l.Url)
	if err != nil {
		elog.ErrorCtx(ctx, "SendLarkMsg", zap.String("resp", string(resp.Body())), zap.Error(err))

	}
	elog.ErrorCtx(ctx, "SendLarkMsg", zap.String("resp", string(resp.Body())))
}
