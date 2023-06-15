package main

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/karta0898098/iam/pb/identity"
)

func main() {

	conn, err := grpc.Dial("localhost:9091", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewIdentityServiceClient(conn)

	ctx := context.Background()
	resp, err := c.Signin(ctx, &pb.SigninReq{
		Username:    "karta0898098",
		Password:    "A12345678",
		IPAddress:   "127.0.0.1",
		IdpProvider: "web",
		Device: &pb.Device{
			Model:     "OS X",
			Name:      "Ray Macbook",
			OSVersion: "11.0",
		},
	})
	if err != nil {
		fmt.Printf("err = %v", err)
		return
	}

	fmt.Printf("%v\n", resp)
}
