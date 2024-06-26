package producer

import (
	"context"
	"fmt"

	"github.com/IBM/sarama"
	"github.com/namnv2496/go-coffee-shop-demo/pkg/configs"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Client interface {
	Produce(ctx context.Context, queueName string, payload []byte) error
}

type client struct {
	saramaSyncProducer sarama.SyncProducer
}

func newSaramaConfig(clientId string) *sarama.Config {
	saramaConfig := sarama.NewConfig()
	saramaConfig.Producer.Retry.Max = 1
	saramaConfig.Producer.RequiredAcks = sarama.WaitForAll
	saramaConfig.Producer.Return.Successes = true
	saramaConfig.ClientID = clientId
	saramaConfig.Metadata.Full = true
	return saramaConfig
}

func NewClient(
	mqConfig configs.Config,
) Client {
	saramaSyncProducer, err := sarama.NewSyncProducer(
		mqConfig.Kafka.Addresses,
		newSaramaConfig(mqConfig.Kafka.ClientID),
	)
	if err != nil {
		fmt.Println("failed to create sarama sync producer: ", err)
		return nil
	}

	return &client{
		saramaSyncProducer: saramaSyncProducer,
	}
}

func (c client) Produce(ctx context.Context, queueName string, payload []byte) error {

	if _, _, err := c.saramaSyncProducer.SendMessage(&sarama.ProducerMessage{
		Topic: queueName,
		Value: sarama.ByteEncoder(payload),
	}); err != nil {
		return status.Error(codes.Internal, "failed to produce message")
	}
	return nil
}
