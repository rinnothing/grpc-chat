package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type approve struct {
	list.Model
	channel chan<- bool
}

func newApprove() *approve {
	return &approve{
		Model: list.New([]list.Item{item{"Yes"}, item{"No"}}, list.NewDefaultDelegate(), 0, 0),
	}
}

func (*approve) Init() tea.Cmd {
	return nil
}

func (a *approve) View() string {
	return a.View()
}

type AskApprovalMsg struct {
	Username       string
	ApproveChannel chan<- bool
}

func (a *approve) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if a.channel == nil {
				break
			}

			if a.Model.Index() == 0 {
				a.channel <- true
			} else {
				a.channel <- false
			}
			a.channel = nil
			return a, nil
		}
	case AskApprovalMsg:
		a.Model.Title = fmt.Sprintf("Incoming message from %s, approve?", msg.Username)
		a.channel = msg.ApproveChannel
		return a, nil
	}

	var cmd tea.Cmd
	a.Model, cmd = a.Model.Update(msg)
	return a, cmd
}
