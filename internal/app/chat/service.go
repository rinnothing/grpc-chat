package chat

import (
	"context"
	"time"

	"github.com/rinnothing/grpc-chat/internal/pkg/model"
	desc "github.com/rinnothing/grpc-chat/pkg/generated/proto/chat"
)

type (
	SendHelloUseCase interface {
		// SendHello makes a request for connection with text and author from message and returns true if allowed
		SendHello(context.Context, *model.Message) (bool, error)
	}

	SendMessageUseCase interface {
		// SendMessage sends message on established connection
		SendMessage(context.Context, *model.Message) error
	}

	SendGoodbyeUseCase interface {
		// SendGoodbye ends connection with given user (also brings time as an additional data)
		SendGoodbye(context.Context, *model.User, time.Time) error
	}
)

type Implementation struct {
	desc.UnimplementedChatInstanceServer
	SendHelloUseCase   SendHelloUseCase
	SendMessageUseCase SendMessageUseCase
	SendGoodbyeUseCase SendGoodbyeUseCase
}
