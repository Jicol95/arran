package config

import (
	"os"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type ArranConfig struct {
	PostgresConfig PostgresConfig
	Kafka          KafkaConfig
}

type KafkaConfig struct {
	BaseConfig           kafka.ConfigMap
	ExampleConsumerTopic string
}

type PostgresConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	Database string
}

func NewArranConfig() ArranConfig {
	return ArranConfig{
		PostgresConfig: PostgresConfig{
			Host:     os.Getenv("POSTGRES_HOST"),
			Port:     os.Getenv("POSTGRES_PORT"),
			Username: os.Getenv("POSTGRES_USERNAME"),
			Password: os.Getenv("POSTGRES_PASSWORD"),
			Database: os.Getenv("POSTGRES_DATABASE"),
		},
		Kafka: KafkaConfig{
			BaseConfig: kafka.ConfigMap{
				"bootstrap.servers":     os.Getenv("KAFKA_BOOTSTRAP_SERVERS"),
				"broker.address.family": "v4",
				"group.id":              os.Getenv("KAFKA_GROUP_ID"),
				"session.timeout.ms":    os.Getenv("KAFKA_SESSION_TIMEOUT_MS"),
				"auto.offset.reset":     os.Getenv("KAFKA_AUTO_OFFSET_RESET"),
				"enable.auto.commit":    false,
			},
			ExampleConsumerTopic: os.Getenv("KAFKA_EXAMPLE_CONSUMER_TOPICS"),
		},
	}
}
