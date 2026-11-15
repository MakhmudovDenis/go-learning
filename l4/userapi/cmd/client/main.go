package main

import (
	"context"
	"fmt"

	userapipb "userapi/gen/userapi"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func main() {
	cc, err := grpc.NewClient(
		"localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		fmt.Printf("conn err: %v\n", err)
		return
	}

	client := userapipb.NewUserapiServiceClient(cc)

	createReq := userapipb.CreateRequest{
		Name:            "asd",
		Email:           "asd@s.d",
		Password:        "123321",
		PasswordConfirm: "123321",
		Role:            userapipb.Role_ROLE_ADMIN,
	}

	createRes, err := client.Create(context.Background(), &createReq)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(createRes)

	getReq := userapipb.GetRequest{
		Id: 11,
	}

	getRes, err := client.Get(context.Background(), &getReq)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(getRes)

	updateReq := userapipb.UpdateRequest{
		Id:    13,
		Name:  &wrapperspb.StringValue{Value: "wer"},
		Email: &wrapperspb.StringValue{Value: "dsa@s.d"},
	}

	updateRes, err := client.Update(context.Background(), &updateReq)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("res: %v\n", updateRes)

	deleteReq := userapipb.DeleteRequest{
		Id: 44,
	}

	deleteRes, err := client.Delete(context.Background(), &deleteReq)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("res: %v\n", deleteRes)
}
