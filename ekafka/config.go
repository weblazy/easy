package fkafka

import (
	"strings"

	"errors"

	"github.com/IBM/sarama"
)

const (
	passwordAuthType = "password"
	noAuthType       = "none"
)

const (
	offsetOldest = "oldest"
	offsetNewest = "newest"
)

type Config struct {
	Debug                      bool
	EnableAccessInterceptorReq bool // 是否开启记录 publish 消息，默认开启
	EnableAccessInterceptorRes bool // 是否开启记录 consumer 消费消息, 默认开启

	ClientConfig         ClientConfig
	ProducerConfig       ProducerConfig
	ConsumerGroupConfigs map[string]ConsumerGroupConfig
}

type ClientConfig struct {
	Brokers      []string // Brokers brokers地址
	AuthType     string   // 鉴权方式, password / none. 默认为: none
	SaslUsername string   // 鉴权方式为 password 时, 必填
	SaslPassword string   // 鉴权方式为 password 时, 必填
	Version      string   // kafka version, 默认为 2.0.0.0
}

type ProducerConfig struct {
	MaxMessageBytes int
}

type RetryConfig struct {
	MaxRetries int64 // consumer 消费重试次数, 默认 0 不重试
}

type ConsumerGroupConfig struct {
	Topics        []string
	GroupID       string
	InitialOffset string // 初始化 offset, oldest / newest, 默认 oldest
	RetryConfig   RetryConfig
}

func DefaultConfig() *Config {
	return &Config{
		Debug:                      false,
		EnableAccessInterceptorReq: true,
		EnableAccessInterceptorRes: true,
		ClientConfig:               ClientConfig{AuthType: noAuthType},
	}
}

func (c *Config) toSaramaConfig() (*sarama.Config, error) { //nolint
	sc := sarama.NewConfig()
	clientConfig := c.ClientConfig

	if len(c.ClientConfig.Brokers) == 0 {
		return nil, errors.New("empty brokers")
	}

	if clientConfig.AuthType == passwordAuthType {
		if clientConfig.SaslUsername == "" || clientConfig.SaslPassword == "" {
			return nil, errors.New("username and password are required when using password auth type")
		}
		sc.Net.SASL.Enable = true
		sc.Net.SASL.User = clientConfig.SaslUsername
		sc.Net.SASL.Password = clientConfig.SaslPassword
		sc.Net.SASL.Mechanism = sarama.SASLTypePlaintext
	}

	if clientConfig.Version != "" {
		v, err := sarama.ParseKafkaVersion(clientConfig.Version)
		if err != nil {
			return nil, errors.New("invalid kafka version")
		}
		sc.Version = v
	} else {
		sc.Version = sarama.V2_0_0_0
	}

	return sc, nil
}

func (c *Config) brokers() string { //nolint
	return strings.Join(c.ClientConfig.Brokers, ",")
}
