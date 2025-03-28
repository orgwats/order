package gapi

import (
	"encoding/json"
	"log"

	"github.com/adshao/go-binance/v2/futures"
	pb "github.com/orgwats/idl/gen/go/order"
	"github.com/orgwats/order/internal/config"
	"github.com/orgwats/order/internal/kafka"
)

type Server struct {
	pb.UnimplementedOrderServer

	// TODO: 임시
	cfg   *config.Config
	store interface{}
}

func NewServer(cfg *config.Config, store interface{}) *Server {
	return &Server{
		cfg:   cfg,
		store: store,
	}
}

func (s *Server) Run() {
	// testnet
	futures.UseTestnet = true

	// s.cfg.~~~
	apiKey := ""
	secretKey := ""
	orderPlaceService, err := futures.NewOrderPlaceWsService(apiKey, secretKey)
	if err != nil {
		log.Fatal("cannot connect binance websocket:", err)
	}

	// s.cfg.~~~
	addr := "localhost:9092"
	consumer, err := kafka.NewConsumer(addr)
	if err != nil {
		log.Fatal("cannot create kafka consumer:", err)
	}

	// s.cfg.~~~
	topic := "order"
	partitionConsumer, err := consumer.NewPartitionConsumer(topic)
	if err != nil {
		log.Fatal("cannot create partitionConsumer:", err)
	}

	// s.cfg.~~~
	addr = "localhost:9092"
	producer, err := kafka.NewProducer(addr)
	if err != nil {
		log.Fatal("cannot create kafka producer:", err)
	}

	orderMap := make(map[string]string)

	go func(orderPlaceService *futures.OrderPlaceWsService) {
		for {
			select {
			case err := <-partitionConsumer.Errors():
				log.Printf("Consumer error: %v", err)
			case msg := <-partitionConsumer.Messages():
				var order kafka.Order
				err := json.Unmarshal(msg.Value, &order)
				if err != nil {
					log.Println("cannot unmarshal kafka message")
				}

				orderMap[order.Key] = string(msg.Key)

				err = orderPlaceService.Do(
					order.Key,
					futures.NewOrderPlaceWsRequest().
						Symbol(order.Symbol).
						Side(order.OrderSide).
						Type(futures.OrderTypeMarket). // 시장가 주문
						Quantity(order.Size),
				)
				if err != nil {
					log.Println("failed order:", err)
				}
			}
		}
	}(orderPlaceService)

	go func(orderPlaceService *futures.OrderPlaceWsService, producer *kafka.Producer) {
		for msg := range orderPlaceService.GetReadChannel() {
			var response futures.CreateOrderWsResponse
			err := json.Unmarshal(msg, &response)
			if err != nil {
				log.Println("cannot unmarshal order response")
			}

			analyzerId := orderMap[response.Id]

			result := &kafka.OrderResult{
				OrderSide: response.Result.Side,
				IsFilled:  response.Status == 200,
			}

			_, _, err = producer.Send(analyzerId, result)
			if err != nil {
				// 필요에 따라 재시도 로직
			}

			log.Println("order place response", string(msg))
		}
	}(orderPlaceService, producer)

	go func(orderPlaceService *futures.OrderPlaceWsService) {
		for err := range orderPlaceService.GetReadErrorChannel() {
			log.Println("order place error", err)
		}
	}(orderPlaceService)
}
