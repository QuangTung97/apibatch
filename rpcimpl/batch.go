package rpcimpl

import (
	"apibatch/rpc/batchpb"
	"context"
	"fmt"
)

// BatchServer ...
type BatchServer struct {
	batchpb.UnimplementedBatchServiceServer
}

// NewBatchServer ...
func NewBatchServer() *BatchServer {
	return &BatchServer{}
}

// Inc ...
func (s *BatchServer) Inc(ctx context.Context, req *batchpb.IncRequest) (*batchpb.IncResponse, error) {
	fmt.Println("Inc Called:", req)
	return &batchpb.IncResponse{}, nil
}
