package chat

import (
	"context"
	"errors"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	desc "github.com/rinnothing/grpc-chat/pkg/generated/proto/chat"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (i *Implementation) SendGoodbye(ctx context.Context, req *desc.SendGoodbyeRequest) (*desc.SendGoodbyeResponse, error) {
	if err := validateSendGoodbyeRequest(ctx, req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return i.SendGoodbyeUseCase.SendGoodbye(ctx, req)
}

func validateSendGoodbyeRequest(ctx context.Context, req *desc.SendGoodbyeRequest) error {
	if req == nil {
		return errors.New("request is nil")
	}
	return validation.ValidateStructWithContext(
		ctx,
		req,
		validation.Field(&req.Sender, validation.Required, isCorrectCredentials),
		validation.Field(&req.Time, validation.Required, isCorrectCredentials),
	)
}
