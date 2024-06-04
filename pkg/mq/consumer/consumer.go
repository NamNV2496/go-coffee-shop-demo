package consumer

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/IBM/sarama"
	"github.com/namnv2496/go-coffee-shop-demo/pkg/configs"
)

type HandlerFunc func(ctx context.Context, queueName string, payload []byte) error

type consumerHandler struct {
	handlerFunc       HandlerFunc
	exitSignalChannel chan os.Signal
}

func newConsumerHandler(
	handlerFunc HandlerFunc,
	exitSignalChannel chan os.Signal,
) *consumerHandler {
	return &consumerHandler{
		handlerFunc:       handlerFunc,
		exitSignalChannel: exitSignalChannel,
	}
}

func (h consumerHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h consumerHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h consumerHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message, ok := <-claim.Messages():
			if !ok {
				session.Commit()
				return nil
			}

			if err := h.handlerFunc(session.Context(), message.Topic, message.Value); err != nil {
				return err
			}

		case <-h.exitSignalChannel:
			session.Commit()
			break
		}
	}
}

// consumer

type Consumer interface {
	RegisterHandler(queueName string, handlerFunc HandlerFunc)
}

type consumer struct {
	worker                    sarama.ConsumerGroup
	queueNameToHandlerFuncMap map[string]HandlerFunc
}

func newSaramaConfig(mqConfig configs.Kafka) *sarama.Config {
	saramaConfig := sarama.NewConfig()
	saramaConfig.ClientID = mqConfig.ClientID
	saramaConfig.Metadata.Full = true
	return saramaConfig
}
func NewConsumer(
	configs configs.Config,
) Consumer {

	saramaConsumer, err := sarama.NewConsumerGroup(configs.Kafka.Addresses, configs.Kafka.ClientID, newSaramaConfig(configs.Kafka))
	if err != nil {
		fmt.Println("failed to create sarama consumer: ", err)
		return nil
	}

	return &consumer{
		worker:                    saramaConsumer,
		queueNameToHandlerFuncMap: make(map[string]HandlerFunc),
	}
}

func (c consumer) RegisterHandler(queueName string, handlerFunc HandlerFunc) {

	exitSignalChannel := make(chan os.Signal, 1)
	signal.Notify(exitSignalChannel, os.Interrupt)

	go func(queueName string, handlerFunc HandlerFunc) {
		if err := c.worker.Consume(
			context.Background(),
			[]string{queueName},
			newConsumerHandler(handlerFunc, exitSignalChannel),
		); err != nil {
			fmt.Println("Fail to consume")
		}
	}(queueName, handlerFunc)
}
