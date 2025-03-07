package main

import (
	"context"
	"log"
	"time"

	pb "github.com/SirishaGopigiri/sample-grpc-server/user"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.NewClient("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewUsersClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := client.SayHello(ctx, &pb.UserRequest{Name: "Gopher"})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	log.Println("Greeting:", resp.Message)
}
