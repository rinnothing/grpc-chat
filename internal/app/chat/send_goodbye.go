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

func (i *Implementation) SendGoodbye(ctx context.Context, req *desc.SendGoodbyeRequest) (*desc.SendGoodbyeResponse, error) {
	if err := validateSendGoodbyeRequest(ctx, req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	err := i.SendGoodbyeUseCase.SendGoodbye(ctx, convert.Credentials2User(req.Sender, 0), req.Time.AsTime())
	if err != nil {
		return nil, status.Error(codes.Internal, "internal connection closing error")
	}

	return &desc.SendGoodbyeResponse{
		Addressee: req.Sender,
		Time:      timestamppb.Now(),
	}, nil
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
