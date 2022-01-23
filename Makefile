export PATH := $(GOPATH)/bin:$(PATH)

all: build

build: fmt frps frpc

frps:
	go build -o bin/mfrps.exe ./cmd/mfrps

frpc:
	go build -o bin/mfrpc.exe ./cmd/mfrpc

fmt:
	@go fmt ./...

test:
	@go test ./...
