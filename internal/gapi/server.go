package gapi

import (
	pb "github.com/orgwats/idl/gen/go/order"
	"github.com/orgwats/order/internal/config"
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
