package send_hello

import (
	"context"
	"errors"
	"time"

	"github.com/rinnothing/grpc-chat/internal/config"
	"github.com/rinnothing/grpc-chat/internal/pkg/model"
	"github.com/rinnothing/grpc-chat/internal/pkg/presenter/dialogue"
	"github.com/rinnothing/grpc-chat/internal/pkg/repository/connections"
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

func (uc *UseCase) SendHello(ctx context.Context, request *model.Message) (bool, *model.User, error) {
	// getting myself
	myself := uc.identify.Myself(ctx)
	if myself == nil {
		return false, nil, errors.New("can't identify myself")
	}

	// getting user id
	var err error
	request.User.ID, err = uc.repo.GetID(ctx, request.User.Username)
	if err != nil {
		return false, myself, err
	}

	// trying to make connect
	err = uc.connections.Connect(ctx, request.User)
	if err != nil {
		if errors.Is(err, connections.ErrAlreadyConnected) {
			return false, myself, nil
		}
		return false, myself, err
	}

	// getting message id
	request.ID, err = uc.message.PutID(ctx, request.User, request.Text, request.Time)
	if err != nil {
		return false, myself, err
	}

	// check for automatic accept
	if !config.MustGetAcceptAll() {
		// ask if connection request is accepted
		err = uc.dialogue.AskAccept(ctx, request.User, request)
		if err != nil {
			if errors.Is(err, dialogue.ErrNotAllowed) {
				return false, myself, nil
			}
			return false, myself, err
		}
	}

	// displaying new message in chat
	err = uc.chat.NewChat(ctx, request.User, request)
	if err != nil {
		return true, myself, err
	}

	return true, myself, nil
}
