package user

import (
	"context"
	"sync"
)

type Repo struct {
	mx    *sync.RWMutex
	users map[string]int
	cnt   int
}

func NewRepo() *Repo {
	return &Repo{
		mx:    new(sync.RWMutex),
		users: make(map[string]int),
	}
}

func (r *Repo) GetID(ctx context.Context, username string) (int, error) {
	r.mx.RLock()
	var once sync.Once
	defer once.Do(r.mx.RUnlock)

	if err := ctx.Err(); err != nil {
		return 0, err
	}

	val, ok := r.users[username]
	if ok {
		return val, nil
	}

	once.Do(r.mx.RUnlock)

	r.mx.Lock()
	defer r.mx.Unlock()

	if err := ctx.Err(); err != nil {
		return 0, err
	}

	val, ok = r.users[username]
	if ok {
		return val, nil
	}

	r.users[username] = r.cnt
	r.cnt++

	return r.cnt - 1, nil
}
