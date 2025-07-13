package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"article-processing-microservice/database"
	"article-processing-microservice/proto"
	"article-processing-microservice/server"

	"google.golang.org/grpc"
)

func main() {

	if err := database.ConnectToMongoDB(); err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer func() {
		if err := database.CloseMongoDB(); err != nil {
			log.Printf("Error closing MongoDB connection: %v", err)
		}
	}()

	
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	
	grpcServer := grpc.NewServer()
	articleServer := server.NewArticleServer()
	proto.RegisterArticleServiceServer(grpcServer, articleServer)

	log.Println("Starting gRPC server on :50051...")
	log.Println("Server is ready to handle requests!")

	
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan
		log.Println("Shutting down server...")
		grpcServer.GracefulStop()
	}()

	
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
