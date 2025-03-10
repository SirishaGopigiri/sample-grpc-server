package main

import (
	"context"
	"log"
	"net"

	"github.com/SirishaGopigiri/sample-grpc-server/database"
	pb "github.com/SirishaGopigiri/sample-grpc-server/user"
	user "github.com/SirishaGopigiri/sample-grpc-server/user_server"

	"google.golang.org/grpc"
)

func loggingInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (interface{}, error) {
	log.Printf("Received request for method: %s", info.FullMethod)
	return handler(ctx, req)
}

func main() {
	ctx := context.Background()
	config, err := database.LoadConfig("config.json")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}
	db, err := database.Connect_to_DB(ctx, config)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer(grpc.UnaryInterceptor(loggingInterceptor))
	pb.RegisterUsersServer(s, &user.Server{DB: db})

	log.Println("Server is running on port 50051...")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
