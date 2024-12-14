package send_message

import (
	"context"
	"errors"
	"time"

	"github.com/rinnothing/grpc-chat/internal/pkg/convert"
	"github.com/rinnothing/grpc-chat/internal/pkg/model"
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
	IsConnected(ctx context.Context, user *model.User) error
}

type MessageRepo interface {
	PutID(ctx context.Context, user *model.User, text string, time time.Time) (int, error)
}

type Chat interface {
	NewMessage(ctx context.Context, user *model.User, message *model.Message) error
}

type Identify interface {
	Myself(ctx context.Context) *model.User
}

type UserCase struct {
	repo        UserRepo
	connections Connections
	message     MessageRepo
	chat        Chat
	identify    Identify
}

func New(repo UserRepo, connections Connections, message MessageRepo, chat Chat, identify Identify) *UserCase {
	return &UserCase{
		repo:        repo,
		connections: connections,
		message:     message,
		chat:        chat,
		identify:    identify,
	}
}

func (uc *UserCase) SendMessage(ctx context.Context, req *desc.SendMessageRequest) (*desc.SendMessageResponse, error) {
	id, err := uc.repo.GetID(ctx, req.Sender.Username)
	if err != nil {
		return nil, status.Error(codes.Internal, "repository error")
	}

	user := convert.Credentials2User(req.Sender, id)

	err = uc.connections.IsConnected(ctx, user)
	if err != nil {
		if errors.Is(err, connections.ErrNotConnected) {
			return nil, status.Error(codes.Unauthenticated, err.Error())
		}
		return nil, status.Error(codes.Internal, "repository error")
	}

	sentTime := req.Message.Time.AsTime()
	messageID, err := uc.message.PutID(ctx, user, req.Message.Text, sentTime)
	if err != nil {
		return nil, status.Error(codes.Internal, "repository error")
	}

	message := convert.Text2Message(req.Message, user, messageID)

	err = uc.chat.NewMessage(ctx, user, message)
	if err != nil {
		return nil, status.Error(codes.Internal, "chat error")
	}

	return &desc.SendMessageResponse{
		Addressee: convert.User2Credentials(user),
		Time:      timestamppb.New(sentTime),
	}, nil
}
