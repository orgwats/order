package kafka

import (
	"encoding/json"
	"log"

	"github.com/IBM/sarama"
	"github.com/adshao/go-binance/v2/futures"
)

type Producer struct {
	sarama.SyncProducer
}

type OrderResult struct {
	OrderSide futures.SideType
	IsFilled  bool
}

func NewProducer(addr string) (*Producer, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll // 모든 브로커 부터 응답 대기
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer([]string{addr}, config)
	if err != nil {
		return nil, err
	}

	return &Producer{producer}, nil
}

func (p *Producer) Send(analyzerId string, result *OrderResult) (partition int32, offset int64, err error) {
	data, _ := json.Marshal(result)

	msg := &sarama.ProducerMessage{
		Topic: "order-result",
		Key:   sarama.StringEncoder(analyzerId),
		Value: sarama.StringEncoder(string(data)),
	}

	partition, offset, err = p.SendMessage(msg)
	if err != nil {
		log.Printf("Failed to send message: %v", err)
		return
	}

	return
}
