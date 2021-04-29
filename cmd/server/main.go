package main

import (
	"apibatch/rpc/batchpb"
	"apibatch/rpcimpl"
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"
)

func main() {
	server := grpc.NewServer()
	s := rpcimpl.NewBatchServer()
	batchpb.RegisterBatchServiceServer(server, s)

	runGRPCAndHTTPServers(server)
}

func registerGateway(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) {
	err := batchpb.RegisterBatchServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	if err != nil {
		panic(err)
	}
}

func runGRPCAndHTTPServers(server *grpc.Server) {
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, os.Kill)

	listener, err := net.Listen("tcp", ":9000")
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	registerGateway(ctx, mux, "localhost:9000", opts)

	http.Handle("/api/", mux)
	httpServer := http.Server{
		Addr: ":9080",
	}

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()

		err := server.Serve(listener)
		if err != nil {
			panic(err)
		}
	}()

	go func() {
		defer wg.Done()
		err := httpServer.ListenAndServe()
		if err == http.ErrServerClosed {
			return
		}
		if err != nil {
			panic(err)
		}
	}()

	<-exit
	server.GracefulStop()

	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()
	err = httpServer.Shutdown(ctx)
	if err != nil {
		panic(err)
	}

	wg.Wait()
	fmt.Println("Stop successfully")
}
