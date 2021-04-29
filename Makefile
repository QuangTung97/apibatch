.PHONY: all format lint install-tools generate

RPC_DIR := rpc

CURRENT_DIR := $(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))

define generate
	mkdir -p rpc/$(1) && \
		cd proto/$(2) && \
		protoc -I.:${CURRENT_DIR}/proto \
			--go_out=paths=source_relative:${CURRENT_DIR}/rpc/$(1) \
			--go-grpc_out=paths=source_relative:${CURRENT_DIR}/rpc/$(1) \
			--grpc-gateway_out=logtostderr=true,paths=source_relative:${CURRENT_DIR}/rpc/$(1) \
			$(3)
endef

all:
	go build -o main cmd/server/main.go

format:
	go fmt ./...
	cd proto && prototool format -w

lint:
	golint ./...
	go vet ./...
	errcheck ./...
	gocyclo -over 10 .
	cd proto && prototool lint

install-tools:
	go install github.com/fzipp/gocyclo
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway
	go install github.com/kisielk/errcheck
	go install golang.org/x/lint/golint
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc
	go install google.golang.org/protobuf/cmd/protoc-gen-go

proto/google/api/httpbody.proto:
	mkdir -p proto/google/api
	wget -P proto/google/api/ https://raw.githubusercontent.com/googleapis/googleapis/master/google/api/annotations.proto
	wget -P proto/google/api/ https://raw.githubusercontent.com/googleapis/googleapis/master/google/api/field_behavior.proto
	wget -P proto/google/api/ https://raw.githubusercontent.com/googleapis/googleapis/master/google/api/http.proto
	wget -P proto/google/api/ https://raw.githubusercontent.com/googleapis/googleapis/master/google/api/httpbody.proto

generate: proto/google/api/httpbody.proto
	rm -rf rpc
	$(call generate,batchpb/,batch,batch.proto)
