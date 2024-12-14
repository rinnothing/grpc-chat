package send_goodbye

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
	Disconnect(ctx context.Context, user *model.User) error
}

type Chat interface {
	CloseChat(ctx context.Context, user *model.User, time time.Time) error
}

type Identify interface {
	Myself(ctx context.Context) *model.User
}

type UseCase struct {
	repo        UserRepo
	connections Connections
	chat        Chat
	identify    Identify
}

func New(repo UserRepo, chat Chat, connections Connections, identify Identify) *UseCase {
	return &UseCase{
		repo:        repo,
		connections: connections,
		chat:        chat,
		identify:    identify,
	}
}

func (uc *UseCase) SendGoodbye(ctx context.Context, req *desc.SendGoodbyeRequest) (*desc.SendGoodbyeResponse, error) {
	id, err := uc.repo.GetID(ctx, req.Sender.Username)
	if err != nil {
		return nil, status.Error(codes.Internal, "repository error")
	}

	user := convert.Credentials2User(req.Sender, id)

	err = uc.connections.Disconnect(ctx, user)
	if err != nil {
		if errors.Is(err, connections.ErrAlreadyConnected) {
			return nil, status.Error(codes.Unauthenticated, err.Error())
		}
		return nil, status.Error(codes.Internal, "repository error")
	}

	sentTime := req.Time.AsTime()
	err = uc.chat.CloseChat(ctx, user, sentTime)
	if err != nil {
		return nil, status.Error(codes.Internal, "chat error")
	}

	return &desc.SendGoodbyeResponse{
		Addressee: convert.User2Credentials(uc.identify.Myself(ctx)),
		Time:      timestamppb.Now(),
	}, nil
}
