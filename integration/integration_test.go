package integration

import (
	"context"
	"github.com/rinnothing/grpc-chat/internal/pkg/convert"
	"github.com/rinnothing/grpc-chat/internal/pkg/model"
	"google.golang.org/protobuf/types/known/timestamppb"
	"net"
	"testing"

	desc "github.com/rinnothing/grpc-chat/pkg/generated/proto/chat"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)


func TestDefaultScenario(t *testing.T) {
	ctx := context.Background()
	service := newInstance(t, &model.User{
		Username: "Bob_Bro",
		IPv4:     net.IP{127, 0, 0, 1},
	})

	{
		resp, err := service.sendHello(ctx, "hello")
		require.NoError(t, err)
		require.Equal(t, true, resp.Allowed)
		require.Equal(t, "rinnothing", resp.Addressee.Username)
	}

	{
		resp, err := service.sendMessage(ctx, "how are you?")
		require.NoError(t, err)
		require.Equal(t, "rinnothing", resp.Addressee.Username)
	}


	{
		resp, err := service.sendMessage(ctx, "are you ok?")
		require.NoError(t, err)
		require.Equal(t, "rinnothing", resp.Addressee.Username)
	}


	{
		resp, err := service.sendGoodbye(ctx)
		require.NoError(t, err)
		require.Equal(t, "rinnothing", resp.Addressee.Username)
	}
}

func TestSequence(t *testing.T) {
	ctx := context.Background()
	service := newInstance(t, &model.User{
		Username: "Bob_Bro",
		IPv4:     net.IP{127, 0, 0, 1},
	})

	// goodbye before connection shall not be possible
	{
		_, err := service.sendGoodbye(ctx)
		require.Error(t, err)
	}

	// messaging before connection shall not be possible
	{
		_, err := service.sendMessage(ctx, "no answer")
		require.Error(t, err)
	}
}


type instance struct {
	client desc.ChatInstanceClient
	user *model.User
}

func newInstance(t *testing.T, user *model.User) *instance {
	t.Helper()

	client := newGRPCClient(t)
	return &instance{
		client: client,
		user:   user,
	}
}


func (i *instance) sendHello(ctx context.Context, text string) (*desc.SendHelloResponse, error) {
	return i.client.SendHello(ctx, &desc.SendHelloRequest{
		Sender:      convert.User2Credentials(i.user),
		RequestText: &desc.Message{
			Text: text,
			Time: timestamppb.Now(),
		},
	})
}

func (i *instance) sendMessage(ctx context.Context, text string) (*desc.SendMessageResponse, error) {
	return i.client.SendMessage(ctx, &desc.SendMessageRequest{
		Sender:      convert.User2Credentials(i.user),
		Message: &desc.Message{
			Text: text,
			Time: timestamppb.Now(),
		},
	})
}

func (i *instance) sendGoodbye(ctx context.Context) (*desc.SendGoodbyeResponse, error) {
	return i.client.SendGoodbye(ctx, &desc.SendGoodbyeRequest{
		Sender:      convert.User2Credentials(i.user),
		Time:        timestamppb.Now(),
	})
}

func newGRPCClient(t *testing.T) desc.ChatInstanceClient {
	t.Helper()

	addr := "127.0.0.1:8081"
	c, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)

	return desc.NewChatInstanceClient(c)
}
