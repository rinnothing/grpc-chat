package main

import (
	"context"
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	service "github.com/rinnothing/grpc-chat/internal/app/chat"
	"github.com/rinnothing/grpc-chat/internal/app/tui"
	"github.com/rinnothing/grpc-chat/internal/config"
	"github.com/rinnothing/grpc-chat/internal/pkg/model"
	tuiPres "github.com/rinnothing/grpc-chat/internal/pkg/presenter/tui"
	"github.com/rinnothing/grpc-chat/internal/pkg/repository/connections"
	"github.com/rinnothing/grpc-chat/internal/pkg/repository/identify"
	"github.com/rinnothing/grpc-chat/internal/pkg/repository/message"
	"github.com/rinnothing/grpc-chat/internal/pkg/repository/user"
	"github.com/rinnothing/grpc-chat/internal/pkg/usecases/send_goodbye"
	"github.com/rinnothing/grpc-chat/internal/pkg/usecases/send_hello"
	"github.com/rinnothing/grpc-chat/internal/pkg/usecases/send_message"
	desc "github.com/rinnothing/grpc-chat/pkg/generated/proto/chat"

	tea "github.com/charmbracelet/bubbletea"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

// tempSender is temporary structure used before I write an actual implementation
type tempSender struct{}

func (t *tempSender) Send(ctx context.Context, msg *model.Message) error {
	//TODO: replace with an actual implementation
	return nil
}

type tempGetMessages struct {
	message []*model.Message
}

func (t *tempGetMessages) GetMessages(ctx context.Context, usr *model.User) []*model.Message {
	// TODO: replace with an actual implementation
	return t.message
}

func main() {
	flag.Parse()

	// creating repos to store data
	userRepo := user.NewRepo()
	messageRepo := message.NewRepo()
	connectionsRepo := connections.NewRepo()
	identifyRepo := identify.NewRepo(config.MustGetIPv4(), config.MustGetUsername())

	// creating chat instances (i.e. interface)

	//dialoguePresenter := dialogue.NewPresenter(os.Stdin, os.Stdout)
	//chatPresenter := chat.NewPresenter(os.Stdout)

	usr := &model.User{
		ID:       10,
		Username: "mike_miami",
		IPv4:     net.IPv4(127, 0, 0, 1),
	}
	startTime := time.Now().Truncate(2 * time.Hour)
	mdl := tui.New(&tempSender{}, &tempGetMessages{[]*model.Message{
		{3, usr, "hello", startTime},
		{5, usr, "world", startTime.Add(time.Hour)},
	}})
	program := tea.NewProgram(mdl)
	presenter := tuiPres.New(program)

	// initializing usecases
	sendHello := send_hello.New(
		userRepo,
		connectionsRepo,
		messageRepo,
		presenter,
		presenter,
		identifyRepo,
	)
	sendGoodbye := send_goodbye.New(
		userRepo,
		presenter,
		connectionsRepo,
		identifyRepo,
	)
	sendMessage := send_message.New(
		userRepo,
		connectionsRepo,
		messageRepo,
		presenter,
		identifyRepo,
	)

	// starting a service
	appService := service.New(sendHello, sendMessage, sendGoodbye)

	// starting
	i := implementation{}
	go i.runGrpc(config.MustGetPort(), appService)

	tuiChan := make(chan interface{})
	// running tui
	go func() {
		if _, err := program.Run(); err != nil {
			// TODO: add proper logging
			log.Println(err)
		}

		close(tuiChan)
	}()

	// registering interrupts
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	// waiting for interrupt or interface close
	select {
	case <-tuiChan:
	case <-c:
	}

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
