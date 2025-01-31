package chat

import (
	"errors"
	"regexp"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	desc "github.com/rinnothing/grpc-chat/pkg/generated/proto/chat"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	isUsernameFormat     = validation.Match(regexp.MustCompile("^[A-Za-z0-9](_?[A-Za-z0-9])*$"))
	isCorrectCredentials = validation.By(func(value interface{}) error {
		cred, ok := value.(*desc.Credentials)
		if !ok {
			return errors.New("is not type credentials")
		}

		return validation.ValidateStruct(
			value,
			validation.Field(&cred.Username, validation.Required, isUsernameFormat, validation.Length(5, 20)))
	})

	isCorrectSentTimestamp = validation.By(func(value interface{}) error {
		timestamp, ok := value.(*timestamppb.Timestamp)
		if !ok {
			return errors.New("is not type timestamp")
		}

		return validation.Validate(
			timestamp.AsTime(),
			validation.Max(time.Now()),
		)
	})

	isCorrectMessage = validation.By(func(value interface{}) error {
		message, ok := value.(*desc.Message)
		if !ok {
			return errors.New("is not type message")
		}

		return validation.ValidateStruct(
			message,
			validation.Field(&message.Time, validation.Required, isCorrectSentTimestamp),
			validation.Field(&message.Text, validation.Required, validation.Length(0, 1000)),
		)
	})
)
