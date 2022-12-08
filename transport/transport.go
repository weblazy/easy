package transport

import (
	"context"
	"strings"

	"github.com/weblazy/easy/set"

	"go.opentelemetry.io/otel/propagation"
)

const (
	PrefixPass = "x-pass-"

	lenPP = len(PrefixPass)
)

type ctxKeyType struct{}

var ctxKey ctxKeyType

var customKeyStore = contextKeyStore{
	keyArr: make([]string, 0),
	keySet: set.NewStringSet(),
}

type contextKeyStore struct {
	keyArr []string
	keySet *set.StringSet
}

// RegisterCustomKeys 注册全局 custom keys
// key 必须为纯小写否则 panic
// 只应该在程序启动前全局注册一次
// 是以 append 的形式, 支持多次调用
// 为 lib 库预留扩展空间
func RegisterCustomKeys(keys []string) {
	// 去重
	for _, k := range keys {
		kk := strings.ToLower(k)
		if k != kk {
			panic("custom key only support lowercase")
		}
		customKeyStore.keySet.Add(k)
	}
	customKeyStore.keyArr = customKeyStore.keySet.ToArray()
}

// for test only
func reset() {
	customKeyStore = contextKeyStore{
		keyArr: make([]string, 0),
		keySet: set.NewStringSet(),
	}
}

func getMap(ctx context.Context) map[string]string {
	if ctx != nil {
		if val, ok := ctx.Value(ctxKey).(map[string]string); ok {
			return val
		}
	}

	return make(map[string]string, 0)
}

func setMap(ctx context.Context, m map[string]string) context.Context {
	if ctx == nil {
		return nil
	}

	return context.WithValue(ctx, ctxKey, m)
}

func GetMapFromContext(ctx context.Context) map[string]string {
	mp := getMap(ctx)

	if len(mp) > 0 {
		m := make(map[string]string, len(mp))
		for k, v := range mp {
			m[k] = v
		}
		return m
	}

	return mp
}

func GetMapFromPropagator(carrier propagation.TextMapCarrier) map[string]string {
	mp := make(map[string]string)
	for _, k := range carrier.Keys() {
		kk := strings.ToLower(k)

		if customKeyStore.keySet.Has(kk) || (len(kk) > lenPP && strings.HasPrefix(kk, PrefixPass)) {
			v := carrier.Get(kk)
			if v != "" {
				mp[kk] = v
			}
		}
	}
	return mp
}

type CustomKeys struct{}

var CustomKeysMapPropagator propagation.TextMapPropagator = &CustomKeys{}

func (c *CustomKeys) Inject(ctx context.Context, carrier propagation.TextMapCarrier) {
	mp := getMap(ctx)
	for k, v := range mp {
		if v != "" {
			carrier.Set(k, v)
		}
	}
}

func (c *CustomKeys) Extract(ctx context.Context, carrier propagation.TextMapCarrier) context.Context {
	mp := getMap(ctx)

	cmp := GetMapFromPropagator(carrier)
	if len(cmp) == 0 {
		return ctx
	}

	for k, v := range cmp {
		mp[k] = v
	}

	return setMap(ctx, mp)
}

// Fields not used
func (c *CustomKeys) Fields() []string {
	return nil
}
