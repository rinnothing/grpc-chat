package tui

import (
	"context"
	"sync"

	"github.com/rinnothing/grpc-chat/internal/pkg/model"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

type modelState int

const (
	showChats modelState = iota
	browseMessages
	writeMessage
)

type MessageSender interface {
	Send(ctx context.Context, msg *model.Message) error
}

type MessageRepo interface {
	GetMessages(ctx context.Context, usr *model.User) []*model.Message
}

type Model struct {
	status modelState

	chatsList    table.Model
	rowToUser []*model.User
	rowsMx    sync.Mutex

	chatContent  viewport.Model
	msgRepo MessageRepo
	chatContentStr string
	curUser *model.User

	messageInput textarea.Model
	sender MessageSender
}

func New(sender MessageSender, msgRepo MessageRepo) *Model {
	listColumns := []table.Column{
		{Title: "Sender", Width: 20},
		{Title: "Last message", Width: 15},
		{Title: "Status", Width: 5},
	}

	return &Model{
		status: showChats,
		chatsList: table.New(
			table.WithColumns(listColumns),
			table.WithFocused(true),
		),
		chatContent:  viewport.Model{},
		messageInput: textarea.New(),
		sender:       sender,
		rowToUser:    make([]*model.User, 0),
		msgRepo:      msgRepo,
	}
}

func (m *Model) Init() tea.Cmd {
	return nil
}

// Update passes message to it's two parts or processes it itself when key
func (m *Model) Update(msg tea.Msg) (*Model, tea.Cmd) {
	switch msg := msg.(type) {
	// key messages only belong to one update
	case tea.KeyMsg:
		switch m.status {
		case showChats:
			// check for exit in showChats
			switch msg.String() {
			case "q", "ctrl+c":
				return m, tea.Quit
			}
		case browseMessages:
			return m.updateBrowseMessages(msg)
		case writeMessage:
			return m.updateWriteMessage(msg)
		default:
			return m, nil
		}
	}

	// otherwise call messages on both parts and combine the outputs
	var mdl *Model
	cmds := make([]tea.Cmd, 3)
	mdl, cmds[0] = m.updateChats(msg)
	mdl, cmds[1] = mdl.updateBrowseMessages(msg)
	mdl, cmds[2] = mdl.updateWriteMessage(msg)

	return mdl, tea.Batch(cmds...)
}

func (m *Model) View() string {
	switch m.status {
	case showChats:
		return m.viewChats()
	case browseMessages:
		fallthrough
	case writeMessage:
		return m.viewMessages()
	default:
		panic("unknown model status")
	}
}
