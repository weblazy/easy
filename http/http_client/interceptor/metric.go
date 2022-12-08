package interceptor

import (
	"net/url"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/prometheus/client_golang/prometheus"
)

var (

	// ClientHandleCounter ...
	ClientHandleCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "",
			Name:      "http_client_handle_total",
		}, []string{"name", "method", "path", "peer", "code"})

	// ClientHandleHistogram ...
	ClientHandleHistogram = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "",
		Name:      "http_client_handle_seconds",
	}, []string{"name", "method", "path", "peer"})
)

func init() {
	prometheus.MustRegister(ClientHandleCounter)
	prometheus.MustRegister(ClientHandleHistogram)
}

type MetricPathRewriter func(origin string) string

func NoopMetricPathRewriter(origin string) string {
	return origin
}

func MetricInterceptor(name, addr string, r MetricPathRewriter) (resty.RequestMiddleware, resty.ResponseMiddleware, resty.ErrorHook) {
	rewriter := NoopMetricPathRewriter
	if r != nil {
		rewriter = r
	}

	afterFn := func(cli *resty.Client, res *resty.Response) error {
		method := res.Request.Method
		path := rewriter(res.Request.RawRequest.URL.Path)
		ClientHandleCounter.WithLabelValues(name, method, path, addr, strconv.Itoa(res.StatusCode())).Inc()
		ClientHandleHistogram.WithLabelValues(name, method, path, addr).Observe(res.Time().Seconds())
		return nil
	}

	errorFn := func(req *resty.Request, err error) {
		method := req.Method
		var path string

		// OnBeforeRequest 有错误时, 拿不到 req.RawRequest
		u, err2 := url.Parse(req.URL)
		if err2 != nil {
			path = "invalidUrl"
		} else {
			path = rewriter(u.Path)
		}

		if v, ok := err.(*resty.ResponseError); ok {
			ClientHandleCounter.WithLabelValues(name, method, path, addr, strconv.Itoa(v.Response.StatusCode())).Inc()
		} else {
			ClientHandleCounter.WithLabelValues(name, method, path, addr, "unknown").Inc()
		}

		ClientHandleHistogram.WithLabelValues(name, method, path, addr).Observe(time.Since(GetStartTime(req.Context())).Seconds())
	}

	return nil, afterFn, errorFn
}
