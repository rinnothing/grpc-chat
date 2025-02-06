package tui

import (
	"context"
	"sync"

	"github.com/rinnothing/grpc-chat/internal/pkg/model"

	"github.com/charmbracelet/bubbles/list"
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
	askApproval
)

type MessageSender interface {
	Send(ctx context.Context, msg *model.Message) error
}

type MessageRepo interface {
	GetMessages(ctx context.Context, usr *model.User) []*model.Message
}

type Model struct {
	status modelState

	chatsList table.Model
	rowToUser []*model.User
	rowsMx    sync.Mutex

	chatContent    viewport.Model
	msgRepo        MessageRepo
	chatContentStr string
	curUser        *model.User

	messageInput textarea.Model
	sender       MessageSender

	approveMenu    list.Model
	approveChannel chan<- bool
}

type item struct {
	text string
}

func (i item) FilterValue() string {
	return i.text
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
		rowToUser:    make([]*model.User, 0),
		chatContent:  viewport.Model{},
		msgRepo:      msgRepo,
		messageInput: textarea.New(),
		sender:       sender,
		approveMenu:  list.New([]list.Item{item{"Yes"}, item{"No"}}, list.NewDefaultDelegate(), 0, 0),
	}
}

func (m *Model) Init() tea.Cmd {
	return nil
}

type ErrorMsg struct {
	err error
}

// Update passes message to it's two parts or processes it itself when key
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
		case askApproval:
			return m.updateApprove(msg)
		default:
			return m, nil
		}
	}

	// otherwise call messages on both parts and combine the outputs
	var mdl *Model
	cmds := make([]tea.Cmd, 4)
	mdl, cmds[0] = m.updateChats(msg)
	mdl, cmds[1] = mdl.updateBrowseMessages(msg)
	mdl, cmds[2] = mdl.updateWriteMessage(msg)
	mdl, cmds[3] = mdl.updateApprove(msg)

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
	case askApproval:
		return m.viewApprove()
	default:
		panic("unknown model status")
	}
}

