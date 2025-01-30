api := api
proto_target := pkg/generated/proto

LOCAL_BIN := $(CURDIR)/bin

# Detect the operating system
UNAME_S := $(shell uname -s)

ifeq ($(UNAME_S),Darwin)
	CMD = brew install protobuf
else ifeq ($(UNAME_S),Linux)
    CMD = apt install -y protobuf-compiler
else ifeq ($(OS), Windows_NT)
	CMD = winget install protobuf --force
else
	$(error OS not supported)
endif

.PHONY: install-deps
install-deps:
	# installing protoc and it's modules
	$(CMD)
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

generate: install-deps
	protoc -I=${api} --go_out=${proto_target} --go_opt=paths=source_relative --go-grpc_out=${proto_target} \
	--go-grpc_opt=paths=source_relative ${api}/chat/chat.proto
	go mod tidy
