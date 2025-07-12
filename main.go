package main

import (
	"log"
	"net"

	"article-processing-microservice/proto"
	"article-processing-microservice/server"

	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	articleServer := server.NewArticleServer()
	proto.RegisterArticleServiceServer(grpcServer, articleServer)

	log.Println("Starting gRPC server on :50051...")
	log.Println("Server is ready to handle requests!")

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
