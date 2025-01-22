package main

import (
	"context"
	"log"
	"time"

	personGRPC "hello/person"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := personGRPC.NewPersonServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	name := "Lauro"
	resp, err := client.PersonCreate(ctx, &personGRPC.PersonCreateRequest{Person: &personGRPC.Person{Name: name}})
	if err != nil {
		log.Fatalf("error on create person: %v", err)
	}
	log.Printf("Person %v created!", resp.Person)

	list, err := client.PersonList(ctx, &personGRPC.PersonListRequest{})
	if err != nil {
		log.Fatalf("error on list persons: %v", err)
	}
	log.Println(list.Persons)
}
