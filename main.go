package main

import (
	"context"
	"log"
	"net"

	pb "github.com/SirishaGopigiri/sample-grpc-server/user"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedUsersServer
}

func (s *server) SayHello(ctx context.Context, req *pb.UserRequest) (*pb.UserResponse, error) {
	return &pb.UserResponse{Message: "Hello " + req.Name}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterUsersServer(s, &server{})

	log.Println("Server is running on port 50051...")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
