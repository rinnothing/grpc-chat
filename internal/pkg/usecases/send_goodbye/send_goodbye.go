package send_goodbye

import (
	"context"
	"errors"
	"time"

	"github.com/rinnothing/grpc-chat/internal/pkg/model"
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

func (uc *UseCase) SendGoodbye(ctx context.Context, sender *model.User, time time.Time) (*model.User, error) {
	// getting myself
	myself := uc.identify.Myself(ctx)
	if myself == nil {
		return nil, errors.New("can't identify myself")
	}

	// getting user id
	var err error
	sender.ID, err = uc.repo.GetID(ctx, sender.Username)
	if err != nil {
		return nil, err
	}

	// disconnecting
	err = uc.connections.Disconnect(ctx, sender)
	if err != nil {
		return nil, err
	}

	// closing chat
	err = uc.chat.CloseChat(ctx, sender, time)
	if err != nil {
		return nil, err
	}

	return myself, nil
}
