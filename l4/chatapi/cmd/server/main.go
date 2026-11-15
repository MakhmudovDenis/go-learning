package main

import (
	"context"
	"fmt"
	"log"
	"net"

	chatapipb "chatapi/gen/chatapi"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	_ "google.golang.org/protobuf/types/known/timestamppb"
)

type ChatService struct {
	chatapipb.UnimplementedChatServiceServer
}

func (s *ChatService) Create(ctx context.Context, in *chatapipb.CreateReq) (*chatapipb.CreateRes, error) {
	if len(in.Usernames) < 1 {
		return nil, status.Error(codes.Aborted, "cannot create chat with less than 1 user")
	}

	fmt.Printf("crated chat for %v\n", in.Usernames)

	return &chatapipb.CreateRes{Id: 20}, nil
}

func (s *ChatService) Delete(ctx context.Context, in *chatapipb.DeleteReq) (*emptypb.Empty, error) {
	if in.Id == 42 {
		return nil, status.Error(codes.Aborted, "cannot delete chat 42")
	}

	fmt.Printf("deleted chat %v\n", in.Id)

	return &emptypb.Empty{}, nil
}

func (s *ChatService) SendMessage(ctx context.Context, in *chatapipb.SendMessageReq) (*emptypb.Empty, error) {
	fmt.Printf("%v\n\t%s says: %s", in.Timestamp.AsTime(), in.From, in.Text)

	return &emptypb.Empty{}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		fmt.Println(err)
		return
	}

	grpcServer := grpc.NewServer()

	chatapipb.RegisterChatServiceServer(grpcServer, &ChatService{})

	fmt.Println("listening on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
