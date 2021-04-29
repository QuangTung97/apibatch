package rpcimpl

import (
	"apibatch/rpc/batchpb"
	"context"
	"github.com/jmoiron/sqlx"
)

// BatchServer ...
type BatchServer struct {
	batchpb.UnimplementedBatchServiceServer
	db *sqlx.DB
}

// NewBatchServer ...
func NewBatchServer(db *sqlx.DB) *BatchServer {
	return &BatchServer{
		db: db,
	}
}

// Inc ...
func (s *BatchServer) Inc(ctx context.Context, req *batchpb.IncRequest) (*batchpb.IncResponse, error) {
	query := `
INSERT INTO counters (id, val)
VALUES (?, ?)
ON DUPLICATE KEY UPDATE val = val + VALUES(val)
`
	_, err := s.db.ExecContext(ctx, query, req.Id, req.Value)
	if err != nil {
		return nil, err
	}

	return &batchpb.IncResponse{}, nil
}
