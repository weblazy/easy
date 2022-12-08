package fkafka

import (
	"context"

	"github.com/Shopify/sarama"
	"go.opentelemetry.io/contrib/instrumentation/github.com/Shopify/sarama/otelsarama"
	"go.uber.org/zap"

	"github.com/weblazy/easy/elog"

	"github.com/weblazy/easy/etrace"
)

const (
	CodeOK    = "OK"
	CodeError = "Error"
)

type Producer struct {
	config   *Config
	producer sarama.SyncProducer
}

func newProducer(config *Config, sc *sarama.Config) (*Producer, error) { //nolint
	c := &Producer{
		config: config,
	}

	producer, err := getSyncProducer(config, *sc)
	if err != nil {
		return nil, err
	}

	c.producer = producer

	return c, nil
}

func (c *Producer) SendMessage(ctx context.Context, msg *Message) error {
	smsg := &sarama.ProducerMessage{
		Topic:    msg.Topic,
		Key:      sarama.ByteEncoder(msg.Key),
		Value:    sarama.ByteEncoder(msg.Value),
		Headers:  msg.Headers,
		Metadata: msg.Metadata,
	}

	partition, offset, err := c.producer.SendMessage(smsg)

	labels := make([]zap.Field, 10)
	labels = append(labels, zap.String("topic", smsg.Topic))

	if c.config.EnableAccessInterceptorReq {
		labels = append(labels, zap.Any("req", msg.ToMap()))
	}

	if tid := etrace.ExtractTraceID(ctx); tid != "" {
		labels = append(labels, elog.FieldTrace(tid))
	}

	if err != nil {
		labels = append(labels, elog.FieldError(err))
		elog.ErrorCtx(ctx, "kafka publish failed", labels...)
		kafkaPublishCounter.WithLabelValues(c.config.brokers(), msg.Topic, CodeError).Inc()
		return err
	}

	labels = append(labels, zap.Int64("partition", int64(partition)), zap.Int64("offset", offset))

	elog.InfoCtx(ctx, "kafka publish success", labels...)
	kafkaPublishCounter.WithLabelValues(c.config.brokers(), msg.Topic, CodeOK).Inc()

	return nil
}

func (c *Producer) Close() error {
	if c.producer != nil {
		elog.InfoCtx(context.Background(), "producer exit")
		return c.producer.Close()
	}

	return nil
}

func getSyncProducer(config *Config, sc sarama.Config) (sarama.SyncProducer, error) { //nolint
	// Add SyncProducer specific properties to copy of base config
	sc.Producer.RequiredAcks = sarama.WaitForAll
	sc.Producer.Retry.Max = 5
	sc.Producer.Return.Successes = true

	maxMessageBytes := config.ProducerConfig.MaxMessageBytes

	if maxMessageBytes > 0 {
		sc.Producer.MaxMessageBytes = maxMessageBytes
	}

	producer, err := sarama.NewSyncProducer(config.ClientConfig.Brokers, &sc)
	if err != nil {
		return nil, err
	}

	// wrap tracing
	return otelsarama.WrapSyncProducer(&sc, producer), nil
}
