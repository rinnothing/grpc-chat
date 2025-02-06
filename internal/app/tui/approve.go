package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

func (m *Model) viewApprove() string {
	return m.approveMenu.View()
}

type AskApprovalMsg struct {
	Username       string
	ApproveChannel chan<- bool
}

func (m *Model) updateApprove(msg tea.Msg) (*Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if m.approveChannel == nil {
				break
			}

			if m.approveMenu.Index() == 0 {
				m.approveChannel <- true
			} else {
				m.approveChannel <- false
			}
			m.approveChannel = nil
		}
	case AskApprovalMsg:
		m.approveMenu.Title = fmt.Sprintf("Incoming message from %s, approve?", msg.Username)
		m.approveChannel = msg.ApproveChannel
	}

	var cmd tea.Cmd
	m.approveMenu, cmd = m.approveMenu.Update(msg)
	return m, cmd
}
