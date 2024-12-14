package connections

import (
	"bytes"
	"context"
	"errors"
	"sync"

	"github.com/rinnothing/grpc-chat/internal/pkg/model"
)

var (
	ErrAlreadyConnected = errors.New("already connected")
	ErrNotConnected     = errors.New("not connected")
)

type Repo struct {
	mx   *sync.RWMutex
	data map[int]*model.User
}

func NewIdentifyRepo() *Repo {
	return &Repo{
		mx:   new(sync.RWMutex),
		data: make(map[int]*model.User),
	}
}

func (r *Repo) Connect(ctx context.Context, user *model.User) error {
	r.mx.RLock()
	var once sync.Once
	defer once.Do(r.mx.RUnlock)

	if ctx.Err() != nil {
		return ctx.Err()
	}

	val, ok := r.data[user.ID]
	if ok && !bytes.Equal(val.IPv4, user.IPv4) {
		return ErrAlreadyConnected
	}

	r.mx.RUnlock()

	r.mx.Lock()
	defer r.mx.Unlock()

	if ctx.Err() != nil {
		return ctx.Err()
	}

	val, ok = r.data[user.ID]
	if ok && !bytes.Equal(val.IPv4, user.IPv4) {
		return ErrAlreadyConnected
	}

	r.data[user.ID] = user

	return nil
}

func (r *Repo) IsConnected(ctx context.Context, user *model.User) error {
	r.mx.RLock()
	defer r.mx.RUnlock()

	if err := ctx.Err(); err != nil {
		return err
	}

	_, ok := r.data[user.ID]
	if !ok {
		return ErrNotConnected
	}
	return nil
}

func (r *Repo) Disconnect(ctx context.Context, user *model.User) error {
	r.mx.RLock()
	var once sync.Once
	defer once.Do(r.mx.RUnlock)

	if ctx.Err() != nil {
		return ctx.Err()
	}

	_, ok := r.data[user.ID]
	if !ok {
		return ErrNotConnected
	}

	r.mx.RUnlock()

	r.mx.Lock()
	defer r.mx.Unlock()

	if ctx.Err() != nil {
		return ctx.Err()
	}

	_, ok = r.data[user.ID]
	if !ok {
		return ErrNotConnected
	}

	delete(r.data, user.ID)
	return nil
}
