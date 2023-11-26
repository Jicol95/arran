package consumer

import (
	"fmt"
	"time"

	"github.com/jicol-95/arran/config"
	"github.com/jicol-95/arran/domain"
	"github.com/labstack/echo/v4"
)

type ExapmleResourceConsumer struct {
	Consumer KafkaConsumer
	svc      domain.ExampleResourceService
}

func ProcessExampleResourceTopic(cfg config.KafkaConfig, svc domain.ExampleResourceService, logger echo.Logger) (ExapmleResourceConsumer, error) {
	k, err := newConsumer(cfg, logger)

	if err != nil {
		return ExapmleResourceConsumer{}, err
	}

	if err := k.start(cfg.ExampleConsumerTopic); err != nil {
		return ExapmleResourceConsumer{}, err
	}

	c := ExapmleResourceConsumer{
		Consumer: k,
		svc:      svc,
	}

	go c.processMessages()

	return c, nil
}

func (c *ExapmleResourceConsumer) processMessages() {
	for {
		msg, err := c.Consumer.kafka.ReadMessage(time.Second * 5)
		if err != nil {
			c.Consumer.logger.Error(err)
			continue
		}

		c.Consumer.logger.Info(fmt.Sprintf("Message received: %s", string(msg.Value)))
		c.svc.CreateExampleResource(string(msg.Value))
	}
}
