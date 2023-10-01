package ekafka

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	kafkaPublishCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "kafka_publish_total",
	}, []string{"brokers", "topic", "code"})

	kafkaConsumerGroupCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "kafka_consumer_group_handle_total",
	}, []string{"brokers", "group_id", "topic", "code"})

	kafkaConsumerGroupHistogram = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "kafka_consumer_group_handle_seconds",
		Buckets: []float64{.025, .05, .1, .25, .5, 1, 2.5, 5, 10, 30},
	}, []string{"brokers", "group_id", "topic"})
)

func init() {
	prometheus.MustRegister(kafkaPublishCounter)
	prometheus.MustRegister(kafkaConsumerGroupCounter)
	prometheus.MustRegister(kafkaConsumerGroupHistogram)
}
