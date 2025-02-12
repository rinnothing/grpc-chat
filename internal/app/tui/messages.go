package tui

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/rinnothing/grpc-chat/internal/pkg/model"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

const gap = "\n"

func (m *messages) View() string {
	return fmt.Sprintf("%s%s%s", m.chatContent.View(), gap, m.messageInput.View())
}

// NewMessageMsg is the message that signals that new message has arrived
type NewMessageMsg struct {
	Msg *model.Message
}

// OpenBranchMsg is the message that signals to open message branch of given user
type OpenBranchMsg struct {
	User *model.User
}

type GotMessagesMsg struct {
	Usr *model.User
	Msg []*model.Message
}

func renderMessage(msg *model.Message) string {
	return fmt.Sprintf("%s\n%s\t%s\n", msg.User.Username, timeFormat(msg.Time), msg.Text)
}

func renderChat(user *model.User, msgs []*model.Message) string {
	b := strings.Builder{}
	b.WriteString(user.Username)
	b.WriteString("\n")

	for _, msg := range msgs {
		b.WriteString(renderMessage(msg))
	}
	return b.String()
}

func newMessages(msgRepo MessageRepo, sender MessageSender) *messages {
	return &messages{
		chatContent:  viewport.Model{},
		msgRepo:      msgRepo,
		messageInput: textarea.New(),
		sender:       sender,
	}
}

type messages struct {
	Mode modelState

	chatContent    viewport.Model
	msgRepo        MessageRepo
	chatContentStr string
	curUser        *model.User

	messageInput textarea.Model
	sender       MessageSender
}

func (m *messages) Init() tea.Cmd {
	return nil
}

func (m *messages) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.Mode {
	case browseMessages:
		return m.updateBrowseMessages(msg)
	case writeMessage:
		return m.updateWriteMessage(msg)
	default:
		panic("invalid mode")
	}
}

// updateBrowseMessages returns nil as first argument if message isn't supported
func (m *messages) updateBrowseMessages(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		// leave to chats
		case "esc":
			return m, sendMsg(UpdateStatusMsg{showChats})
		// start typing message
		case "enter":
			m.messageInput.Focus()
			return m, sendMsg(UpdateStatusMsg{writeMessage})
		}
	case OpenBranchMsg:
		// setting loading screen
		m.chatContent.SetContent(fmt.Sprintf("%s\nloading...", msg.User.Username))
		m.curUser = msg.User
		// send command to retrieve content
		return m, sendMsg(GotMessagesMsg{msg.User, m.msgRepo.GetMessages(context.TODO(), msg.User)})
	case GotMessagesMsg:
		m.chatContentStr = renderChat(msg.Usr, msg.Msg)
		m.chatContent.SetContent(m.chatContentStr)
		return m, nil
	case NewMessageMsg:
		// skip if not in current chat
		if m.curUser == nil || msg.Msg.User.Username != m.curUser.Username {
			return m, nil
		}

		// otherwise add new message to list
		m.chatContentStr = fmt.Sprintf("%s\n%s", m.chatContentStr, renderMessage(msg.Msg))
		return m, nil
	}

	// otherwise pass it down the line
	var cmd tea.Cmd
	m.chatContent, cmd = m.chatContent.Update(msg)
	return m, cmd
}

// updateWriteMessage returns nil as first argument if message isn't supported
func (m *messages) updateWriteMessage(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		// leave to the chat
		case msg.String() == "esc":
			m.messageInput.Blur()
			return m, sendMsg(UpdateStatusMsg{browseMessages})
		case msg.String() == "enter" && msg.Alt:
			var cmd tea.Cmd
			m.messageInput, cmd = m.messageInput.Update(tea.KeyEnter)
			return m, cmd
		case msg.String() == "enter":
			if m.messageInput.Length() == 0 {
				return m, nil
			}

			curText := m.messageInput.Value()
			curUser := m.curUser
			return m, func() tea.Msg {
				err := m.sender.Send(context.TODO(), &model.Message{
					User: curUser,
					Text: curText,
					Time: time.Now(),
				})

				if err != nil {
					return ErrorMsg{err}
				}
				return nil
			}
		}
	}

	// otherwise pass it down the line
	var cmd tea.Cmd
	m.messageInput, cmd = m.messageInput.Update(msg)
	return m, cmd
}
