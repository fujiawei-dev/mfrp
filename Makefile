export PATH := $(GOPATH)/bin:$(PATH)

all: build

build: frps frpc

frps:
	go build -o bin/mfrps.exe ./cmd/mfrps

frpc:
	go build -o bin/mfrpc.exe ./cmd/mfrpc
