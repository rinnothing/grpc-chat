package chat

import (
	"context"

	desc "github.com/rinnothing/grpc-chat/pkg/generated/proto/chat"
)

type (
	SendHelloUseCase interface {
		SendHello(context.Context, *desc.SendHelloRequest) (*desc.SendHelloResponse, error)
	}

	SendMessageUseCase interface {
		SendMessage(context.Context, *desc.SendMessageRequest) (*desc.SendMessageResponse, error)
	}

	SendGoodbyeUseCase interface {
		SendGoodbye(context.Context, *desc.SendGoodbyeRequest) (*desc.SendGoodbyeResponse, error)
	}
)

type Implementation struct {
	desc.UnimplementedChatInstanceServer
	SendHelloUseCase SendHelloUseCase
	SendMessageUseCase SendMessageUseCase
	SendGoodbyeUseCase SendGoodbyeUseCase
}
