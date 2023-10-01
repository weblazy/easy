package ekafka

import "github.com/IBM/sarama"

// Message sarama.ProducerMessage for kafka publish
type Message struct {
	Topic string
	Key   []byte
	Value []byte
	// The headers are key-value pairs that are transparently passed
	// by Kafka between producers and consumers.
	Headers []sarama.RecordHeader

	// This field is used to hold arbitrary data you wish to include, so it
	// will be available when receiving on the Successes and Errors channels.
	// Sarama completely ignores this field and is only to be used for
	// pass-through data.
	Metadata interface{}
}

func (m *Message) ToMap() map[string]interface{} {
	mp := make(map[string]interface{}, 5)
	mp["topic"] = m.Topic
	mp["value"] = string(m.Value)

	if m.Key != nil {
		mp["key"] = m.Key
	}

	if len(m.Headers) > 0 {
		headers := make(map[string]string, len(m.Headers))
		for _, h := range m.Headers {
			headers[string(h.Key)] = string(h.Value)
		}
		mp["headers"] = headers
	}

	if m.Metadata != nil {
		mp["metadata"] = m.Metadata
	}

	return mp
}

func ConsumerMessageToMap(m *sarama.ConsumerMessage) map[string]interface{} {
	mp := make(map[string]interface{}, 7)
	mp["topic"] = m.Topic
	mp["value"] = string(m.Value)
	mp["partition"] = m.Partition
	mp["offset"] = m.Offset
	mp["timestamp"] = m.Timestamp

	if m.Key != nil {
		mp["key"] = m.Key
	}

	if len(m.Headers) > 0 {
		headers := make(map[string]string, len(m.Headers))
		for _, h := range m.Headers {
			headers[string(h.Key)] = string(h.Value)
		}
		mp["headers"] = headers
	}

	return mp
}
