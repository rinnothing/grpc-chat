package chat

import (
	"context"
	"errors"
	desc "github.com/rinnothing/grpc-chat/pkg/generated/proto/chat"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (i *Implementation) SendHello(ctx context.Context, req *desc.SendHelloRequest) (*desc.SendHelloResponse, error) {
	if err := validateSendHelloRequest(ctx, req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return i.SendHelloUseCase.SendHello(ctx, req)
}

func validateSendHelloRequest(ctx context.Context, req *desc.SendHelloRequest) error {
	if req == nil {
		return errors.New("request is nil")
	}
	return validation.ValidateStructWithContext(
		ctx,
		req,
		validation.Field(&req.Sender, validation.Required, isCorrectCredentials),
		validation.Field(&req.RequestText, validation.Required, isCorrectMessage),
	)
}
