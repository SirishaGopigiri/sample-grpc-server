package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"sync"
	"time"

	pb "github.com/SirishaGopigiri/sample-grpc-server/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
)

var wg sync.WaitGroup

func main() {
	conn, err := grpc.NewClient(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewUsersClient(conn)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for {
		fmt.Println("GRPC Client with options:")
		fmt.Println("1. Stream Request")
		fmt.Println("2. Stream Response")
		fmt.Println("3. Stream Req Response")
		fmt.Println("0. Exit")
		fmt.Print("Please select and option from above: ")
		var value int
		fmt.Scan(&value)
		switch value {
		case 1:
			fmt.Println("Sending Stream Request")
			SendStreamrequest(ctx, client)
		case 2:
			fmt.Println("Sending Stream Response Request")
			SendReqStreamResp(ctx, client)
		case 3:
			fmt.Println("Sending Stream Request and getting Stream Response:")
			SendReqRespClient(ctx, client)
		case 0:
			return
		default:
			fmt.Println("Unknown input, please select proper input!")
		}
	}
}

func SendReqStreamResp(ctx context.Context, client pb.UsersClient) {
	fmt.Println("Enter name to fetch:")
	var name string
	fmt.Scan(&name)
	user := &pb.UserReq{Name: name}
	stream, err := client.StreamResponse(ctx, user)
	if err != nil {
		log.Fatalf("Error calling StreamReqResp: %v", err)
	}
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error receiving stream: %v", err)
		}
		time.Sleep(1 * time.Second)
		fmt.Println("Server Response:", msg.GetMessage())
	}
}

func SendStreamrequest(ctx context.Context, client pb.UsersClient) {
	stream, err := client.StreamRequest(ctx)
	if err != nil {
		log.Fatalf("Error calling StreamReqResp: %v", err)
	}
	for {
		fmt.Println("Enter name to fetch:")
		var name string
		fmt.Scan(&name)
		if name == "exit" {
			break
		}
		stream.Send(&pb.UserReq{Name: name})
		time.Sleep(time.Second)
	}
	resp, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error receiving response: %v", err)
	}
	fmt.Println("Server Response:", resp.GetMessage())
}

func SendReqRespClient(ctx context.Context, client pb.UsersClient) {
	stream, err := client.StreamReqResp(ctx)
	if err != nil {
		log.Fatalf("Error calling StreamReqResp: %v", err)
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			fmt.Println("Enter name to fetch:")
			var name string
			fmt.Scan(&name)
			if name == "exit" {
				stream.CloseSend()
				return
			}
			stream.Send(&pb.UserReq{Name: name})
			time.Sleep(time.Second)
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			msg, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error receiving: %v", err)
			}
			fmt.Println("Server:", msg.GetMessage())
		}
	}()
	wg.Wait()
}
