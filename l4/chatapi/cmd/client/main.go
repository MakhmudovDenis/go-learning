package main

import (
	"context"
	"fmt"

	chatapipb "chatapi/gen/chatapi"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func main() {
	cc, err := grpc.NewClient(
		":50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		fmt.Printf("conn error:%v", err)
	}

	client := chatapipb.NewChatServiceClient(cc)

	createReq := chatapipb.CreateReq{
		Usernames: []string{"sergey", "anton", "nikolay"},
	}

	createRes, err := client.Create(context.Background(), &createReq)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("created chat with id: %d", createRes.Id)

	deleteReq := chatapipb.DeleteReq{
		Id: 41,
	}

	deleteRes, err := client.Delete(context.Background(), &deleteReq)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("delete response: %v\n", deleteRes)

	sendMessageReq := chatapipb.SendMessageReq{
		From:      "grigory",
		Text:      "hello there",
		Timestamp: timestamppb.Now(),
	}

	sendMessageRes, err := client.SendMessage(context.Background(), &sendMessageReq)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("send message response: %v\n", sendMessageRes)
}
