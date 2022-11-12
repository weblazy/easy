package ftrace

import (
	"context"

	"go.opentelemetry.io/otel/baggage"
)

// GetBaggageValue get baggage info from context, if key not exists, return "", false.
func GetBaggageValue(ctx context.Context, key string) (string, bool) {
	b := baggage.FromContext(ctx)
	m := b.Member(key)

	if m.Key() == "" {
		return "", false
	}

	return m.Value(), true
}

// WithBaggage append baggage by string key val.
func WithBaggage(parent context.Context, key, val string) (context.Context, error) {
	member, err := baggage.NewMember(key, val)
	if err != nil {
		return parent, err
	}

	b := baggage.FromContext(parent)
	b, err = b.SetMember(member)
	if err != nil {
		return parent, err
	}

	return baggage.ContextWithBaggage(parent, b), nil
}

// AppendBaggageByMap append map kvs to current ctx baggage, will return origin ctx if error.
func AppendBaggageByMap(ctx context.Context, mp map[string]string) (context.Context, error) {
	b := baggage.FromContext(ctx)

	for k, v := range mp {
		m, err := baggage.NewMember(k, v)
		if err != nil {
			return ctx, err
		}

		b, err = b.SetMember(m)
		if err != nil {
			return ctx, err
		}
	}

	return baggage.ContextWithBaggage(ctx, b), nil
}
