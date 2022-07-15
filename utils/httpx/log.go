package httpx

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/spf13/cast"
	"github.com/weblazy/easy/utils/glog"
	"github.com/weblazy/easy/utils/glog/sls"
)

type Log interface {
	Info(obj *LogObject) error
}

type LogObject struct {
	Url             string        `json:"url"`
	Method          string        `json:"method"`
	RequestHders    http.Header   `json:"request_headers"`
	RequestRawBody  interface{}   `json:"request_raw_body"`
	ResponseHeaders http.Header   `json:"response_headers"`
	ResponseBody    string        `json:"response_body"`
	StartTime       string        `json:"start_time"`
	Duration        time.Duration `json:"duration"`
	Status          int           `json:"status"`
}

func (h *HttpClient) SetLog(log Log) *HttpClient {
	err := h.Client.OnAfterResponse(func(client *resty.Client, resp *resty.Response) error {
		r := resp.Request
		obj := &LogObject{
			Url:             r.URL,
			Method:          r.Method,
			RequestHders:    r.Header,
			RequestRawBody:  r.Body,
			ResponseHeaders: resp.Header(),
			ResponseBody:    string(resp.Body()),
			StartTime:       r.Time.Format("2006-01-02 15:04:05"),
			Duration:        resp.Time() / time.Millisecond,
			Status:          resp.StatusCode(),
		}
		log.Info(obj)
		return nil
	})
	if err != nil {
		fmt.Printf("%#v\n", err)
	}
	return h
}

type GocoreLog struct {
}

func NewGocoreLog() *GocoreLog {
	return &GocoreLog{}
}

func (l *GocoreLog) Info(obj *LogObject) error {
	data, _ := json.Marshal(obj)
	glog.InfoF(string(data))
	return nil
}

func NewAliyunLog(topic string) *AliyunLog {
	return &AliyunLog{topic: topic}
}

type AliyunLog struct {
	topic string
}

//  使用阿里云日志需要提前调用sls.InitLog初始化
func (l *AliyunLog) Info(obj *LogObject) error {
	requestHeaderBytes, _ := json.Marshal(obj.RequestHders)
	requestBodyBytes, _ := json.Marshal(obj.RequestRawBody)
	responseHeaderBytes, _ := json.Marshal(obj.ResponseHeaders)
	_ = sls.Info(l.topic, map[string]string{
		"url":              obj.Url,
		"method":           obj.Method,
		"request_headers":  string(requestHeaderBytes),
		"request_raw_body": string(requestBodyBytes),
		"response_headers": string(responseHeaderBytes),
		"start_time":       obj.StartTime,
		"duration":         cast.ToString(obj.Duration),
		"status":           cast.ToString(obj.Status),
	})
	return nil
}
