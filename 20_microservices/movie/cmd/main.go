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
	"movie.com/movie/internal/controller/movie"
	grpcmetadatagate "movie.com/movie/internal/gateway/metadata/grpc"
	httpmetadatagate "movie.com/movie/internal/gateway/metadata/http"
	grpcratinggate "movie.com/movie/internal/gateway/rating/grpc"
	httpratinggate "movie.com/movie/internal/gateway/rating/http"
	grpchandler "movie.com/movie/internal/handler/grpc"
	httphandler "movie.com/movie/internal/handler/http"
	"movie.com/pkg/discovery"
	"movie.com/pkg/discovery/consul"
)

const serviceName = "movie"

func main() {
	var port int
	flag.IntVar(&port, "port", 8083, "API handler port")
	var handler string
	flag.StringVar(&handler, "handler", "grpc", "API handler type")
	flag.Parse()
	log.Printf("Starting the movie service on port %d", port)
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

	switch handler {
	case "grpc":
		metadataGateway := grpcmetadatagate.New(registry)
		ratingGateway := grpcratinggate.New(registry)
		ctrl := movie.New(ratingGateway, metadataGateway)

		h := grpchandler.New(ctrl)
		lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%v", port))
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		srv := grpc.NewServer()
		reflection.Register(srv)
		gen.RegisterMovieServiceServer(srv, h)
		if err := srv.Serve(lis); err != nil {
			panic(err)
		}
	case "json":
		metadataGateway := httpmetadatagate.New(registry)
		ratingGateway := httpratinggate.New(registry)
		ctrl := movie.New(ratingGateway, metadataGateway)

		h := httphandler.New(ctrl)
		http.Handle("/movie", http.HandlerFunc(h.GetMovieDetails))
		if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
			panic(err)
		}
	}
}
