package ekafka

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go.uber.org/zap"

	"emperror.dev/errors"
	"github.com/cenkalti/backoff/v4"
	"github.com/weblazy/easy/elog"
	"github.com/weblazy/easy/run"

	"github.com/weblazy/easy/retry"

	"github.com/IBM/sarama"
)

type Handler func(ctx context.Context, message *sarama.ConsumerMessage) error

type ConsumerGroup struct {
	config               *Config
	consumerGroupConfig  *ConsumerGroupConfig
	handler              Handler
	backOffConfig        retry.Config
	consumeRetryInterval time.Duration

	cancel context.CancelFunc
	cg     sarama.ConsumerGroup
}

func newConsumerGroup(config *Config, consumerGroupConfig *ConsumerGroupConfig, sc *sarama.Config) (*ConsumerGroup, error) {
	rc := retry.DefaultConfig()
	rc.MaxRetries = consumerGroupConfig.RetryConfig.MaxRetries

	s := &ConsumerGroup{
		config:               config,
		consumerGroupConfig:  consumerGroupConfig,
		consumeRetryInterval: 100 * time.Millisecond,
		backOffConfig:        rc,
	}

	cg, err := s.getConsumerGroup(*sc)
	if err != nil {
		return nil, err
	}

	s.cg = cg

	return s, nil
}

func (s *ConsumerGroup) SetHandler(h Handler) {
	if s.handler != nil {
		elog.ErrorCtx(context.Background(), "handler already set")
		return
	}
	s.handler = h
}

func (s *ConsumerGroup) Start() {
	if s.handler == nil {
		elog.ErrorCtx(context.Background(), "empty handler")
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	s.cancel = cancel

	go func() {
		elog.InfoCtx(ctx, fmt.Sprintf("Subscribed and listening to topics: %v", s.consumerGroupConfig.Topics))
		for {
			elog.InfoCtx(ctx, "Starting loop to consume.")

			// Consume the requested topics
			bo := backoff.WithContext(backoff.NewConstantBackOff(s.consumeRetryInterval), ctx)

			innerErr := retry.RetryWithLog(ctx, func() error {
				elog.InfoCtx(ctx, "ekafkaConsume", zap.Any("consumerGroupConfig", s.consumerGroupConfig))
				return s.cg.Consume(ctx, s.consumerGroupConfig.Topics, s)
				// return s.cg.Consume(ctx, s.consumerGroupConfig.Topics, otelsarama.WrapConsumerGroupHandler(s))
			}, bo, "ekafka consumeRetry")

			if innerErr != nil && !errors.Is(innerErr, context.Canceled) {
				elog.ErrorCtx(ctx, fmt.Sprintf("Permanent error consuming %v", s.consumerGroupConfig.Topics), elog.FieldError(innerErr))
			}

			// If the context was canceled, as is the case when handling SIGINT and SIGTERM below, then this pops
			// us out of the consumer loop
			if ctx.Err() != nil {
				return
			}
		}
	}()
}

func (s *ConsumerGroup) Close() error {
	if s.cg != nil {
		elog.InfoCtx(context.Background(), "consumer group close")
		s.cancel()
		return s.cg.Close()
	}

	return nil
}

func (s *ConsumerGroup) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (s *ConsumerGroup) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (s *ConsumerGroup) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	if s.handler == nil {
		return errors.New("nil handler")
	}

	for message := range claim.Messages() {
		start := time.Now()

		// ctx := otel.GetTextMapPropagator().Extract(context.Background(), otelsarama.NewConsumerMessageCarrier(message))
		ctx := context.Background()
		b := s.backOffConfig.NewBackOffWithContext(session.Context())

		labels := make([]zap.Field, 0)
		labels = append(labels, zap.String("topic", message.Topic))

		// if tid := etrace.ExtractTraceID(ctx); tid != "" {
		// 	labels = append(labels, elog.FieldTrace(tid))
		// }

		if s.config.EnableAccessInterceptorRes {
			labels = append(labels, zap.Any("res", ConsumerMessageToMap(message)))
		}

		err := retry.RetryWithLog(ctx, func() error {
			return run.RunSafeWrap(ctx, func() error {
				return s.handler(ctx, message)
			})
		}, b, fmt.Sprintf("ekafka message handler retry, topic %s partition %d offset %d", message.Topic, message.Partition, message.Offset))

		duration := time.Since(start)
		labels = append(labels, elog.FieldCost(duration))

		if err != nil {
			labels = append(labels, elog.FieldError(err))
			elog.ErrorCtx(ctx, "kafka handler message error", labels...)
			kafkaConsumerGroupCounter.WithLabelValues(s.config.brokers(), s.consumerGroupConfig.GroupID, message.Topic, CodeError).Inc()
		} else {
			elog.InfoCtx(ctx, "kafka handler message success", labels...)
			kafkaConsumerGroupCounter.WithLabelValues(s.config.brokers(), s.consumerGroupConfig.GroupID, message.Topic, CodeOK).Inc()
		}

		kafkaConsumerGroupHistogram.WithLabelValues(s.config.brokers(), s.consumerGroupConfig.GroupID, message.Topic).Observe(duration.Seconds())

		session.MarkMessage(message, "")
	}

	return nil
}

func (s *ConsumerGroup) getConsumerGroup(sc sarama.Config) (sarama.ConsumerGroup, error) { //nolint
	initialOffset, err := parseInitialOffset(s.consumerGroupConfig.InitialOffset)
	if err != nil {
		return nil, err
	}

	sc.Consumer.Offsets.Initial = initialOffset

	cg, err := sarama.NewConsumerGroup(s.config.ClientConfig.Brokers, s.consumerGroupConfig.GroupID, &sc)
	if err != nil {
		return nil, err
	}

	return cg, nil
}

func parseInitialOffset(value string) (initialOffset int64, err error) {
	switch {
	case strings.EqualFold(value, offsetOldest):
		initialOffset = sarama.OffsetOldest
	case strings.EqualFold(value, offsetNewest):
		initialOffset = sarama.OffsetNewest
	case value != "":
		return 0, fmt.Errorf("kafka error: invalid initialOffset: %s", value)
	default:
		initialOffset = sarama.OffsetOldest // Default
	}

	return initialOffset, err
}
