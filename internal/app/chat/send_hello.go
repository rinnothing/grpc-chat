package chat

import (
	"context"
	"errors"

	"github.com/rinnothing/grpc-chat/internal/pkg/convert"
	desc "github.com/rinnothing/grpc-chat/pkg/generated/proto/chat"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (i *Implementation) SendHello(ctx context.Context, req *desc.SendHelloRequest) (*desc.SendHelloResponse, error) {
	if err := validateSendHelloRequest(ctx, req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	allowed, err := i.SendHelloUseCase.SendHello(ctx,
		convert.Text2Message(req.RequestText, convert.Credentials2User(req.Sender, 0), 0))
	if err != nil {
		return nil, status.Error(codes.Internal, "internal connection establishment error")
	}

	return &desc.SendHelloResponse{
		Addressee: req.Sender,
		Allowed:   allowed,
		Time:      timestamppb.Now(),
	}, nil
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
