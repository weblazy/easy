package http_client

import (
	"net"
	"net/http"
	"net/http/cookiejar"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/weblazy/easy/utils/http/http_client/http_client_config"
	"github.com/weblazy/easy/utils/http/http_client/interceptor"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"golang.org/x/net/publicsuffix"
)

type HttpClient struct {
	config  *http_client_config.Config
	Client  *resty.Client
	Request *resty.Request
}

func NewHttpClient(c *http_client_config.Config) *HttpClient {
	if c == nil {
		c = http_client_config.DefaultConfig()
	}
	// resty的默认方法，无法设置长连接个数，和是否开启长连接，这里重新构造http client。
	cookieJar, _ := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List}) // nolint

	client := resty.NewWithClient(&http.Client{Transport: createTransport(c), Jar: cookieJar}).
		SetDebug(c.RawDebug).
		SetTimeout(c.ReadTimeout).
		SetBaseURL(c.Addr)
		// 中间件
	onBefore, onAfter, onErr := interceptor.HeaderCarrierInterceptor()
	AddInterceptors(client, onBefore, onAfter, onErr)

	onBefore, onAfter, onErr = interceptor.SetStartTimeInterceptor()
	AddInterceptors(client, onBefore, onAfter, onErr)

	onBefore, onAfter, onErr = interceptor.LogInterceptor(c)
	AddInterceptors(client, onBefore, onAfter, onErr)

	if c.EnableMetricInterceptor {
		onBefore, onAfter, onErr := interceptor.MetricInterceptor(c.Name, c.Addr, nil)
		AddInterceptors(client, onBefore, onAfter, onErr)
	}

	return &HttpClient{
		Client:  client,
		Request: client.R(),
	}
}

func createTransport(c *http_client_config.Config) http.RoundTripper {
	dialer := &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
		DualStack: true,
	}

	t := &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		DialContext:           dialer.DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          c.MaxIdleConns,
		IdleConnTimeout:       c.IdleConnTimeout,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		DisableKeepAlives:     !c.EnableKeepAlives,
		MaxIdleConnsPerHost:   c.MaxIdleConnsPerHost,
	}

	if c.EnableTraceInterceptor {
		return otelhttp.NewTransport(t)
	}

	return t
}

func (h *HttpClient) EnableMetricInterceptor(metricPathRewriter interceptor.MetricPathRewriter) {
	onBefore, onAfter, onErr := interceptor.MetricInterceptor(h.config.Name, h.config.Addr, metricPathRewriter)
	AddInterceptors(h.Client, onBefore, onAfter, onErr)
}

func AddInterceptors(client *resty.Client, onBefore resty.RequestMiddleware, onAfter resty.ResponseMiddleware, onErr resty.ErrorHook) {
	if onBefore != nil {
		client.OnBeforeRequest(onBefore)
	}
	if onAfter != nil {
		client.OnAfterResponse(onAfter)
	}
	if onErr != nil {
		client.OnError(onErr)
	}
}

func (h *HttpClient) SetTrace(header interface{}) *HttpClient {
	trace := SetHeader(header)
	h.Request.Header = trace.HttpHeader
	return h
}
