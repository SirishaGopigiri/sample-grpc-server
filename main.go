package main

import (
	"log"
	"net"

	pb "github.com/SirishaGopigiri/sample-grpc-server/user"
	user "github.com/SirishaGopigiri/sample-grpc-server/user_server"

	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterUsersServer(s, &user.Server{})

	log.Println("Server is running on port 50051...")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
