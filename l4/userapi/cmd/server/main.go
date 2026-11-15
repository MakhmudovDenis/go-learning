package main

import (
	"context"
	"fmt"
	"log"
	"net"

	userapipb "userapi/gen/userapi"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type UserService struct {
	userapipb.UnimplementedUserapiServiceServer
}

func (s *UserService) Create(ctx context.Context, req *userapipb.CreateRequest) (*userapipb.CreateResponse, error) {
	fmt.Printf("CREATE:\n%v\n", req)
	if req.Password != req.PasswordConfirm {
		return nil, status.Error(codes.Aborted, "passwords do not match")
	}

	return &userapipb.CreateResponse{
		Id: 1,
	}, nil
}

func (s *UserService) Get(ctx context.Context, req *userapipb.GetRequest) (*userapipb.GetResponse, error) {
	fmt.Printf("GET:\n%v\n", req)

	return &userapipb.GetResponse{
		Id:        req.GetId(),
		Name:      "vlad",
		Email:     "vl@d.d",
		Role:      userapipb.Role_ROLE_ADMIN,
		CreatedAt: timestamppb.Now(),
		UpdatedAt: timestamppb.Now(),
	}, nil
}

func (s *UserService) Update(ctx context.Context, req *userapipb.UpdateRequest) (*emptypb.Empty, error) {
	fmt.Printf("UPDATE:\n%v\n", req)

	return &emptypb.Empty{}, nil
}

func (s *UserService) Delete(ctx context.Context, req *userapipb.DeleteRequest) (*emptypb.Empty, error) {
	fmt.Printf("DELETE:\n%v\n", req)

	return &emptypb.Empty{}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer()

	log.Println(":50051")

	userapipb.RegisterUserapiServiceServer(grpcServer, &UserService{})

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
