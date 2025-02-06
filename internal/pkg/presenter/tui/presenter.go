package tui

import (
	"context"
	"time"

	msg "github.com/rinnothing/grpc-chat/internal/app/tui"
	"github.com/rinnothing/grpc-chat/internal/pkg/model"
	"github.com/rinnothing/grpc-chat/internal/pkg/presenter/dialogue"

	tea "github.com/charmbracelet/bubbletea"
)

type Presenter struct {
	model *tea.Program
}

func New(model *tea.Program) *Presenter {
	return &Presenter{model: model}
}

func (p *Presenter) NewChat(_ context.Context, user *model.User, message *model.Message) error {
	message.User = user
	p.model.Send(msg.NewChatMsg{Msg: message})

	// todo: add a way to check for errors
	return nil
}

func (p *Presenter) NewMessage(_ context.Context, user *model.User, message *model.Message) error {
	message.User = user
	p.model.Send(msg.NewMessageMsg{Msg: message})

	return nil
}

func (p *Presenter) CloseChat(_ context.Context, user *model.User, time time.Time) error {
	// todo: add something to send when call close chat

	return nil
}

func (p *Presenter) AskAccept(_ context.Context, user *model.User, message *model.Message) error {
	message.User = user

	ch := make(chan bool)
	p.model.Send(msg.AskApprovalMsg{
		Username:       user.Username,
		ApproveChannel: ch,
	})

	answer := <-ch
	if !answer {
		return dialogue.ErrNotAllowed
	}
	return nil
}
