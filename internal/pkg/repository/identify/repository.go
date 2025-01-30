package identify

import (
	"context"
	"net"

	"github.com/rinnothing/grpc-chat/internal/pkg/model"
)

type Repo struct {
	user model.User
}

func NewRepo(IPv4 net.IP, username string) *Repo {
	return &Repo{
		user: model.User{
			Username: username,
			IPv4:     IPv4,
		},
	}
}

func (r *Repo) Myself(_ context.Context) *model.User {
	return &r.user
}
