package main

import (
	"log"
	"net"

	pb "github.com/RafaelPiloto10/networking-basics/chat"
	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":8000")
	if err != nil {
		log.Fatalf("Failed to listen on port 9000: %v", err)
	}

	s := grpc.NewServer()

	pb.RegisterChatServiceServer(s, &pb.Server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC server over port 8000: %v", err)
	}
}
