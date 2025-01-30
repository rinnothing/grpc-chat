package message

import (
	"context"
	"sync"
	"time"

	"github.com/rinnothing/grpc-chat/internal/pkg/model"
)

// Repo will be in memory repo for now
// obviously I will replace it with database later
type Repo struct {
	mx       *sync.RWMutex
	messages []Entry
}

type Entry struct {
	userID  int
	message model.Message
}

func NewRepo() *Repo {
	return &Repo{
		mx:       new(sync.RWMutex),
		messages: make([]Entry, 0),
	}
}

func (r *Repo) PutID(ctx context.Context, user *model.User, text string, time time.Time) (int, error) {
	r.mx.Lock()
	defer r.mx.Unlock()

	id := len(r.messages)
	r.messages = append(r.messages, Entry{user.ID, model.Message{
		ID:   id,
		User: user,
		Text: text,
		Time: time,
	}})

	return id, nil
}
