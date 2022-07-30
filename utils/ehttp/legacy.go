package httpx

import "github.com/go-resty/resty/v2"

// AddMetricInterceptor 给 sdk 中的裸 resty 老代码使用
func AddMetricInterceptor(client *resty.Client, name, addr string, rewriter MetricPathRewriter) {
	// context start time
	fb, fa, fe := fixedInterceptor(name, nil)
	AddInterceptors(client, fb, fa, fe)
	mb, ma, me := MetricInterceptor(name, addr, rewriter)
	AddInterceptors(client, mb, ma, me)
}
