package main

import (
	"apibatch/rpc/batchpb"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"sync"
	"sync/atomic"
	"time"
)

// ConcurrentConfig ...
type ConcurrentConfig struct {
	ElementCount int
	ThreadCount  int
	Spacing      int
}

// Concurrent ...
func Concurrent(config ConcurrentConfig, fn func(i int), breakFn func()) {
	div := config.ElementCount / config.ThreadCount
	numHigher := config.ElementCount - config.ThreadCount*div

	count := uint64(0)

	var wg sync.WaitGroup
	wg.Add(config.ThreadCount)

	first := 0
	last := 0
	for t := 0; t < config.ThreadCount; t++ {
		if t < numHigher {
			last += div + 1
		} else {
			last += div
		}

		firstInner := first
		lastInner := last

		go func() {
			defer wg.Done()

			for i := firstInner; i < lastInner; i++ {
				fn(i)
				newCount := atomic.AddUint64(&count, 1)
				if config.Spacing > 0 && newCount%uint64(config.Spacing) == 0 {
					breakFn()
				}
			}
		}()

		first = last
	}

	wg.Wait()
}

func callConcurrent(conn *grpc.ClientConn) {
	a := time.Now()

	var startMut sync.Mutex
	start := time.Now()

	client := batchpb.NewBatchServiceClient(conn)
	Concurrent(ConcurrentConfig{
		ElementCount: 1000,
		ThreadCount:  10,
		Spacing:      100,
	}, func(i int) {
		_, err := client.Inc(context.Background(), &batchpb.IncRequest{
			Id:    100,
			Value: 1,
		})
		if err != nil {
			panic(err)
		}
	}, func() {
		startMut.Lock()
		now := time.Now()
		d := now.Sub(start)
		start = now
		startMut.Unlock()
		fmt.Println("TIME:", d)
	})
	fmt.Println("TOTAL:", time.Since(a))
}

func callSingleThread(conn *grpc.ClientConn) {
	start := time.Now()

	client := batchpb.NewBatchServiceClient(conn)
	for i := 0; i < 1000; i++ {
		_, err := client.Inc(context.Background(), &batchpb.IncRequest{
			Id:    100,
			Value: 1,
		})
		if err != nil {
			panic(err)
		}
	}

	fmt.Println("TIME:", time.Since(start))
}

func main() {
	conn, err := grpc.Dial("localhost:6000", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	callConcurrent(conn)
}
