package convert

import (
	"github.com/rinnothing/grpc-chat/internal/pkg/model"
	desc "github.com/rinnothing/grpc-chat/pkg/generated/proto/chat"
)

func User2Credentials(user *model.User) *desc.Credentials {
	return &desc.Credentials{
		Username: user.Username,
		IPv4:     IP2int(user.IPv4),
	}
}

func Credentials2User(cred *desc.Credentials, id int) *model.User {
	return &model.User{
		ID:       id,
		Username: cred.Username,
		IPv4:     Int2IP(cred.IPv4),
	}
}

func Text2Message(text *desc.Message, user *model.User, messageID int) *model.Message {
	return &model.Message{
		ID:   messageID,
		User: user,
		Text: text.Text,
		Time: text.Time.AsTime(),
	}
}
