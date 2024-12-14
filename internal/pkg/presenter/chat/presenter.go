package chat

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/rinnothing/grpc-chat/internal/pkg/model"
)

type Presenter struct {
	output io.Writer
}

func NewChatPresenter(output io.Writer) *Presenter {
	return &Presenter{
		output: output,
	}
}

func (p *Presenter) printMessage(user *model.User, message *model.Message) error {
	_, err := fmt.Fprintf(p.output, "got new message from %s at %s:\n%s\n",
		user.Username, message.Time.String(), message.Text)
	return err
}

func (p *Presenter) NewChat(_ context.Context, user *model.User, message *model.Message) error {
	_, err := fmt.Fprintf(p.output, "got new chat with %s at %s\n", user.Username, message.Time.String())
	if err != nil {
		return err
	}

	return p.printMessage(user, message)
}

func (p *Presenter) NewMessage(_ context.Context, user *model.User, message *model.Message) error {
	return p.printMessage(user, message)
}

func (p *Presenter) CloseChat(_ context.Context, user *model.User, time time.Time) error {
	_, err := fmt.Fprintf(p.output, "close the chat with %s at %s\n", user.Username, time.String())
	return err
}
