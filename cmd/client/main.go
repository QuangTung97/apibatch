package main

import (
	"apibatch/rpc/batchpb"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"time"
)

func main() {
	conn, err := grpc.Dial("localhost:6000", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	start := time.Now()
	for i := 0; i < 1000; i++ {
		client := batchpb.NewBatchServiceClient(conn)
		_, err = client.Inc(context.Background(), &batchpb.IncRequest{
			Id:    100,
			Value: 1,
		})
		if err != nil {
			panic(err)
		}
	}
	fmt.Println("TIME:", time.Since(start))
}
