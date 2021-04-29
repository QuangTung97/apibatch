package rpcimpl

import (
	"apibatch/rpc/batchpb"
	"bytes"
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
)

// Request ...
type Request struct {
	ID     uint64
	Value  uint64
	WaitCh chan<- struct{}
}

// BatchServer ...
type BatchServer struct {
	ch chan<- Request
	batchpb.UnimplementedBatchServiceServer
	db *sqlx.DB
}

// NewBatchServer ...
func NewBatchServer(db *sqlx.DB) *BatchServer {
	ch := make(chan Request, 1000)

	s := &BatchServer{
		ch: ch,
		db: db,
	}
	go func() {
		s.doInBackground(ch)
	}()
	return s
}

// Inc ...
func (s *BatchServer) Inc(ctx context.Context, req *batchpb.IncRequest) (*batchpb.IncResponse, error) {
	return s.incWithBatching(ctx, req)
}

func (s *BatchServer) incWithoutBatching(ctx context.Context, req *batchpb.IncRequest) (*batchpb.IncResponse, error) {
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

func (s *BatchServer) incWithBatching(ctx context.Context, req *batchpb.IncRequest) (*batchpb.IncResponse, error) {
	waitCh := make(chan struct{}, 1)
	s.ch <- Request{
		ID:     req.Id,
		Value:  req.Value,
		WaitCh: waitCh,
	}

	select {
	case <-ctx.Done():
	case <-waitCh:
	}
	return &batchpb.IncResponse{}, nil
}

func (s *BatchServer) doInBackground(ch <-chan Request) {
	requests := make([]Request, 0, 1000)
	for {
		first := <-ch
		requests = append(requests, first)

	BatchLoop:
		for i := 1; i < 1000; i++ {
			select {
			case req := <-ch:
				requests = append(requests, req)
			default:
				break BatchLoop
			}
		}

		values := map[uint64]uint64{}
		for _, req := range requests {
			values[req.ID] += req.Value
		}

		args := make([]interface{}, 0, 2*len(values))
		for key, value := range values {
			args = append(args, key, value)
		}

		var buf bytes.Buffer
		buf.WriteString("(?, ?)")
		for i := 0; i < len(values)-1; i++ {
			buf.WriteString(",(?, ?)")
		}
		fmt.Println("BATCH:", values)

		query := `
INSERT INTO counters (id, val)
VALUES %s
ON DUPLICATE KEY UPDATE val = val + VALUES(val)
`
		query = fmt.Sprintf(query, buf.String())
		_, err := s.db.Exec(query, args...)
		if err != nil {
			panic(err)
		}

		for _, req := range requests {
			req.WaitCh <- struct{}{}
			req.WaitCh = nil
		}

		requests = requests[:0]
	}
}
