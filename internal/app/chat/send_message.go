package chat

import (
	"context"
	"errors"
	"github.com/rinnothing/grpc-chat/internal/pkg/repository/connections"
	"time"

	"github.com/rinnothing/grpc-chat/internal/pkg/convert"
	desc "github.com/rinnothing/grpc-chat/pkg/generated/proto/chat"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (i *Implementation) SendMessage(ctx context.Context, req *desc.SendMessageRequest) (*desc.SendMessageResponse, error) {
	if err := validateSendMessageRequest(ctx, req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	usr, err := i.SendMessageUseCase.SendMessage(ctx,
		convert.Text2Message(req.Message, convert.Credentials2User(req.Sender, 0), 0))
	if err != nil {
		if errors.Is(err, connections.ErrNotConnected) {
			return nil, status.Error(codes.Unauthenticated, "can't send message on non-existing connection")
		}
		return nil, status.Error(codes.Internal, "internal incoming message proceeding error")
	}

	return &desc.SendMessageResponse{
		Addressee: convert.User2Credentials(usr),
		Time:      timestamppb.New(time.Now()),
	}, nil
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
