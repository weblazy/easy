package ekafka

const (
	PackageNameProducer      = "ekafka.producer"
	PackageNameConsumerGroup = "ekafka.consumerGroup"
)

func NewProducer(name string, config *Config) (*Producer, error) {
	sc, err := config.toSaramaConfig()

	if err != nil {
		return nil, err
	}

	return newProducer(config, sc)
}

func NewConsumerGroup(name string, config *Config, groupConfig *ConsumerGroupConfig) (*ConsumerGroup, error) {
	sc, err := config.toSaramaConfig()

	if err != nil {
		return nil, err
	}

	return newConsumerGroup(config, groupConfig, sc)
}
