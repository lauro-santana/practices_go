package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"movie.com/gen"
	"movie.com/pkg/discovery"
	"movie.com/pkg/discovery/consul"
	"movie.com/rating/internal/controller/rating"
	grpchandler "movie.com/rating/internal/handler/grpc"
	httphandler "movie.com/rating/internal/handler/http"
	"movie.com/rating/internal/repository/memory"
)

const serviceName = "rating"

func main() {
	var port int
	flag.IntVar(&port, "port", 8082, "API handler port")
	var handler string
	flag.StringVar(&handler, "handler", "grpc", "API handler type")
	flag.Parse()
	log.Printf("Starting the rating service on port %d", port)
	registry, err := consul.NewRegistry("localhost:8500")
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	instanceID := discovery.GenerateInstanceID(serviceName)
	if err := registry.Register(ctx, instanceID, serviceName, fmt.Sprintf("localhost:%d", port)); err != nil {
		panic(err)
	}
	go func() {
		for {
			if err := registry.ReportHealthyState(instanceID, serviceName); err != nil {
				log.Println("Failed to report healthy state: " + err.Error())
			}
			time.Sleep(1 * time.Second)
		}
	}()
	defer registry.Deregister(ctx, instanceID, serviceName)
	repo := memory.New()
	ctrl := rating.New(repo)
	switch handler {
	case "grpc":
		h := grpchandler.New(ctrl)
		lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%v", port))
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		srv := grpc.NewServer()
		reflection.Register(srv)
		gen.RegisterRatingServiceServer(srv, h)
		if err := srv.Serve(lis); err != nil {
			panic(err)
		}
	case "json":
		h := httphandler.New(ctrl)
		http.Handle("/rating", http.HandlerFunc(h.Handle))
		if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
			panic(err)
		}
	}
}
