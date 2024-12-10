api=api
proto_target=pkg/generated/proto

generate:
	protoc -I=${api} --go_out=${proto_target} --go_opt=paths=source_relative --go-grpc_out=${proto_target} \
	--go-grpc_opt=paths=source_relative ${api}/chat/chat.proto
	go mod tidy
