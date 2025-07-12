package main

import (
	"log"
	"net"

	"article-processing-microservice/proto"
	"article-processing-microservice/server"

	"google.golang.org/grpc"
)

func main() {
	// Create a TCP listener on port 50051
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Create a new gRPC server
	grpcServer := grpc.NewServer()

	// Create and register the article service
	articleServer := server.NewArticleServer()
	proto.RegisterArticleServiceServer(grpcServer, articleServer)

	log.Println("Starting gRPC server on :50051...")
	log.Println("Server is ready to handle requests!")

	// Start the server
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
