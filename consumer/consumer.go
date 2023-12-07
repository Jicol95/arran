package consumer

import (
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/jicol-95/arran/config"
	"github.com/labstack/echo/v4"
)

type KafkaConsumer struct {
	kafka  *kafka.Consumer
	logger echo.Logger
}

func (c *KafkaConsumer) start(topic string) error {
	if err := c.kafka.Subscribe(topic, nil); err != nil {
		c.logger.Fatal(err)
		return err
	}

	return nil
}

func newConsumer(cfc config.KafkaConfig, logger echo.Logger) (KafkaConsumer, error) {
	k, err := kafka.NewConsumer(&cfc.BaseConfig)

	if err != nil {
		logger.Fatal(err)
		return KafkaConsumer{}, err
	}

	return KafkaConsumer{kafka: k, logger: logger}, nil
}
