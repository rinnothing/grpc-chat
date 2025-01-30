package send_message

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

func (uc *UserCase) SendMessage(ctx context.Context, message *model.Message) (*model.User, error) {
	// getting myself
	myself := uc.identify.Myself(ctx)
	if myself == nil {
		return nil, errors.New("can't identify myself")
	}

	// getting user id
	var err error
	message.User.ID, err = uc.repo.GetID(ctx, message.User.Username)
	if err != nil {
		return nil, err
	}

	// checking connection
	err = uc.connections.IsConnected(ctx, message.User)
	if err != nil {
		return nil, err
	}

	// registering incoming message
	message.ID, err = uc.message.PutID(ctx, message.User, message.Text, message.Time)
	if err != nil {
		return nil, err
	}

	// notifying chat about the new message
	err = uc.chat.NewMessage(ctx, message.User, message)
	if err != nil {
		return nil, err
	}

	return myself, nil
}
