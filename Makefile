.PHONY: \
	build \
	build-httpserver \
	build-grpcserver \
	run \
	run-httpserver \
	run-grpcserver \
	fmt

build: build-httpserver build-grpcserver

build-httpserver:
	go build -o dist/http_server cmd/http_server/*

build-grpcserver:
	go build -o dist/grpc_server cmd/grpc_server/*

run: run-httpserver run-grpcserver

run-httpserver:
	go run cmd/http_server/*

run-grpcserver:
	go run cmd/grpc_server/*

lint:
	golangci-lint run ./...

test:
	go test -race ./...

testv:
	go test -race -v ./...

fmt:
	go fmt ./...
	goimports -w cmd internal
