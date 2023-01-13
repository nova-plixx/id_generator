package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pb "plixx.dev/id_generator/proto"
	"time"
)

func main() {
	conn, err := grpc.Dial("localhost:4040", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer func(conn *grpc.ClientConn) {
		_ = conn.Close()
	}(conn)
	client := pb.NewIdServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := client.GenerateId(ctx, &pb.GenerateIdRequest{})
	if err != nil {
		panic(err)
	}
	fmt.Println(r.Id)
}
