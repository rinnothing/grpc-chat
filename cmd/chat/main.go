package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	service "github.com/rinnothing/grpc-chat/internal/app/chat"
	"github.com/rinnothing/grpc-chat/internal/config"
	"github.com/rinnothing/grpc-chat/internal/pkg/presenter/chat"
	"github.com/rinnothing/grpc-chat/internal/pkg/presenter/dialogue"
	"github.com/rinnothing/grpc-chat/internal/pkg/repository/connections"
	"github.com/rinnothing/grpc-chat/internal/pkg/repository/identify"
	"github.com/rinnothing/grpc-chat/internal/pkg/repository/message"
	"github.com/rinnothing/grpc-chat/internal/pkg/repository/user"
	"github.com/rinnothing/grpc-chat/internal/pkg/usecases/send_goodbye"
	"github.com/rinnothing/grpc-chat/internal/pkg/usecases/send_hello"
	"github.com/rinnothing/grpc-chat/internal/pkg/usecases/send_message"
	desc "github.com/rinnothing/grpc-chat/pkg/generated/proto/chat"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func main() {
	// creating repos to store data
	userRepo := user.NewRepo()
	messageRepo := message.NewRepo()
	connectionsRepo := connections.NewRepo()
	identifyRepo := identify.NewRepo(config.MustGetIPv4(), config.MustGetUsername())

	// creating chat instances (i.e. interface)
	dialoguePresenter := dialogue.NewPresenter(os.Stdin, os.Stdout)
	chatPresenter := chat.NewPresenter(os.Stdout)

	// initializing usecases
	sendHello := send_hello.New(
		userRepo,
		connectionsRepo,
		messageRepo,
		dialoguePresenter,
		chatPresenter,
		identifyRepo,
	)
	sendGoodbye := send_goodbye.New(
		userRepo,
		chatPresenter,
		connectionsRepo,
		identifyRepo,
	)
	sendMessage := send_message.New(
		userRepo,
		connectionsRepo,
		messageRepo,
		chatPresenter,
		identifyRepo,
	)

	// starting a service
	appService := service.New(sendHello, sendMessage, sendGoodbye)

	// starting
	i := implementation{}
	go i.runGrpc(config.MustGetPort(), appService)

	// registering interrupts
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// waiting for interrupt
	<-ctx.Done()
	stop()

	// shutting down the server
	i.server.GracefulStop()
	log.Println("performed graceful stop")
}

type implementation struct {
	server *grpc.Server
}

// runGrpc is a method used to start a grpc service
func (i *implementation) runGrpc(port string, lib desc.ChatInstanceServer) {
	// trying to listen on given grpc port
	lis, err := net.Listen("tcp", "localhost:"+port)
	if err != nil {
		//todo: replace log with slog
		log.Fatalf("failed to listen: %v", err)
	}

	// creating new server and adding given lib for handling
	i.server = grpc.NewServer()
	desc.RegisterChatInstanceServer(i.server, lib)
	grpc_health_v1.RegisterHealthServer(i.server, health.NewServer())
	//todo: replace log with slog
	log.Printf("server listening at %v", lis.Addr())
	if err = i.server.Serve(lis); err != nil {
		panic(err)
	}
}
