package main

import (
	"context"
	"log"
	"net"

	personGRPC "hello/person"

	"google.golang.org/grpc"
)

type personServer struct {
	personGRPC.UnimplementedPersonServiceServer
}

var persons []*personGRPC.Person

func (s *personServer) PersonList(ctx context.Context, req *personGRPC.PersonListRequest) (*personGRPC.PersonListReply, error) {
	return &personGRPC.PersonListReply{Persons: persons}, nil
}

func (s *personServer) PersonCreate(ctx context.Context, req *personGRPC.PersonCreateRequest) (*personGRPC.PersonCreateReply, error) {
	log.Printf("Received: %v", req.Person)
	//service logic... database things... bla bla bla...
	persons = append(persons, &personGRPC.Person{Id: uint64(len(persons) + 1), Name: req.Person.Name})
	return &personGRPC.PersonCreateReply{Person: &personGRPC.Person{Id: uint64(len(persons)), Name: req.Person.Name}}, nil
}

func main() {
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	persons = make([]*personGRPC.Person, 0)

	grpcServer := grpc.NewServer()

	personGRPC.RegisterPersonServiceServer(grpcServer, &personServer{})
	log.Println("Server is listening on port 50051")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
