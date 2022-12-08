package interceptor

import (
	"context"
	"time"

	"github.com/go-resty/resty/v2"
)

// https://stackoverflow.com/questions/40891345/fix-should-not-use-basic-type-string-as-key-in-context-withvalue-golint
// https://blog.golang.org/context#TOC_3.2.
// https://golang.org/pkg/context/#WithValue ，这边文章说明了用struct，可以避免分配
type startTimeKey struct{}

func SetStartTimeInterceptor() (resty.RequestMiddleware, resty.ResponseMiddleware, resty.ErrorHook) {
	return func(cli *resty.Client, req *resty.Request) error {
		req.SetContext(context.WithValue(req.Context(), startTimeKey{}, time.Now()))
		return nil
	}, nil, nil
}

func GetStartTime(ctx context.Context) time.Time {
	startTime, _ := ctx.Value(startTimeKey{}).(time.Time)
	return startTime
}
