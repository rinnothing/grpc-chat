package dialogue

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/rinnothing/grpc-chat/internal/pkg/model"
)

var ErrNotAllowed = errors.New("user not allowed")

type Presenter struct {
	input  io.Reader
	output io.Writer
}

func NewDialoguePresenter(input io.Reader, output io.Writer) *Presenter {
	return &Presenter{
		input:  input,
		output: output,
	}
}

func (p *Presenter) AskAccept(_ context.Context, user *model.User, message *model.Message) error {
	_, err := fmt.Fprintf(p.output,
		"User %s wants to chat with you:\n"+
			"%s\n"+
			"allow | deny\n",
		user.Username, message.Text,
	)
	if err != nil {
		return err
	}

	var answer string
	for {
		_, err = fmt.Fscanln(p.input, &answer)
		if err != nil {
			return err
		}

		switch answer {
		case "allow":
			return nil
		case "deny":
			return ErrNotAllowed
		default:
			_, err = fmt.Fprintf(p.output, "%s is not a correct option, choose allow or deny\n", answer)
		}
	}
}
