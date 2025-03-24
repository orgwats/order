package main

import (
	"log"
	"net"
	"sync"

	pb "github.com/orgwats/idl/gen/go/order"
	"github.com/orgwats/order/internal/config"
	"github.com/orgwats/order/internal/gapi"
	"google.golang.org/grpc"
)

func main() {
	// TODO 임시
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	// TODO 임시
	store := struct{}{}
	server := gapi.NewServer(cfg, store)
	grpcServer := grpc.NewServer()

	pb.RegisterOrderServer(grpcServer, server)

	listener, err := net.Listen("tcp", ":50053")
	if err != nil {
		log.Fatal("cannot listen network address:", err)
	}

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Printf("start order GRPC server at %s", listener.Addr().String())

		if err := grpcServer.Serve(listener); err != nil {
			log.Fatal("order server failed to serve:", err)
		}
	}()
	wg.Wait()
}
