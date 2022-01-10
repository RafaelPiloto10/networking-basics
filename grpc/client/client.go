package main

import (
	"context"
	"log"
	"time"

	pb "github.com/RafaelPiloto10/networking-basics/chat"
	"google.golang.org/grpc"
)


const (
	address = "localhost:8000"
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())

	if err != nil {
		log.Fatalf("got err triyng to dial conn: %v", err)
	}

	defer conn.Close()

	c := pb.NewChatServiceClient(conn)
	
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	message := pb.Message{Body: "Hello from the client!"}
	response, err := c.SayHello(ctx, &message)

	if err != nil {
		log.Fatalf("got error sending message: %v", err)
	}

	log.Printf("server response: %v", response)
}
