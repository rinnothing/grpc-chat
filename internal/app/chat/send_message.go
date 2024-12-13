package chat

import (
	"context"
	"errors"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	desc "github.com/rinnothing/grpc-chat/pkg/generated/proto/chat"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (i *Implementation) SendMessage(ctx context.Context, req *desc.SendMessageRequest) (*desc.SendMessageResponse, error) {
	if err := validateSendMessageRequest(ctx, req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return i.SendMessageUseCase.SendMessage(ctx, req)
}

func validateSendMessageRequest(ctx context.Context, req *desc.SendMessageRequest) error {
	if req == nil {
		return errors.New("request is nil")
	}
	return validation.ValidateStructWithContext(
		ctx,
		req,
		validation.Field(&req.Sender, validation.Required, isCorrectCredentials),
		validation.Field(&req.Message, validation.Required, isCorrectMessage),
	)
}
