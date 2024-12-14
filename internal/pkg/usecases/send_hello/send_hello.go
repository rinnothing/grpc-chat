package send_hello

import (
	"context"
	"errors"
	"time"

	"github.com/rinnothing/grpc-chat/internal/pkg/convert"
	"github.com/rinnothing/grpc-chat/internal/pkg/model"
	"github.com/rinnothing/grpc-chat/internal/pkg/presenter/dialogue"
	"github.com/rinnothing/grpc-chat/internal/pkg/repository/connections"
	desc "github.com/rinnothing/grpc-chat/pkg/generated/proto/chat"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type UserRepo interface {
	GetID(ctx context.Context, username string) (int, error)
}

type Connections interface {
	Connect(ctx context.Context, user *model.User) error
}

type MessageRepo interface {
	PutID(ctx context.Context, user *model.User, text string, time time.Time) (int, error)
}

type Dialogue interface {
	AskAccept(ctx context.Context, user *model.User, message *model.Message) error
}

type Chat interface {
	NewChat(ctx context.Context, user *model.User, message *model.Message) error
}

type Identify interface {
	Myself(ctx context.Context) *model.User
}

type UseCase struct {
	repo        UserRepo
	connections Connections
	message     MessageRepo
	dialogue    Dialogue
	chat        Chat
	identify    Identify
}

func New(repo UserRepo, connections Connections, message MessageRepo, dialogue Dialogue, chat Chat, identify Identify) *UseCase {
	return &UseCase{
		repo:        repo,
		connections: connections,
		message:     message,
		dialogue:    dialogue,
		chat:        chat,
		identify:    identify,
	}
}

func (uc *UseCase) makeResponse(ctx context.Context, allowed bool) *desc.SendHelloResponse {
	return &desc.SendHelloResponse{
		Addressee: convert.User2Credentials(uc.identify.Myself(ctx)),
		Allowed:   allowed,
		Time:      timestamppb.Now(),
	}
}

func (uc *UseCase) SendHello(ctx context.Context, req *desc.SendHelloRequest) (*desc.SendHelloResponse, error) {
	id, err := uc.repo.GetID(ctx, req.Sender.Username)
	if err != nil {
		return nil, status.Error(codes.Internal, "repository error")
	}

	user := convert.Credentials2User(req.Sender, id)

	err = uc.connections.Connect(ctx, user)
	if err != nil {
		if errors.Is(err, connections.ErrAlreadyConnected) {
			return uc.makeResponse(ctx, false), nil
		}
		return nil, status.Error(codes.Internal, "repository error")
	}

	sentTime := req.RequestText.Time.AsTime()
	messageID, err := uc.message.PutID(ctx, user, req.RequestText.Text, sentTime)
	if err != nil {
		return nil, status.Error(codes.Internal, "repository error")
	}

	message := convert.Text2Message(req.RequestText, user, messageID)

	err = uc.dialogue.AskAccept(ctx, user, message)
	if err != nil {
		if errors.Is(err, dialogue.ErrNotAllowed) {
			return uc.makeResponse(ctx, false), nil
		}
		return nil, status.Error(codes.Internal, "dialogue error")
	}

	err = uc.chat.NewChat(ctx, user, message)
	if err != nil {
		return nil, status.Error(codes.Internal, "dialogue error")
	}

	return uc.makeResponse(ctx, true), nil
}
