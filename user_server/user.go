package user_server

import (
	"context"
	"fmt"

	pb "github.com/SirishaGopigiri/sample-grpc-server/user"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var users = map[string]*pb.User{}

type Server struct {
	pb.UnimplementedUsersServer
}

func (s *Server) GetUser(ctx context.Context, req *pb.UserReq) (*pb.User, error) {
	user_name := req.Name
	if _, exists := users[user_name]; !exists {
		error_str := fmt.Sprintf("user not found: %s", user_name)
		return nil, status.Error(codes.NotFound, error_str)
	}
	return users[user_name], nil
}

func (s *Server) GetUsers(ctx context.Context, req *pb.EmptyRequest) (*pb.UserList, error) {
	userslist := []*pb.User{}
	for _, value := range users {
		userslist = append(userslist, value)
	}
	return &pb.UserList{Users: userslist}, nil
}

func (s *Server) CreateUser(ctx context.Context, req *pb.User) (*pb.UserResponse, error) {
	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "name cannot be empty")
	} else if _, exists := users[req.Name]; exists {
		return nil, status.Error(codes.AlreadyExists, "user already exists, cannot create with same name")
	}
	users[req.Name] = req
	return &pb.UserResponse{Message: "Successfully created user"}, nil
}

func (s *Server) UpdateUser(ctx context.Context, req *pb.User) (*pb.UserResponse, error) {
	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "name cannot be empty")
	} else if _, exists := users[req.Name]; !exists {
		return nil, status.Error(codes.NotFound, "user doesn't exists to update, cannot update")
	}
	users[req.Name] = req
	return &pb.UserResponse{Message: "Successfully updated user"}, nil
}

func (s *Server) DeleteUser(ctx context.Context, req *pb.UserReq) (*pb.UserResponse, error) {
	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "name cannot be empty")
	} else if _, exists := users[req.Name]; !exists {
		return nil, status.Error(codes.NotFound, "user doesn't exists to update, cannot update")
	}
	delete(users, req.Name)
	return &pb.UserResponse{Message: "Successfully deleted user"}, nil
}
