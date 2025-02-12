package tui

import (
	"fmt"
	"sync"
	"time"

	"github.com/rinnothing/grpc-chat/internal/pkg/model"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

type chats struct {
	table.Model
	rowToUser []*model.User
	rowsMx    sync.Mutex
}

func newChats() *chats {
	listColumns := []table.Column{
		{Title: "Sender", Width: 20},
		{Title: "Last message", Width: 15},
		{Title: "Status", Width: 5},
	}

	return &chats{
		Model: table.New(
			table.WithColumns(listColumns),
			table.WithFocused(true),
		),
		rowToUser: make([]*model.User, 0),
		rowsMx:    sync.Mutex{},
	}
}

func (c *chats) Init() tea.Cmd {
	return nil
}

func (c *chats) View() string {
	return c.Model.View()
}

// NewChatMsg is the message that signals that new chat has opened
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

func (c *chats) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		// open selected chat
		case "enter":
			if len(c.rowToUser) == 0 {
				return c, nil
			}

			// tell messages to open chat of user, then switch the status
			c.Update(OpenBranchMsg{c.rowToUser[c.Cursor()]})
			//c.Blur()
			return c, tea.Batch(sendMsg(UpdateStatusMsg{browseMessages}),
				sendMsg(OpenBranchMsg{c.rowToUser[c.Cursor()]}))
		}
	case NewChatMsg:
		c.rowsMx.Lock()

		// adding user to list
		c.rowToUser = append(c.rowToUser, msg.Msg.User)

		// adding new row to user interface
		updatedRows := c.Rows()
		updatedRows = append(updatedRows, msgToRow(msg.Msg))
		c.SetRows(updatedRows)

		c.rowsMx.Unlock()
		return c, nil
	}

	// otherwise pass it to table
	var cmd tea.Cmd
	c.Model, cmd = c.Model.Update(msg)
	return c, cmd
}
