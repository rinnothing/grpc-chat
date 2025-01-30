package chat

import (
	"context"
	"errors"
	"github.com/rinnothing/grpc-chat/internal/pkg/repository/connections"

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

	usr, err := i.SendGoodbyeUseCase.SendGoodbye(ctx, convert.Credentials2User(req.Sender, 0), req.Time.AsTime())
	if err != nil {
		if errors.Is(err, connections.ErrNotConnected) {
			return nil, status.Error(codes.Unauthenticated, "can't close non-existing connection")
		}
		return nil, status.Error(codes.Internal, "internal connection closing error")
	}

	return &desc.SendGoodbyeResponse{
		Addressee: convert.User2Credentials(usr),
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
