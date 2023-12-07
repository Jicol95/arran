package consumer

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/jicol-95/arran/config"
	"github.com/jicol-95/arran/domain"
	"github.com/labstack/echo/v4"
)

type ExampleResourceMessage struct {
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
			c.Consumer.logger.Error(err)
			continue
		}

		c.Consumer.logger.Info(fmt.Sprintf("Message received: %s", msg))
		c.svc.CreateExampleResource(msg.Name)
	}
}

func (c *ExampleResourceConsumer) readMessage() (*ExampleResourceMessage, error) {
	msg, err := c.Consumer.kafka.ReadMessage(time.Second * 5)

	if err != nil {
		return nil, err
	}

	var message ExampleResourceMessage
	if err := json.Unmarshal(msg.Value, &message); err != nil {
		return nil, err
	}

	return &message, nil
}
