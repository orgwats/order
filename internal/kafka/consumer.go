package kafka

import (
	"github.com/IBM/sarama"
	"github.com/adshao/go-binance/v2/futures"
)

type Consumer struct {
	sarama.Consumer
}

type Order struct {
	OrderSide futures.SideType
	Key       string
	Symbol    string
	Size      string
}

func NewConsumer(addr string) (*Consumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	consumer, err := sarama.NewConsumer([]string{addr}, config)
	if err != nil {
		return nil, err
	}

	return &Consumer{consumer}, nil
}

func (c *Consumer) NewPartitionConsumer(topic string) (sarama.PartitionConsumer, error) {
	partitionConsumer, err := c.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if err != nil {
		return nil, err

	}
	return partitionConsumer, nil
}
