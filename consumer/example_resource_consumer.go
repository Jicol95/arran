package consumer

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/jicol-95/arran/config"
	"github.com/jicol-95/arran/domain"
	"github.com/labstack/echo/v4"
)

type ExampleResourceMessage struct {
	Message KafkaMessage
	Data    ExampleResourceMessageData
}

type ExampleResourceMessageData struct {
	Name string `json:"name"`
}

type ExampleResourceConsumer struct {
	Consumer KafkaConsumer
	svc      domain.ExampleResourceService
}

func ProcessExampleResourceTopic(cfg config.KafkaConfig, svc domain.ExampleResourceService, logger echo.Logger) (ExampleResourceConsumer, error) {
	k, err := newConsumer(cfg, logger)

	if err != nil {
		return ExampleResourceConsumer{}, err
	}

	if err := k.start(cfg.ExampleConsumerTopic); err != nil {
		return ExampleResourceConsumer{}, err
	}

	c := ExampleResourceConsumer{
		Consumer: k,
		svc:      svc,
	}

	go c.processMessages()

	return c, nil
}

func (c *ExampleResourceConsumer) processMessages() {
	for {
		msg, err := c.readMessage()
		if err != nil {
			if kafkaErr, ok := err.(kafka.Error); ok {
				if kafkaErr.IsTimeout() {
					continue
				}
			}

			// TODO: this will cause poison pill scenario, need to implement dead letter queue.
			c.Consumer.logger.Error(fmt.Sprintf("Error reading message: %s", err))
			continue
		}

		exampleResource, err := unmarshalMessage(msg)

		if err != nil {
			c.Consumer.logger.Error(fmt.Sprintf("Error unmarshalling message: %s", err))
		}

		c.Consumer.logger.Info(fmt.Sprintf("Message received: %s", exampleResource.Data))
		_, err = c.svc.CreateExampleResource(exampleResource.Data.Name)

		if err != nil {
			// TODO: this will cause poison pill scenario, need to implement dead letter queue.
			c.Consumer.logger.Error("Error processing message")
		} else {
			c.Consumer.logger.Info("Message processed successfully")
			c.Consumer.kafka.CommitMessage(msg)
		}
	}
}

func (c *ExampleResourceConsumer) readMessage() (*kafka.Message, error) {
	msg, err := c.Consumer.kafka.ReadMessage(time.Second * 5)

	if err != nil {
		return nil, err
	}

	return msg, nil
}

func unmarshalMessage(msg *kafka.Message) (ExampleResourceMessage, error) {
	var data ExampleResourceMessageData
	if err := json.Unmarshal(msg.Value, &data); err != nil {
		return ExampleResourceMessage{}, err
	}

	return ExampleResourceMessage{
		Message: KafkaMessage{
			Topic:  msg.TopicPartition.Topic,
			Offset: int64(msg.TopicPartition.Offset),
		},
		Data: data,
	}, nil
}
