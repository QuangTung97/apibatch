package server

import (
	"apibatch/rpc/batchrpc"
	"context"
	"fmt"
)

// BatchServer ...
type BatchServer struct {
	batchrpc.UnimplementedBatchServiceServer
}

func (s *BatchServer) Inc(ctx context.Context, req *batchrpc.IncRequest) (*batchrpc.IncResponse, error) {
	fmt.Println("Inc Called")
	return &batchrpc.IncResponse{}, nil
}
