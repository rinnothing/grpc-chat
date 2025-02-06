package tui

import (
	"fmt"
	"time"

	"github.com/rinnothing/grpc-chat/internal/pkg/model"

	tea "github.com/charmbracelet/bubbletea"
)

func (m *Model) viewChats() string {
	return m.chatsList.View()
}

//NewChatMsg is the message that signals that new chat has opened
type NewChatMsg struct {
	Msg *model.Message
}

func timeFormat(t time.Time) string {
	return fmt.Sprintf("%d.%d.%d %d:%d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute())
}

func msgToRow(msg *model.Message) []string {
	return []string{
		msg.User.Username,
		timeFormat(msg.Time),
		"Unread",
	}
}

// updateChats returns nil as first argument if message isn't supported
func (m *Model) updateChats(msg tea.Msg) (*Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		// open selected chat
		case "enter":
			if len(m.rowToUser) == 0 {
				return m, nil
			}

			// tell messages to open chat of user, then switch the status
			m.Update(OpenBranchMsg{m.rowToUser[m.chatsList.Cursor()]})
			m.chatsList.Blur()
			m.status = browseMessages
		}
	case NewChatMsg:
		m.rowsMx.Lock()

		// adding user to list
		m.rowToUser = append(m.rowToUser, msg.Msg.User)

		// adding new row to user interface
		updatedRows := m.chatsList.Rows()
		updatedRows = append(updatedRows, msgToRow(msg.Msg))
		m.chatsList.SetRows(updatedRows)

		m.rowsMx.Unlock()
	}

	// otherwise pass it to table
	var cmd tea.Cmd
	m.chatsList, cmd = m.chatsList.Update(msg)
	return m, cmd
}
