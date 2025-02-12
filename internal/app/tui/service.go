package tui

import (
	"context"
	"github.com/rinnothing/grpc-chat/internal/pkg/model"

	tea "github.com/charmbracelet/bubbletea"
)

func sendMsg(msg tea.Msg) tea.Cmd {
	return func() tea.Msg {
		return msg
	}
}

type modelState int

const (
	showChats modelState = iota
	browseMessages
	writeMessage
	askApproval
)

// UpdateStatusMsg is a message sent to make Model change its status
type UpdateStatusMsg struct {
	status modelState
}

type MessageSender interface {
	Send(ctx context.Context, msg *model.Message) error
}

type MessageRepo interface {
	GetMessages(ctx context.Context, usr *model.User) []*model.Message
}

type Model struct {
	status modelState

	chats *chats

	messages *messages

	approve *approve
}

type item struct {
	text string
}

func (i item) FilterValue() string {
	return i.text
}

func New(sender MessageSender, msgRepo MessageRepo) *Model {
	m := &Model{
		status:   showChats,
		chats:    newChats(),
		messages: newMessages(msgRepo, sender),
		approve:  newApprove(),
	}

	return m
}

func (m *Model) Init() tea.Cmd {
	return nil
}

type ErrorMsg struct {
	err error
}

func (m *Model) updateState(state modelState) {
	m.status = state
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
			m.messages.Mode = browseMessages
			return m.Update(msg)
		case writeMessage:
			m.messages.Mode = writeMessage
			return m.Update(msg)
		case askApproval:
			return m.approve.Update(msg)
		default:
			return m, nil
		}
	}

	// otherwise call messages on both parts and combine the outputs
	cmds := make([]tea.Cmd, 3)
	// let's forget about references for a minute
	_, cmds[0] = m.chats.Update(msg)
	_, cmds[1] = m.messages.Update(msg)
	_, cmds[2] = m.approve.Update(msg)

	return m, tea.Batch(cmds...)
}

func (m *Model) View() string {
	switch m.status {
	case showChats:
		return m.chats.View()
	case browseMessages:
		fallthrough
	case writeMessage:
		return m.messages.View()
	case askApproval:
		return m.approve.View()
	default:
		panic("unknown model status")
	}
}
