package main

import (
	"apibatch/rpc/batchpb"
	"context"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:9000", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	client := batchpb.NewBatchServiceClient(conn)
	_, err = client.Inc(context.Background(), &batchpb.IncRequest{
		Id:    100,
		Value: 200,
	})
	if err != nil {
		panic(err)
	}
}
