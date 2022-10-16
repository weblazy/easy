package interceptor

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
)

// https://stackoverflow.com/questions/40891345/fix-should-not-use-basic-type-string-as-key-in-context-withvalue-golint
// https://blog.golang.org/context#TOC_3.2.
// https://golang.org/pkg/context/#WithValue ，这边文章说明了用struct，可以避免分配
type startTimeKey struct{}

func SetStartTimeInterceptor() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request.WithContext(context.WithValue(c.Request.Context(), startTimeKey{}, time.Now()))
	}
}

func GetStartTime(ctx context.Context) time.Time {
	startTime, _ := ctx.Value(startTimeKey{}).(time.Time)
	return startTime
}
